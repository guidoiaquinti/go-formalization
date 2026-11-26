package main

import (
	"math/bits"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	check.Exhaustive("godiff-reverse16", 1<<16, func(x int) (int, int) {
		result := 0
		for i := range 16 {
			result |= ((x >> i) & 1) << (16 - 1 - i)
		}
		return int(bits.Reverse16(uint16(x))), result
	})
}
