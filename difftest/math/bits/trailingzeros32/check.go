package main

import (
	"math/bits"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const wordWidth = 32

func expected(x uint64) uint64 {
	for i := range uint64(wordWidth) {
		if x&(1<<i) != 0 {
			return i
		}
	}
	return wordWidth
}

func main() {
	check.SampleUnary("godiff-trailingzeros32", wordWidth, func(x uint64) uint64 {
		return uint64(bits.TrailingZeros32(uint32(x)))
	}, expected)
}
