package main

import (
	"math/bits"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const wordWidth = 64

func expected(x uint64) uint64 {
	count := uint64(0)
	for i := range wordWidth {
		count += (x >> i) & 1
	}
	return count
}

func main() {
	check.SampleUnary("godiff-onescount64", wordWidth, func(x uint64) uint64 {
		return uint64(bits.OnesCount64(uint64(x)))
	}, expected)
}
