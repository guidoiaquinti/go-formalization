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

type input struct{ x, y uint64 }

func main() {
	words := append(binary64.Boundaries(), 0x3ff0000000000001, 0x3fefffffffffffff,
		0x3ca0000000000000, 0xbca0000000000000)
	var edges []input
	for _, x := range words {
		for _, y := range words {
			edges = append(edges, input{x, y})
		}
	}
	check.CompareSamples("godiff-math-dim", edges, sample,
		func(c input) string { return fmt.Sprintf("%d %d", c.x, c.y) },
		func(c input) string { return observe(math.Dim(binary64.Store(c.x), binary64.Store(c.y))) },
		expected)
}

func sample(rng *rand.Rand, i int) input {
	x, y := rng.Uint64(), rng.Uint64()
	switch i % 5 {
	case 1:
		y = x + 1
	case 2:
		y = x - 1
	case 3:
		x &= (1 << 63) | (1<<52 - 1)
		y &= (1 << 63) | (1<<52 - 1)
	case 4:
		y = x
	}
	return input{x, y}
}

func expected(c input) string {
	x, y := binary64.Store(c.x), binary64.Store(c.y)
	switch {
	case math.IsNaN(x) || math.IsNaN(y):
		return "NaN"
	case math.IsInf(x, 0) && x == y:
		return "NaN"
	case math.IsInf(x, 1) || math.IsInf(y, -1):
		return observe(math.Inf(1))
	case math.IsInf(x, -1) || math.IsInf(y, 1):
		return "0"
	}
	left, right := new(big.Rat).SetFloat64(x), new(big.Rat).SetFloat64(y)
	left.Sub(left, right)
	if left.Sign() <= 0 {
		return "0"
	}
	rounded, _ := left.Float64()
	return observe(rounded)
}

func observe(value float64) string {
	if math.IsNaN(value) {
		return "NaN"
	}
	return strconv.FormatUint(binary64.Load(value), 10)
}
