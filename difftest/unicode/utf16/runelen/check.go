package main

import (
	"math"
	"math/rand"
	"strconv"
	"unicode/utf16"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	edges := []rune{math.MinInt32, -1, 0x10000, 0x10001, 0x10ffff, 0x110000, math.MaxInt32}
	for r := rune(0); r < 0x10000; r++ {
		edges = append(edges, r)
	}
	check.CompareSamples("godiff-utf16-runelen", edges, func(rng *rand.Rand, i int) rune {
		if i%2 == 0 {
			return 0x10000 + rng.Int31n(0x100000)
		}
		return rune(rng.Uint32())
	}, func(r rune) string {
		return strconv.FormatInt(int64(r), 10)
	}, func(r rune) string {
		return strconv.Itoa(utf16.RuneLen(r))
	}, expected)
}

func expected(r rune) string {
	plane, block := r>>16, r>>10
	if r < 0 || plane >= 17 || block == 54 || block == 55 {
		return "-1"
	}
	if plane == 0 {
		return "1"
	}
	return "2"
}
