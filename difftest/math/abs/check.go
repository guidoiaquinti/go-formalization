package main

import (
	"math"
	"math/rand"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	words := []uint64{
		0, 1, 0x000fffffffffffff, 0x0010000000000000,
		0x3ff0000000000000, 0x7fefffffffffffff, 0x7ff0000000000000,
		0x7ff0000000000001, 0x7ff7ffffffffffff, 0x7ff8000000000000,
		0x7ff8000000000001, 0x7fffffffffffffff,
	}
	var edges []uint64
	for _, sign := range []uint64{0, 1 << 63} {
		for _, word := range words {
			edges = append(edges, word+sign)
		}
	}
	check.CompareSamples("godiff-abs", edges, func(rng *rand.Rand, _ int) uint64 {
		return rng.Uint64()
	}, decimal, func(bits uint64) string {
		return decimal(math.Float64bits(math.Abs(math.Float64frombits(bits))))
	}, func(bits uint64) string {
		return decimal(bits % (1 << 63))
	})
}

func decimal(bits uint64) string {
	return strconv.FormatUint(bits, 10)
}
