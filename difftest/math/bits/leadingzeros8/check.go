package main

import (
	"math/bits"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	check.Exhaustive("godiff-leadingzeros8", 1<<8, func(x int) (int, int) {
		// Count zeros from the most significant bit until the first set bit.
		count := 0
		for mask := 1 << 7; mask != 0 && x&mask == 0; mask >>= 1 {
			count++
		}
		return bits.LeadingZeros8(uint8(x)), count
	})
}
