package main

import (
	"log"
	"math/bits"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const wordWidth = 64

func expected(x uint64) uint64 {
	width := uint64(0)
	for x != 0 {
		x /= 2
		width++
	}
	return width
}

func main() {
	if bits.UintSize != wordWidth {
		log.Fatal("Len comparison requires 64-bit uint")
	}
	check.SampleUnary("godiff-len", wordWidth, func(x uint64) uint64 {
		return uint64(bits.Len(uint(x)))
	}, expected)
}
