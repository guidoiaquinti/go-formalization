package main

import (
	"math/bits"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	check.Exhaustive("godiff-leadingzeros16", 1<<16, func(x int) (int, int) {
		count := 0
		for mask := 1 << 15; mask != 0 && x&mask == 0; mask >>= 1 {
			count++
		}
		return int(bits.LeadingZeros16(uint16(x))), count
	})
}
