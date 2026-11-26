package main

import (
	"math/bits"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	check.Exhaustive("godiff-trailingzeros8", 1<<8, func(x int) (int, int) {
		zeros := 0
		for zeros < 8 && x&(1<<zeros) == 0 {
			zeros++
		}
		return bits.TrailingZeros8(uint8(x)), zeros
	})
}
