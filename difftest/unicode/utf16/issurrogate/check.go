package main

import (
	"math"
	"math/rand"
	"strconv"
	"unicode/utf16"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	edges := []rune{math.MinInt32, -0x2800, -1, 0x10000, 0x10ffff, 0x110000, math.MaxInt32}
	// Include every BMP code point, especially both complete surrogate blocks.
	for r := rune(0); r < 0x10000; r++ {
		edges = append(edges, r)
	}
	check.CompareSamples("godiff-utf16-issurrogate", edges, func(rng *rand.Rand, _ int) rune {
		return rune(rng.Uint32())
	}, func(r rune) string {
		return strconv.FormatInt(int64(r), 10)
	}, func(r rune) string {
		return strconv.FormatBool(utf16.IsSurrogate(r))
	}, func(r rune) string {
		// The two surrogate blocks have the signed high-bit prefixes 54 and 55.
		prefix := r >> 10
		return strconv.FormatBool(prefix == 54 || prefix == 55)
	})
}
