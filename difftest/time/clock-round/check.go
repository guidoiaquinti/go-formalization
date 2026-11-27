package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Wall                uint64
	Ext, Duration, Half int64
	Loc                 int
	Lean                string
}

func main() {
	f, err := os.CreateTemp("", "golean-time-clock-round")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	enc := json.NewEncoder(f)
	compare := func(label string, cases []fixture) {
		var in bytes.Buffer
		for _, c := range cases {
			fmt.Fprintln(&in, c.Wall, c.Ext, c.Loc, c.Duration, c.Half)
		}
		for i, row := range check.LeanRows("godiff-time-clock-round", len(cases), &in) {
			c := cases[i]
			c.Lean = row
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d Lean rows recorded\n", label, len(cases))
	}
	divisors := []int64{math.MinInt64, -1000000000, -1, 0, 1, 2, 3, 5, 7, 10, 1000, 100000000, 999999999, 1000000000, 1000000001, 2000000000, 3600000000000, math.MaxInt64}
	var fixed []fixture
	for _, wall := range []uint64{0, 1, 500000000, 999999999, (1 << 30) - 1, 1 << 30, 1 << 63, (1 << 63) + 999999999, math.MaxUint64} {
		for _, ext := range []int64{math.MinInt64, math.MinInt64 + 1, -59453308800, -2, -1, 0, 1, 59453308800, math.MaxInt64} {
			for _, d := range divisors {
				for _, loc := range []int{0, 3} {
					fixed = append(fixed, fixture{Wall: wall, Ext: ext, Duration: d, Half: d / 2, Loc: loc})
				}
			}
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			d := int64(r.Uint64())
			if i%2 == 0 {
				d = divisors[r.Intn(len(divisors))]
			}
			cases[i] = fixture{Wall: r.Uint64(), Ext: int64(r.Uint64()), Duration: d, Half: int64(r.Uint64()), Loc: r.Intn(4)}
		}
		compare(fmt.Sprintf("seed %d", seed), cases)
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_TIME_CLOCK_ROUND", f.Name()))
	must(check.TestInPackage("time", "TestGoLeanClockRound",
		"difftest/time/location/testdata/export_test.go", "difftest/time/clock-arithmetic/testdata/bridge_test.go",
		"difftest/time/clock-round/testdata/export_test.go", "difftest/time/clock-round/testdata/clock_round_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
