package main

import (
	"math/bits"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	check.Exhaustive("godiff-onescount8", 1<<8, func(x int) (int, int) {
		ones := 0
		for i := range 8 {
			ones += (x >> i) & 1
		}
		return bits.OnesCount8(uint8(x)), ones
	})
}
