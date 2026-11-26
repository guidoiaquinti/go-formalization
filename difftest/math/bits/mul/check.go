package main

import (
	"fmt"
	"log"
	"math/big"
	"math/bits"
	"math/rand"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const wordWidth = 64

type operands struct{ x, y uint64 }

func input(c operands) string { return fmt.Sprintf("%d %d", c.x, c.y) }

func actual(c operands) string {
	first, second := bits.Mul(uint(c.x), uint(c.y))
	return fmt.Sprintf("%d %d", first, second)
}

func expected(c operands) string {
	radix := new(big.Int).Lsh(big.NewInt(1), wordWidth)
	n := new(big.Int).SetUint64(c.x)
	n.Mul(n, new(big.Int).SetUint64(c.y))
	high := new(big.Int).Rsh(new(big.Int).Set(n), wordWidth)
	n.Mod(n, radix)
	return fmt.Sprintf("%d %d", high.Uint64(), n.Uint64())
}

func main() {
	if bits.UintSize != wordWidth {
		log.Fatal("Mul comparison requires 64-bit uint")
	}
	words := check.WordBoundaries(wordWidth)
	var edges []operands
	for _, x := range words {
		for _, y := range words {
			edges = append(edges, operands{x, y})
		}
	}
	mask := ^uint64(0) >> (64 - wordWidth)
	check.CompareSamples("godiff-mul", edges, func(rng *rand.Rand, i int) operands {
		x := rng.Uint64() & mask
		y := rng.Uint64() & mask
		return operands{x, y}
	}, input, actual, expected)
}
