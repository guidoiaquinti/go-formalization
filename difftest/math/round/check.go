package main

import (
	"log"
	"math"
	"math/rand"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	edges := binary64.Boundaries()
	for exponent := uint64(1022); exponent <= 1075; exponent++ {
		midpoint := uint64(1 << 51)
		if exponent >= 1023 && exponent < 1075 {
			midpoint >>= exponent - 1023
		}
		lastTie := (1 << 52) - midpoint
		for _, sign := range []uint64{0, 1 << 63} {
			for _, fraction := range []uint64{0, midpoint - 1, midpoint, midpoint + 1, lastTie - 1, lastTie, 1<<52 - 1} {
				edges = append(edges, sign|exponent<<52|fraction)
			}
		}
	}
	check.CompareSamples("godiff-round", edges, func(rng *rand.Rand, i int) uint64 {
		word := rng.Uint64()
		if i%2 == 0 {
			word = word&^uint64(0x7ff<<52) | uint64(1022+rng.Intn(54))<<52
		}
		return word
	}, decimal, func(word uint64) string {
		return observe(math.Round(binary64.Store(word)))
	}, func(word uint64) string {
		numeric := binary64.Load(expected(binary64.Store(word)))
		if spec := expectedWord(word); spec != numeric {
			log.Fatalf("Round oracles disagree for %#016x: word specification=%#016x, numeric oracle=%#016x", word, spec, numeric)
		}
		return observe(binary64.Store(numeric))
	})
}

func observe(value float64) string {
	if math.IsNaN(value) {
		return "18446744073709551616"
	}
	return decimal(binary64.Load(value))
}

func expected(input float64) float64 {
	if math.IsNaN(input) || input <= -(1<<52) || input >= 1<<52 {
		return input
	}
	// The integer conversion truncates in this range. The difference is
	// exact, since the operands are within a factor of two or whole is zero.
	whole := int64(input)
	fraction := input - float64(whole)
	if fraction >= 0.5 {
		whole++
	} else if fraction <= -0.5 {
		whole--
	}
	if whole == 0 {
		return binary64.Store(binary64.Load(input) & (1 << 63))
	}
	return float64(whole)
}

// expectedWord evaluates the exact quotient-and-residue theorem.
func expectedWord(word uint64) uint64 {
	exponent := word >> 52 & 0x7ff
	if exponent < 1023 {
		sign := word & (1 << 63)
		if exponent == 1022 {
			return sign + (1023 << 52)
		}
		return sign
	}
	if exponent >= 1075 {
		return word
	}
	quantum := uint64(1) << (1075 - exponent)
	whole := word / quantum
	if word%quantum >= quantum/2 {
		whole++
	}
	return whole * quantum
}

func decimal(word uint64) string { return strconv.FormatUint(word, 10) }
