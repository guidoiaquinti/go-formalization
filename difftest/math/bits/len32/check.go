package main

import (
	"math/bits"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const wordWidth = 32

func expected(x uint64) uint64 {
	width := uint64(0)
	for x != 0 {
		x /= 2
		width++
	}
	return width
}

func main() {
	check.SampleUnary("godiff-len32", wordWidth, func(x uint64) uint64 {
		return uint64(bits.Len32(uint32(x)))
	}, expected)
}
