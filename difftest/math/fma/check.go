package main

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct{ x, y, z uint64 }

func main() {
	words := append(binary64.Boundaries(), 0x3fe0000000000000, 0x3ff0000000000001,
		0x4000000000000000, 0x7feffffffffffffe, 0xbff0000000000001, 2)
	var edges []input
	for _, x := range words {
		for _, y := range words {
			for _, z := range words {
				edges = append(edges, input{x, y, z})
			}
		}
	}
	edges = append(edges,
		input{0x3ff0000000000000, 0x3ff0000000000000, 0x3ca0000000000000}, // ties down to 1
		input{0x3ff0000000000000, 0x3ff0000000000001, 0x3ca0000000000000}, // ties up to the even successor
		input{0x7fefffffffffffff, 0x4000000000000000, 0xffefffffffffffff}, // overflow avoided by cancellation
		input{0x8000000000000001, 0x3fe0000000000000, 0},                  // negative underflow
		input{0x3ff0000000000001, 0x3fefffffffffffff, 0xbff0000000000000}, // exact fused residual
	)
	check.CompareSamples("godiff-math-fma", edges, sample, func(c input) string {
		return fmt.Sprintf("%d %d %d", c.x, c.y, c.z)
	}, func(c input) string {
		return observe(math.FMA(binary64.Store(c.x), binary64.Store(c.y), binary64.Store(c.z)))
	}, expected)
}

func sample(r *rand.Rand, i int) input {
	x, y, z := r.Uint64(), r.Uint64(), r.Uint64()
	switch i % 5 {
	case 0:
		// Cancel the rounded product, leaving its exact residual.
		product := binary64.Store(x) * binary64.Store(y)
		z = binary64.Load(-product)
	case 1:
		// Choose normal powers of two whose product is near underflow.
		exponent := r.Intn(1023) - 1022
		x = binary64.Load(math.Ldexp(1, exponent)) | x&(1<<63)
		y = binary64.Load(math.Ldexp(1, -1074-exponent+r.Intn(8)-4)) | y&(1<<63)
		z &= 1<<52 - 1
	case 2:
		// Exact cancellation, closely aligned sums, and ties.
		x = x&((1<<52)-1) | uint64(1023)<<52
		y = 0x3ff0000000000000
		z = (x ^ (1 << 63)) + uint64(r.Intn(5)-2)
	case 3:
		z &= 1 << 63
	}
	return input{x, y, z}
}

func expected(c input) string {
	x, y, z := binary64.Store(c.x), binary64.Store(c.y), binary64.Store(c.z)
	negativeProduct := math.Signbit(x) != math.Signbit(y)
	switch {
	case math.IsNaN(x) || math.IsNaN(y) || math.IsNaN(z):
		return "NaN"
	case math.IsInf(x, 0) || math.IsInf(y, 0):
		if x == 0 || y == 0 || math.IsInf(z, 0) && math.Signbit(z) != negativeProduct {
			return "NaN"
		}
		if negativeProduct {
			return observe(math.Inf(-1))
		}
		return observe(math.Inf(1))
	case math.IsInf(z, 0):
		return observe(z)
	}
	a := new(big.Rat).SetFloat64(x)
	b := new(big.Rat).SetFloat64(y)
	addend := new(big.Rat).SetFloat64(z)
	result := a.Add(a.Mul(a, b), addend)
	if result.Sign() == 0 && z == 0 && (x == 0 || y == 0) && negativeProduct && math.Signbit(z) {
		return observe(math.Copysign(0, -1))
	}
	rounded, _ := result.Float64()
	return observe(rounded)
}

func observe(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(binary64.Load(x), 10)
}
