package main

import (
	"math"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	check.SampleUnary("godiff-signbit", 64, func(repr uint64) uint64 {
		if math.Signbit(math.Float64frombits(repr)) {
			return 1
		}
		return 0
	}, func(repr uint64) uint64 {
		return repr >> 63
	})
}
