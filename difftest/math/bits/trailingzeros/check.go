package main

import (
	"log"
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
	if bits.UintSize != wordWidth {
		log.Fatal("TrailingZeros comparison requires 64-bit uint")
	}
	check.SampleUnary("godiff-trailingzeros", wordWidth, func(x uint64) uint64 {
		return uint64(bits.TrailingZeros(uint(x)))
	}, expected)
}
