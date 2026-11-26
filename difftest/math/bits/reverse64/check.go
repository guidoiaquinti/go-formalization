package main

import (
	"math/bits"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const wordWidth = 64

func expected(x uint64) uint64 {
	result := uint64(0)
	for i := range wordWidth {
		result |= ((x >> i) & 1) << (wordWidth - 1 - i)
	}
	return result
}

func main() {
	check.SampleUnary("godiff-reverse64", wordWidth, func(x uint64) uint64 {
		return uint64(bits.Reverse64(uint64(x)))
	}, expected)
}
