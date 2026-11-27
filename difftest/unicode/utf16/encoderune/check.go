package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"unicode/utf16"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	edges := []rune{
		math.MinInt32, -1, 0, 0xd7ff, 0xd800, 0xdfff, 0xe000, 0xffff,
		0x10000, 0x10001, 0x103ff, 0x10400, 0x10401,
		0x10fbff, 0x10fc00, 0x10ffff, 0x110000, math.MaxInt32,
	}
	check.CompareSamples("godiff-utf16-encoderune", edges, func(rng *rand.Rand, i int) rune {
		if i%2 == 0 {
			return 0x10000 + rng.Int31n(0x100000)
		}
		return rune(rng.Uint32())
	}, func(r rune) string {
		return strconv.FormatInt(int64(r), 10)
	}, func(r rune) string {
		high, low := utf16.EncodeRune(r)
		return fmt.Sprintf("%d %d", high, low)
	}, expected)
}

func expected(r rune) string {
	plane := r >> 16
	if plane < 1 || plane > 16 {
		return "65533 65533"
	}
	payload := r - 65536
	return fmt.Sprintf("%d %d", 55296+payload/1024, 56320+payload%1024)
}
