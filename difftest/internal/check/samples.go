package check

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"runtime"
)

func CompareSamples[T any](binary string, edges []T, sample func(*rand.Rand, int) T,
	input, actual, expected func(T) string) {
	compare := func(label string, cases []T) {
		var stdin bytes.Buffer
		for _, c := range cases {
			fmt.Fprintln(&stdin, input(c))
		}
		for i, row := range LeanRows(binary, len(cases), &stdin) {
			c := cases[i]
			got, want := actual(c), expected(c)
			if got != want || row != want {
				log.Fatalf("%s %s input %s: Go=%q, Lean=%q, expected=%q", binary, label, input(c), got, row, want)
			}
		}
		fmt.Printf("%s %s: %d inputs agree.\n", binary, label, len(cases))
	}
	compare("boundaries", edges)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]T, 20_000)
		for i := range cases {
			cases[i] = sample(rng, i)
		}
		compare(fmt.Sprintf("seed=%d", seed), cases)
	}
}

func WordBoundaries(width int) []uint64 {
	half, top := uint64(1)<<(width/2), uint64(1)<<(width-1)
	mask := ^uint64(0) >> (64 - width)
	return []uint64{0, 1, 2, 3, half - 1, half, half + 1, top - 1, top, top + 1, mask - 1, mask}
}

// FormatDivisionPanic preserves Go's two documented panic classes.
func FormatDivisionPanic(p any) string {
	if e, ok := p.(runtime.Error); ok {
		switch e.Error() {
		case "runtime error: integer divide by zero":
			return "divide-by-zero"
		case "runtime error: integer overflow":
			return "quotient-overflow"
		}
	}
	return fmt.Sprintf("unexpected panic (%T): %v", p, p)
}
