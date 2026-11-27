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
	edges := []int{math.MinInt, math.MinInt + 1, -325, -324, 309, 310, math.MaxInt - 1, math.MaxInt}
	for n := -323; n <= 308; n++ {
		edges = append(edges, n)
	}

	// Derive every potentially nonconstant result once from exact decimals.
	var oracle [632]string
	for n := -323; n <= 308; n++ {
		oracle[n+323] = decimal(expected(n))
	}
	check.CompareSamples("godiff-math-pow10", edges, func(rng *rand.Rand, i int) int {
		if i%2 == 0 {
			return rng.Intn(len(oracle)) - 323
		}
		return int(rng.Uint64())
	}, strconv.Itoa, func(n int) string {
		return decimal(binary64.Load(math.Pow10(n)))
	}, func(n int) string {
		switch {
		case n > 308:
			return decimal(0x7ff0000000000000)
		case n < -323:
			return "0"
		default:
			return oracle[n+323]
		}
	})
}

func expected(n int) uint64 {
	magnitude := max(n, -n)
	block := magnitude / 32 * 32
	if n < 0 {
		block = -block
	}
	coarse := roundedPower(block)
	fine := roundedPower(magnitude % 32)
	if n < 0 {
		coarse.Quo(coarse, fine)
	} else {
		coarse.Mul(coarse, fine)
	}
	result, _ := coarse.Float64()
	return binary64.Load(result)
}

// Construct the ideal decimal operand, round to binary64 once, and
// re-read that rounded value as an exact dyadic rational.
func roundedPower(n int) *big.Rat {
	power := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(max(n, -n))), nil)
	exact := new(big.Rat).SetInt(power)
	if n < 0 {
		exact.Inv(exact)
	}
	rounded, _ := exact.Float64()
	return exact.SetFloat64(rounded)
}

func decimal(word uint64) string { return strconv.FormatUint(word, 10) }
