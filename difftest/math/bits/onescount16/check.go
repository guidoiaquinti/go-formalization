package main

import (
	"math/bits"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	check.Exhaustive("godiff-onescount16", 1<<16, func(x int) (int, int) {
		count := 0
		for i := range 16 {
			count += (x >> i) & 1
		}
		return int(bits.OnesCount16(uint16(x))), count
	})
}
