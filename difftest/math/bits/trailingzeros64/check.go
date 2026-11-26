package main

import (
	"math/bits"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const wordWidth = 64

func expected(x uint64) uint64 {
	for i := range uint64(wordWidth) {
		if x&(1<<i) != 0 {
			return i
		}
	}
	return wordWidth
}

func main() {
	check.SampleUnary("godiff-trailingzeros64", wordWidth, func(x uint64) uint64 {
		return uint64(bits.TrailingZeros64(uint64(x)))
	}, expected)
}
