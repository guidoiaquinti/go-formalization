package main

import (
	"math/bits"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	check.Exhaustive("godiff-len16", 1<<16, func(x int) (int, int) {
		width := 0
		for x >= 1<<width {
			width++
		}
		return int(bits.Len16(uint16(x))), width
	})
}
