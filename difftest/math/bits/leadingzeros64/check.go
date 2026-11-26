package main

import (
	"math/bits"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const wordWidth = 64

func expected(x uint64) uint64 {
	count := uint64(0)
	for mask := uint64(1) << (wordWidth - 1); mask != 0 && x&mask == 0; mask >>= 1 {
		count++
	}
	return count
}

func main() {
	check.SampleUnary("godiff-leadingzeros64", wordWidth, func(x uint64) uint64 {
		return uint64(bits.LeadingZeros64(uint64(x)))
	}, expected)
}
