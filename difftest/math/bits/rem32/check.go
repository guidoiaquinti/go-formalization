package main

import (
	"fmt"
	"math/big"
	"math/bits"
	"math/rand"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const wordWidth = 32

type operands struct{ hi, lo, y uint64 }

func input(c operands) string { return fmt.Sprintf("%d %d %d", c.hi, c.lo, c.y) }

func actual(c operands) (result string) {
	defer func() {
		if p := recover(); p != nil {
			result = check.FormatDivisionPanic(p)
		}
	}()
	return fmt.Sprint(bits.Rem32(uint32(c.hi), uint32(c.lo), uint32(c.y)))
}

func expected(c operands) string {
	if c.y == 0 {
		return "divide-by-zero"
	}
	n := new(big.Int).SetUint64(c.hi)
	n.Lsh(n, wordWidth).Add(n, new(big.Int).SetUint64(c.lo))
	y := new(big.Int).SetUint64(c.y)
	return n.Mod(n, y).String()
}

func main() {
	words := check.WordBoundaries(wordWidth)
	var edges []operands
	for _, hi := range words {
		for _, lo := range words {
			for _, y := range words {
				edges = append(edges, operands{hi, lo, y})
			}
		}
	}
	mask := ^uint64(0) >> (64 - wordWidth)
	check.CompareSamples("godiff-rem32", edges, func(rng *rand.Rand, i int) operands {
		c := operands{rng.Uint64() & mask, rng.Uint64() & mask, rng.Uint64() & mask}
		if offset := i % 8; offset < 3 {
			c.hi = (c.y + uint64(offset) - 1) & mask
		}
		return c
	}, input, actual, expected)
}
