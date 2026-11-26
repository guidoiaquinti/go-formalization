package check

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"strconv"
)

// SampleUnary covers both neighbors of each power of two, single cleared
// bits, and reproducible uniform samples over the selected word width.
func SampleUnary(binary string, width int, actual, expected func(uint64) uint64) {
	mask := ^uint64(0) >> (64 - width)
	compare := func(label string, words []uint64) {
		var input bytes.Buffer
		for _, x := range words {
			fmt.Fprintln(&input, x)
		}
		for i, row := range LeanRows(binary, len(words), &input) {
			x := words[i]
			got, want := actual(x), expected(x)
			if got != want || row != strconv.FormatUint(want, 10) {
				log.Fatalf("%s %s input %#x: Go=%d, Lean=%q, expected=%d", binary, label, x, got, row, want)
			}
		}
		fmt.Printf("%s %s: %d inputs agree.\n", binary, label, len(words))
	}

	edges := []uint64{0, mask}
	for i := range width {
		bit := uint64(1) << i
		edges = append(edges, bit-1, bit, bit+1, mask^bit)
	}
	compare("boundaries", edges)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		words := make([]uint64, 20_000)
		for i := range words {
			words[i] = rng.Uint64() & mask
		}
		compare(fmt.Sprintf("seed=%d", seed), words)
	}
}
