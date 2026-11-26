package main

import (
	"log"
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
	if bits.UintSize != wordWidth {
		log.Fatal("ReverseBytes comparison requires 64-bit uint")
	}
	check.SampleUnary("godiff-reversebytes", wordWidth, func(x uint64) uint64 {
		return uint64(bits.ReverseBytes(uint(x)))
	}, expected)
}
