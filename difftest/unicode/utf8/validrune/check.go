package main

import (
	"math"
	"math/rand"
	"strconv"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	edges := []rune{math.MinInt32, -1, 0x10000, 0x10fffe, 0x10ffff, 0x110000, math.MaxInt32}
	for r := rune(0); r < 0x10000; r++ {
		edges = append(edges, r)
	}
	check.CompareSamples("godiff-utf8-validrune", edges, func(rng *rand.Rand, i int) rune {
		if i%2 == 0 {
			return rune(rng.Uint32())
		}
		return rng.Int31n(0x110000)
	}, func(r rune) string {
		return strconv.FormatInt(int64(r), 10)
	}, func(r rune) string {
		return strconv.FormatBool(utf8.ValidRune(r))
	}, func(r rune) string {
		// A Unicode scalar has plane index 0 through 16 and no surrogate prefix.
		plane, prefix := r>>16, r>>10
		return strconv.FormatBool(0 <= plane && plane <= 16 && prefix != 54 && prefix != 55)
	})
}
