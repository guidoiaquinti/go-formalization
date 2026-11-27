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

type clock struct {
	wall uint64
	ext  int64
}

type pair struct{ t, u clock }

const mono = uint64(1) << 63

func main() {
	file, err := os.CreateTemp("", "golean-time-compare-")
	must(err)
	defer os.Remove(file.Name())
	out := bufio.NewWriter(file)
	var clocks []clock
	for _, ns := range []uint64{0, 1, 499_999_999, 999_999_999, 1<<30 - 1} {
		for _, sec := range []int64{math.MinInt64, -1, 0, 1, 59_453_308_800, 62_135_596_800, math.MaxInt64} {
			clocks = append(clocks, clock{ns, sec})
		}
		for _, field := range []uint64{0, 1, 1 << 32, 1<<33 - 1} {
			for _, ext := range []int64{math.MinInt64, -1, 0, 1, math.MaxInt64} {
				clocks = append(clocks, clock{mono | field<<30 | ns, ext})
			}
		}
	}
	clocks = append(clocks, clock{1 << 62, 0}, clock{1 << 30, 0})
	fixed := make([]pair, 0, len(clocks)*len(clocks))
	for _, t := range clocks {
		for _, u := range clocks {
			fixed = append(fixed, pair{t, u})
		}
	}
	compare("fixed", fixed, out)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			rows := make([]pair, 5000)
			for i := range rows {
				t := clock{rng.Uint64(), int64(rng.Uint64())}
				u := clock{rng.Uint64(), int64(rng.Uint64())}
				switch i % 6 {
				case 0:
					u = t
				case 1:
					u.ext = t.ext
				case 2:
					u.wall = t.wall ^ 1
					u.ext = t.ext
				case 3:
					u.wall &^= mono
					u.ext = referenceSeconds(t.wall, t.ext)
				}
				rows[i] = pair{t, u}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), rows, out)
		}
	}
	must(out.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_CLOCK_COMPARE", file.Name()))
	must(check.TestInPackage("time", "TestGoLeanClockCompare",
		"difftest/time/time-compare/testdata/bridge_test.go",
		"difftest/time/time-compare/testdata/compare_test.go"))
}

func referenceSeconds(wall uint64, ext int64) int64 {
	if wall < mono {
		return ext
	}
	return 59_453_308_800 + int64((wall-mono)/(1<<30))
}

func compare(label string, rows []pair, out *bufio.Writer) {
	var input bytes.Buffer
	for _, r := range rows {
		fmt.Fprintln(&input, r.t.wall, r.t.ext, r.u.wall, r.u.ext)
	}
	for i, got := range check.LeanRows("godiff-time-compare", len(rows), &input) {
		r := rows[i]
		_, err := fmt.Fprintln(out, r.t.wall, r.t.ext, r.u.wall, r.u.ext, got)
		must(err)
	}
	fmt.Printf("time comparison %s: %d Lean rows\n", label, len(rows))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
