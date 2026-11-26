package main

import (
	"math/bits"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const wordWidth = 64

func expected(x uint64) uint64 {
	result := uint64(0)
	for i := range wordWidth / 8 {
		b := (x >> (8 * i)) & 0xff
		result |= b << (wordWidth - 8*(i+1))
	}
	return result
}

func main() {
	check.SampleUnary("godiff-reversebytes64", wordWidth, func(x uint64) uint64 {
		return uint64(bits.ReverseBytes64(uint64(x)))
	}, expected)
}
