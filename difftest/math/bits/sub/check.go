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

type operands struct{ x, y, flag uint64 }

func input(c operands) string { return fmt.Sprintf("%d %d %d", c.x, c.y, c.flag) }

func actual(c operands) string {
	first, second := bits.Sub(uint(c.x), uint(c.y), uint(c.flag))
	return fmt.Sprintf("%d %d", first, second)
}

func expected(c operands) string {
	radix := new(big.Int).Lsh(big.NewInt(1), wordWidth)
	n := new(big.Int).SetUint64(c.x)
	n.Sub(n, new(big.Int).SetUint64(c.y))
	n.Sub(n, new(big.Int).SetUint64(c.flag))
	flag := 0
	if n.Sign() < 0 {
		flag = 1
	}
	n.Mod(n, radix)
	return fmt.Sprintf("%d %d", n.Uint64(), flag)
}

func main() {
	if bits.UintSize != wordWidth {
		log.Fatal("Sub comparison requires 64-bit uint")
	}
	words := check.WordBoundaries(wordWidth)
	var edges []operands
	for _, x := range words {
		for _, y := range words {
			edges = append(edges, operands{x, y, 0}, operands{x, y, 1})
		}
	}
	mask := ^uint64(0) >> (64 - wordWidth)
	check.CompareSamples("godiff-sub", edges, func(rng *rand.Rand, i int) operands {
		x := rng.Uint64() & mask
		y := rng.Uint64() & mask
		if i%4 == 0 {
			y = x
		}
		return operands{x, y, uint64(i%2 ^ i/4%2)}
	}, input, actual, expected)
}
