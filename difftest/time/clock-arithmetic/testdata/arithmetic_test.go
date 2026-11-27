package time_test

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"math"
	"math/big"
	"os"
	"testing"
	"time"
)

const quantum = int64(1_000_000_000)
const nanosField = uint64(1 << 30)
const marked = uint64(1 << 63)
const firstPack = int64(59_453_308_800)

type rawClock struct {
	wall     uint64
	ext      int64
	location int
}

type output struct {
	add       rawClock
	sub, mono int64
}

func isMarked(a rawClock) bool { return a.wall >= marked }
func nanos(a rawClock) int64   { return int64(a.wall % nanosField) }
func seconds(a rawClock) int64 {
	if isMarked(a) {
		return firstPack + int64((a.wall-marked)/nanosField)
	}
	return a.ext
}

func clamp(x *big.Int, lower int64) int64 {
	if x.IsInt64() {
		return x.Int64()
	}
	if x.Sign() < 0 {
		return lower
	}
	return math.MaxInt64
}

func subReading(t, u int64) int64 {
	return clamp(new(big.Int).Sub(big.NewInt(t), big.NewInt(u)), math.MinInt64)
}

func withDuration(a rawClock, d int64) rawClock {
	whole, fractional := d/quantum, d%quantum
	combined := nanos(a) + fractional
	// Clamp the correction to one carry or one borrow, including
	// raw wall nanoseconds past the public Time invariant.
	normalized := combined % quantum
	if normalized < 0 {
		normalized += quantum
	}
	correction := (combined - normalized) / quantum
	if correction > 1 {
		correction = 1
	}
	wallNanos := uint64(combined - correction*quantum)
	wallSeconds := seconds(a)
	delta := whole + correction
	newSec := new(big.Int).Add(big.NewInt(wallSeconds), big.NewInt(delta))
	if !isMarked(a) {
		a.wall = a.wall/nanosField*nanosField + wallNanos
		a.ext = clamp(newSec, -math.MaxInt64)
		return a
	}
	packed := new(big.Int).Sub(new(big.Int).Set(newSec), big.NewInt(firstPack))
	if packed.Sign() >= 0 && packed.BitLen() <= 33 {
		updated := new(big.Int).Add(big.NewInt(a.ext), big.NewInt(d))
		if updated.IsInt64() {
			a.wall = marked + packed.Uint64()*nanosField + wallNanos
			a.ext = updated.Int64()
			return a
		}
	}
	a.wall = wallNanos
	a.ext = clamp(newSec, -math.MaxInt64)
	return a
}

func wallCompare(a, b rawClock) int {
	if order := cmp.Compare(seconds(a), seconds(b)); order != 0 {
		return order
	}
	return cmp.Compare(nanos(a), nanos(b))
}

func wrapped(x *big.Int) int64 {
	// Modulo 2^64 establishes the wrapped two’s-complement word.
	modulus := new(big.Int).Lsh(big.NewInt(1), 64)
	return int64(new(big.Int).Mod(x, modulus).Uint64())
}

func difference(a, b rawClock) int64 {
	if isMarked(a) && isMarked(b) {
		return subReading(a.ext, b.ext)
	}
	sec := new(big.Int).Sub(big.NewInt(seconds(a)), big.NewInt(seconds(b)))
	total := new(big.Int).Add(sec.Mul(sec, big.NewInt(quantum)), big.NewInt(nanos(a)-nanos(b)))
	delta := wrapped(total)
	if wallCompare(withDuration(b, delta), a) == 0 {
		return delta
	}
	if wallCompare(a, b) < 0 {
		return math.MinInt64
	}
	return math.MaxInt64
}

func TestGoLeanClockArithmetic(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_CLOCK_ARITHMETIC"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	in := bufio.NewReader(f)
	location := []*time.Location{nil, time.FixedZone("offset", 123), time.FixedZone("offset", 123)}
	for row := 0; ; row++ {
		var a, b rawClock
		var d int64
		var want output
		n, err := fmt.Fscan(in, &a.wall, &a.ext, &a.location, &b.wall, &b.ext, &b.location, &d,
			&want.add.wall, &want.add.ext, &want.add.location, &want.sub, &want.mono)
		if n == 0 && err == io.EOF && row > 0 {
			t.Logf("%d pinned native rows per declaration", row)
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		left := time.GoLeanArithmeticTime(a.wall, a.ext, location[a.location])
		right := time.GoLeanArithmeticTime(b.wall, b.ext, location[b.location])
		wall, ext, loc := time.GoLeanArithmeticFields(left.Add(time.Duration(d)))
		if loc != location[a.location] {
			t.Fatalf("row %d changed location identity", row)
		}
		got := output{rawClock{wall, ext, a.location}, int64(left.Sub(right)), int64(time.GoLeanArithmeticSubMono(a.ext, b.ext))}
		reference := output{withDuration(a, d), difference(a, b), subReading(a.ext, b.ext)}
		if got != want || got != reference {
			t.Fatalf("row %d a=%+v b=%+v d=%d: Go %+v, Lean %+v, independent %+v", row, a, b, d, got, want, reference)
		}
	}
}
