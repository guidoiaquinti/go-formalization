package main

import (
	"math/bits"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const wordWidth = 32

func expected(x uint64) uint64 {
	result := uint64(0)
	for i := range wordWidth {
		result |= ((x >> i) & 1) << (wordWidth - 1 - i)
	}
	return result
}

func main() {
	check.SampleUnary("godiff-reverse32", wordWidth, func(x uint64) uint64 {
		return uint64(bits.Reverse32(uint32(x)))
	}, expected)
}
