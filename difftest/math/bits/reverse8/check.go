package main

import (
	"math/bits"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	check.Exhaustive("godiff-reverse8", 1<<8, func(x int) (int, int) {
		result := 0
		for i := range 8 {
			result |= ((x >> i) & 1) << (8 - 1 - i)
		}
		return int(bits.Reverse8(uint8(x))), result
	})
}
