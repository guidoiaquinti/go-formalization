package main

import (
	"log"
	"math/bits"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const wordWidth = 64

func main() {
	if bits.UintSize != wordWidth {
		log.Fatal("OnesCount comparison requires 64-bit uint")
	}
	check.SampleUnary("godiff-onescount", wordWidth, func(x uint64) uint64 {
		return uint64(bits.OnesCount(uint(x)))
	}, func(x uint64) uint64 {
		var count uint64
		for i := range wordWidth {
			count += (x >> i) & 1
		}
		return count
	})
}
