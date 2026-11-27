package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type words struct {
	magnitude, sign uint64
}

func main() {
	magnitudes := []uint64{
		0, 1, 0x000fffffffffffff, 0x0010000000000000,
		0x3ff0000000000000, 0x7fefffffffffffff, 0x7ff0000000000000,
		0x7ff0000000000001, 0x7ff7ffffffffffff, 0x7ff8000000000000,
		0x7ff8000000000001, 0x7fffffffffffffff,
	}
	var bits []uint64
	for _, sign := range []uint64{0, 1 << 63} {
		for _, magnitude := range magnitudes {
			bits = append(bits, magnitude+sign)
		}
	}
	var edges []words
	for _, magnitude := range bits {
		for _, sign := range bits {
			edges = append(edges, words{magnitude, sign})
		}
	}
	check.CompareSamples("godiff-copysign", edges, func(rng *rand.Rand, _ int) words {
		return words{rng.Uint64(), rng.Uint64()}
	}, func(w words) string {
		return fmt.Sprintf("%d %d", w.magnitude, w.sign)
	}, func(w words) string {
		got := math.Copysign(math.Float64frombits(w.magnitude), math.Float64frombits(w.sign))
		return strconv.FormatUint(math.Float64bits(got), 10)
	}, func(w words) string {
		const signWeight = 1 << 63
		return strconv.FormatUint(w.magnitude%signWeight+(w.sign/signWeight)*signWeight, 10)
	})
}
