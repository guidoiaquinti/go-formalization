package main

import (
	"math/bits"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	check.Exhaustive("godiff-len8", 1<<8, func(x int) (int, int) {
		// Find the smallest width whose unsigned range contains x.
		width := 0
		for x >= 1<<width {
			width++
		}
		return bits.Len8(uint8(x)), width
	})
}
