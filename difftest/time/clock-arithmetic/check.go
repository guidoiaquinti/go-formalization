package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type stamp struct {
	wall uint64
	ext  int64
	loc  int
}

type fixture struct {
	t, u  stamp
	delta int64
}

func main() {
	f, err := os.CreateTemp("", "golean-clock-arithmetic-")
	must(err)
	defer os.Remove(f.Name())
	out := bufio.NewWriter(f)
	const mark = uint64(1 << 63)
	const fine = uint64(1 << 30)
	edge := []stamp{
		{0, 0, 0}, {0, 1, 1}, {1, 0, 2}, {999999999, 0, 1}, {fine - 1, 0, 0},
		{fine, math.MinInt64, 2}, {fine - 1, math.MinInt64, 1}, {fine - 1, math.MaxInt64, 1},
		{0, math.MinInt64 + 1, 1}, {0, math.MaxInt64 - 1, 2}, {0, math.MaxInt64, 0},
		{mark, 0, 0}, {mark | 17, -1, 1}, {mark | fine | 999999999, 1, 2},
		{mark | fine, math.MaxInt64, 1}, {mark | fine, math.MinInt64, 2},
		{mark | (1 << 63) - fine | 5, math.MaxInt64 - 1, 2}, {mark | fine - 1, math.MinInt64 + 1, 0},
		{0, 59453308800, 1}, {17, 59453308800, 2},
		{854775807, 9223372036, 1}, {854775808, 9223372036, 2},
		{145224192, -9223372037, 1}, {145224191, -9223372037, 2},
	}
	var fixed []fixture
	for _, t := range edge {
		for _, u := range edge {
			for _, d := range []int64{math.MinInt64, -1000000001, -1000000000, -1, 0, 1, 999999999, 1000000000, math.MaxInt64} {
				fixed = append(fixed, fixture{t, u, d})
			}
		}
	}
	compare("fixed", fixed, out)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 5000)
			for i := range cases {
				t := randomStamp(rng)
				u := randomStamp(rng)
				delta := int64(rng.Uint64())
				switch i % 6 {
				case 0:
					u = t
				case 1:
					u = t
					u.ext--
				case 2:
					u = t
					u.wall ^= 1
				case 3:
					delta = int64(rng.Intn(17)) - 8
				case 4:
					u.wall = t.wall &^ mark
					u.ext = int64(rng.Uint64())
				}
				cases[i] = fixture{t, u, delta}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases, out)
		}
	}
	must(out.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_CLOCK_ARITHMETIC", f.Name()))
	must(check.TestInPackage("time", "TestGoLeanClockArithmetic",
		"difftest/time/clock-arithmetic/testdata/bridge_test.go",
		"difftest/time/clock-arithmetic/testdata/arithmetic_test.go"))
}

func randomStamp(rng *rand.Rand) stamp {
	return stamp{rng.Uint64(), int64(rng.Uint64()), rng.Intn(3)}
}

func compare(label string, cases []fixture, out *bufio.Writer) {
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&input, c.t.wall, c.t.ext, c.t.loc, c.u.wall, c.u.ext, c.u.loc, c.delta)
	}
	for i, result := range check.LeanRows("godiff-time-clock-arithmetic", len(cases), &input) {
		c := cases[i]
		_, err := fmt.Fprintln(out, c.t.wall, c.t.ext, c.t.loc, c.u.wall, c.u.ext, c.u.loc, c.delta, result)
		must(err)
	}
	fmt.Printf("Time arithmetic %s: %d Lean rows\n", label, len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
