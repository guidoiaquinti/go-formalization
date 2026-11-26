package main

import (
	"log"
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
	if bits.UintSize != wordWidth {
		log.Fatal("LeadingZeros comparison requires 64-bit uint")
	}
	check.SampleUnary("godiff-leadingzeros", wordWidth, func(x uint64) uint64 {
		return uint64(bits.LeadingZeros(uint(x)))
	}, expected)
}
