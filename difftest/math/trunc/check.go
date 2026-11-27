package main

import (
	"math"
	"math/rand"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const nan = "18446744073709551616"

func main() {
	edges := binary64.Boundaries()
	for _, sign := range []uint64{0, 1 << 63} {
		for _, exponent := range []uint64{1022, 1023, 1024, 1073, 1074, 1075, 1076} {
			for _, fraction := range []uint64{0, 1, (1 << 51) - 1, 1 << 51, (1 << 52) - 1} {
				edges = append(edges, sign|exponent<<52|fraction)
			}
		}
	}
	check.CompareSamples("godiff-trunc", edges, func(rng *rand.Rand, i int) uint64 {
		word := rng.Uint64()
		if i%2 == 0 {
			// Exercise every exponent that can contain fractional and integral bits.
			word = word&^uint64(0x7ff<<52) | uint64(1022+rng.Intn(55))<<52
		}
		return word
	}, decimal, func(word uint64) string {
		return observe(binary64.Load(math.Trunc(binary64.Store(word))))
	}, func(word uint64) string { return observe(expectedWord(word)) })
}

func observe(word uint64) string {
	if math.IsNaN(binary64.Store(word)) {
		return nan
	}
	return decimal(word)
}

func expectedWord(word uint64) uint64 {
	input := binary64.Store(word)
	if math.IsNaN(input) || input <= -(1<<52) || input >= 1<<52 {
		return word
	}
	// The conversion truncates toward zero in this range and its integer
	// result is exactly representable. Keep the sign when it becomes zero.
	whole := int64(input)
	if whole == 0 {
		return word & (1 << 63)
	}
	return binary64.Load(float64(whole))
}

func decimal(word uint64) string { return strconv.FormatUint(word, 10) }
