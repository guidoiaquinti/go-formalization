package netip

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func goLeanU128Wire(u uint128) string { return fmt.Sprintf("%d,%d", u.hi, u.lo) }

func goLeanCatchU128(callback func() string, reason string) (out string) {
	defer func() {
		if v := recover(); v != nil {
			err, ok := v.(runtime.Error)
			if !ok || !strings.Contains(err.Error(), reason) {
				panic(v)
			}
			out = "panic"
		}
	}()
	return callback()
}

func goLeanU128Mask(n int) *big.Int {
	if n < 0 || n > 128 {
		return nil
	}
	return new(big.Int).Lsh(new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), uint(n)), big.NewInt(1)), uint(128-n))
}

func goLeanBigU128(u uint128) *big.Int {
	return new(big.Int).Or(new(big.Int).Lsh(new(big.Int).SetUint64(u.hi), 64), new(big.Int).SetUint64(u.lo))
}

func goLeanBigU128Wire(x *big.Int) string {
	if x == nil {
		return "panic"
	}
	modulus := new(big.Int).Lsh(big.NewInt(1), 128)
	low := new(big.Int).Mod(x, modulus)
	return goLeanU128Wire(uint128{new(big.Int).Rsh(low, 64).Uint64(), low.Uint64()})
}

func TestGoLeanUint128(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_UINT128_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		var c struct {
			N                int
			Bit              uint8
			Missing          bool
			Hi, Lo, VHi, VLo uint64
			Want             string
		}
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		u, v := uint128{c.Hi, c.Lo}, uint128{c.VHi, c.VLo}
		var receiver *uint128
		if !c.Missing {
			copy := u
			receiver = &copy
		}
		halves := goLeanCatchU128(func() string {
			fields := receiver.halves()
			if fields[0] != &receiver.hi || fields[1] != &receiver.lo || fields[0] == fields[1] {
				t.Fatal("incorrect pointer aliases")
			}
			original := uint128{*fields[0], *fields[1]}
			*fields[0] = v.hi
			if receiver.hi != v.hi || receiver.lo != u.lo {
				t.Fatal("writing high reference affected wrong field")
			}
			*fields[1] = v.lo
			if *receiver != v {
				t.Fatal("writing low reference affected wrong field")
			}
			return goLeanU128Wire(original) + ":" + goLeanU128Wire(*receiver)
		}, "nil pointer dereference")
		catchMask := func(op func() uint128) string {
			return goLeanCatchU128(func() string { return goLeanU128Wire(op()) }, "negative shift amount")
		}
		native := []string{
			catchMask(func() uint128 { return mask6(c.N) }), strconv.FormatBool(u.isZero()),
			goLeanU128Wire(u.and(v)), goLeanU128Wire(u.xor(v)), goLeanU128Wire(u.or(v)), goLeanU128Wire(u.not()),
			goLeanU128Wire(u.subOne()), goLeanU128Wire(u.addOne()), halves,
			catchMask(func() uint128 { return u.bitsSetFrom(c.Bit) }),
			catchMask(func() uint128 { return u.bitsClearedFrom(c.Bit) }),
		}
		wantHalves := "panic"
		if !c.Missing {
			wantHalves = goLeanU128Wire(u) + ":" + goLeanU128Wire(v)
		}
		set, clear := "panic", "panic"
		x, y := goLeanBigU128(u), goLeanBigU128(v)
		if mask := goLeanU128Mask(int(c.Bit)); mask != nil {
			set = goLeanBigU128Wire(new(big.Int).Or(x, new(big.Int).Xor(mask, goLeanU128Mask(128))))
			clear = goLeanBigU128Wire(new(big.Int).And(x, mask))
		}
		want := []string{
			goLeanBigU128Wire(goLeanU128Mask(c.N)), strconv.FormatBool(x.Sign() == 0),
			goLeanBigU128Wire(new(big.Int).And(x, y)), goLeanBigU128Wire(new(big.Int).Xor(x, y)),
			goLeanBigU128Wire(new(big.Int).Or(x, y)), goLeanBigU128Wire(new(big.Int).Not(x)),
			goLeanBigU128Wire(new(big.Int).Sub(x, big.NewInt(1))), goLeanBigU128Wire(new(big.Int).Add(x, big.NewInt(1))),
			wantHalves, set, clear,
		}
		if strings.Join(native, "|") != c.Want || strings.Join(want, "|") != c.Want {
			t.Fatalf("%+v: native %v, oracle %v", c, native, want)
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	t.Logf("native uint128 rows per declaration: %d", count)
}
