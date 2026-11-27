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
	for exponent := uint64(1022); exponent <= 1075; exponent++ {
		midpoint := uint64(1 << 51)
		if exponent >= 1023 && exponent < 1075 {
			midpoint >>= exponent - 1023
		}
		for _, sign := range []uint64{0, 1 << 63} {
			for _, payload := range []uint64{0, 1, midpoint - 1, midpoint, midpoint + 1, 1<<52 - 1} {
				edges = append(edges, sign|exponent<<52|payload)
			}
		}
	}
	check.CompareSamples("godiff-math-modf", edges, func(rng *rand.Rand, i int) uint64 {
		word := rng.Uint64()
		if i%2 == 0 {
			word = word&^uint64(2047<<52) | uint64(1022+rng.Intn(54))<<52
		}
		return word
	}, func(word uint64) string {
		return strconv.FormatUint(word, 10)
	}, func(word uint64) string {
		whole, fraction := math.Modf(binary64.Store(word))
		return observe(whole) + " " + observe(fraction)
	}, expected)
}

func expected(word uint64) string {
	x := binary64.Store(word)
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return observe(x) + " NaN"
	}
	exact := new(big.Rat).SetFloat64(x)
	var integer big.Int
	integer.Quo(exact.Num(), exact.Denom())
	exact.Sub(exact, new(big.Rat).SetInt(&integer))
	whole, _ := integer.Float64()
	fraction, _ := exact.Float64()
	return observe(math.Copysign(whole, x)) + " " + observe(math.Copysign(fraction, x))
}

func observe(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(binary64.Load(x), 10)
}
