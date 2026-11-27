package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const (
	firstPacked = 59453308800
	packedSpan  = 8589934591
	marker      = uint64(1 << 63)
	unit        = uint64(1 << 30)
)

type fixture struct {
	Wall, Ext, Argument uint64
	Location            int
	Want                string
}

func oracle(c fixture) string {
	format := func(wall, ext uint64) string { return fmt.Sprintf("%d,%d,%d", wall, ext, c.Location) }
	setWall, setExt := c.Wall, c.Ext
	marked := c.Wall >= marker
	sec := int64(c.Ext)
	if marked {
		setExt = c.Argument
	} else if sec >= firstPacked && sec <= firstPacked+packedSpan {
		setWall = marker + (c.Wall | uint64(sec-firstPacked)*unit)
		setExt = c.Argument
	}
	addWall := c.Wall
	delta := int64(c.Argument)
	if marked {
		oldPacked := int64((c.Wall - marker) / unit)
		// Bounds expressed on the addend avoid signed overflow in the test.
		if -oldPacked <= delta && delta <= packedSpan-oldPacked {
			addWall = marker + uint64(oldPacked+delta)*unit + c.Wall%unit
			return format(setWall, setExt) + "|" + format(addWall, c.Ext)
		}
		addWall = c.Wall % unit
		sec = firstPacked + oldPacked
	}
	var sum big.Int
	sum.Add(big.NewInt(sec), big.NewInt(delta))
	addExt := int64(math.MaxInt64)
	if sum.IsInt64() {
		addExt = sum.Int64()
	} else if sum.Sign() < 0 {
		addExt = -math.MaxInt64
	}
	return format(setWall, setExt) + "|" + format(addWall, uint64(addExt))
}

func main() {
	f, err := os.CreateTemp("", "golean-time-clock-mutation")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	w := bufio.NewWriter(f)
	enc := gob.NewEncoder(w)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprintln(&input, c.Wall, c.Ext, c.Location, c.Argument)
		}
		for i, got := range check.LeanRows("godiff-time-clock-mutation", len(cases), &input) {
			c := cases[i]
			c.Want = oracle(c)
			if got != c.Want {
				log.Fatalf("%s row %d: Lean %q expected %q", label, i, got, c.Want)
			}
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d Lean rows agree\n", label, len(cases))
	}
	var fixed []fixture
	walls := []uint64{0, 17, unit - 1, unit, math.MaxInt64, marker, marker + unit - 1, marker + unit + 17, math.MaxUint64}
	exts := []int64{math.MinInt64, math.MinInt64 + 1, -1, 0, 1, firstPacked - 1, firstPacked, firstPacked + 1, firstPacked + packedSpan, firstPacked + packedSpan + 1, math.MaxInt64 - 1, math.MaxInt64}
	args := []int64{math.MinInt64, math.MinInt64 + 1, -firstPacked - packedSpan, -firstPacked, -packedSpan - 1, -packedSpan, -2, -1, 0, 1, 2, packedSpan, packedSpan + 1, math.MaxInt64 - 1, math.MaxInt64}
	for _, wall := range walls {
		for _, ext := range exts {
			for _, arg := range args {
				for loc := range 3 {
					fixed = append(fixed, fixture{Wall: wall, Ext: uint64(ext), Argument: uint64(arg), Location: loc})
				}
			}
		}
	}
	compare("fixed", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 5000)
			for i := range cases {
				c := fixture{Wall: r.Uint64(), Ext: r.Uint64(), Argument: r.Uint64(), Location: r.Intn(3)}
				switch i % 4 {
				case 0:
					c.Ext = firstPacked + uint64(r.Int63n(packedSpan+1))
				case 1:
					c.Argument = uint64(r.Int63n(2*packedSpan+1) - packedSpan)
				}
				cases[i] = c
			}
			compare(fmt.Sprintf("seed %d block %d", seed, block), cases)
		}
	}
	must(w.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_TIME_CLOCK_MUTATION_ROWS", f.Name()))
	must(check.TestInPackage("time", "TestGoLeanClockMutation", "difftest/time/clock-mutation/testdata/bridge_test.go",
		"difftest/time/clock-mutation/testdata/mutation_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
