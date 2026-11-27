package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"unicode/utf16"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type pair struct {
	first, second rune
}

func main() {
	boundaries := []rune{math.MinInt32, -1, 0, 0xd7ff, 0xd800, 0xdbff, 0xdc00, 0xdfff, 0xe000, 0x10ffff, math.MaxInt32}
	var edges []pair
	for _, first := range boundaries {
		for _, second := range boundaries {
			edges = append(edges, pair{first, second})
		}
	}
	check.CompareSamples("godiff-utf16-decoderune", edges, sample, func(p pair) string {
		return fmt.Sprintf("%d %d", p.first, p.second)
	}, func(p pair) string {
		return strconv.Itoa(int(utf16.DecodeRune(p.first, p.second)))
	}, expected)
}

func sample(rng *rand.Rand, index int) pair {
	high := 0xd800 + rng.Int31n(1024)
	low := 0xdc00 + rng.Int31n(1024)
	switch index % 4 {
	case 0:
		return pair{high, low}
	case 1:
		return pair{high, rune(rng.Uint32())}
	case 2:
		return pair{rune(rng.Uint32()), low}
	default:
		return pair{rune(rng.Uint32()), rune(rng.Uint32())}
	}
}

func expected(p pair) string {
	if p.first>>10 != 54 || p.second>>10 != 55 {
		return "65533"
	}
	first := int(p.first) - 0xd800
	second := int(p.second) - 0xdc00
	return strconv.Itoa(65536 + 1024*first + second)
}
