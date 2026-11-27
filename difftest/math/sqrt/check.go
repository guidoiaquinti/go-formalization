package main

import (
	"math"
	"math/big"
	"math/rand"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	edges := binary64.Boundaries()
	for _, sign := range []uint64{0, 1 << 63} {
		for exponent := uint64(1); exponent < 2047; exponent++ {
			power := exponent << 52
			edges = append(edges, sign|(power-1), sign|power, sign|(power+1))
		}
		for k := range 52 {
			power := uint64(1) << k
			edges = append(edges, sign|(power-1), sign|power, sign|(power+1))
		}
	}
	check.CompareSamples("godiff-math-sqrt", edges, sample,
		func(word uint64) string { return strconv.FormatUint(word, 10) },
		func(word uint64) string { return observe(math.Sqrt(binary64.Store(word))) },
		expected)
}

func sample(rng *rand.Rand, i int) uint64 {
	word := rng.Uint64()
	switch i % 4 {
	case 0:
		word &= (1 << 63) - 1
	case 1:
		word &= (1 << 52) - 1
	case 2:
		integer := float64(uint32(word)) + 1
		word = binary64.Load(integer*integer) + uint64(i%3) - 1
	}
	return word
}

func expected(word uint64) string {
	x := binary64.Store(word)
	switch {
	case x == 0 || math.IsInf(x, 1):
		return observe(x)
	case x < 0 || math.IsNaN(x):
		return "NaN"
	}
	// An exact binary64 input has at most 53 significand bits. Even its
	// smallest positive square root is normal, so 53-bit rounding suffices.
	input := new(big.Float).SetFloat64(x)
	root := new(big.Float).SetPrec(53).SetMode(big.ToNearestEven).Sqrt(input)
	rounded, _ := root.Float64()
	return observe(rounded)
}

func observe(value float64) string {
	if math.IsNaN(value) {
		return "NaN"
	}
	return strconv.FormatUint(binary64.Load(value), 10)
}
