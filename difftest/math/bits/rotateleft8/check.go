package main

import (
	"fmt"
	"math"
	"math/bits"
	"math/rand"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const wordWidth = 8

type rotation struct {
	x uint64
	k int
}

func input(c rotation) string { return fmt.Sprintf("%d %d", c.x, c.k) }

func actual(c rotation) string {
	return fmt.Sprint(bits.RotateLeft8(uint8(c.x), c.k))
}

func expected(c rotation) string {
	shift := c.k % wordWidth
	if shift < 0 {
		shift += wordWidth
	}
	var result uint64
	for i := range wordWidth {
		result |= ((c.x >> i) & 1) << ((i + shift) % wordWidth)
	}
	return fmt.Sprint(result)
}

func main() {
	words := check.WordBoundaries(wordWidth)
	if wordWidth == 8 {
		words = make([]uint64, 256)
		for i := range words {
			words[i] = uint64(i)
		}
	}
	var edges []rotation
	for _, x := range words {
		for k := -wordWidth; k <= wordWidth; k++ {
			edges = append(edges, rotation{x, k})
		}
		for _, k := range []int{math.MinInt, math.MinInt + 1, math.MaxInt - 1, math.MaxInt} {
			edges = append(edges, rotation{x, k})
		}
	}
	mask := ^uint64(0) >> (64 - wordWidth)
	check.CompareSamples("godiff-rotateleft8", edges, func(rng *rand.Rand, _ int) rotation {
		return rotation{rng.Uint64() & mask, int(rng.Uint64())}
	}, input, actual, expected)
}
