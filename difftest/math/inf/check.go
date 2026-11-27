package main

import (
	"math"
	"math/rand"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	edges := []int{math.MinInt, math.MinInt + 1, -2, -1, 0, 1, 2, math.MaxInt - 1, math.MaxInt}
	check.CompareSamples("godiff-inf", edges, func(rng *rand.Rand, _ int) int {
		return int(rng.Uint64())
	}, strconv.Itoa, func(sign int) string {
		return strconv.FormatUint(math.Float64bits(math.Inf(sign)), 10)
	}, expected)
}

func expected(sign int) string {
	const exponent uint64 = 2047 << 52
	var signBit uint64
	if sign < 0 {
		signBit = 1 << 63
	}
	return strconv.FormatUint(signBit+exponent, 10)
}
