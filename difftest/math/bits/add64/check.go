package main

import (
	"fmt"
	"math/big"
	"math/bits"
	"math/rand"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const wordWidth = 64

type operands struct{ x, y, flag uint64 }

func input(c operands) string { return fmt.Sprintf("%d %d %d", c.x, c.y, c.flag) }

func actual(c operands) string {
	first, second := bits.Add64(uint64(c.x), uint64(c.y), uint64(c.flag))
	return fmt.Sprintf("%d %d", first, second)
}

func expected(c operands) string {
	radix := new(big.Int).Lsh(big.NewInt(1), wordWidth)
	n := new(big.Int).SetUint64(c.x)
	n.Add(n, new(big.Int).SetUint64(c.y))
	n.Add(n, new(big.Int).SetUint64(c.flag))
	flag := new(big.Int).Quo(new(big.Int).Set(n), radix)
	n.Mod(n, radix)
	return fmt.Sprintf("%d %d", n.Uint64(), flag.Uint64())
}

func main() {
	words := check.WordBoundaries(wordWidth)
	var edges []operands
	for _, x := range words {
		for _, y := range words {
			edges = append(edges, operands{x, y, 0}, operands{x, y, 1})
		}
	}
	mask := ^uint64(0) >> (64 - wordWidth)
	check.CompareSamples("godiff-add64", edges, func(rng *rand.Rand, i int) operands {
		x := rng.Uint64() & mask
		y := rng.Uint64() & mask
		if i%4 == 0 {
			y = mask ^ x
		}
		return operands{x, y, uint64(i%2 ^ i/4%2)}
	}, input, actual, expected)
}
