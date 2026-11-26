package main

import (
	"math/bits"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	check.Exhaustive("godiff-trailingzeros16", 1<<16, func(x int) (int, int) {
		count := 0
		for i := range 16 {
			if x&(1<<i) != 0 {
				break
			}
			count++
		}
		return int(bits.TrailingZeros16(uint16(x))), count
	})
}
