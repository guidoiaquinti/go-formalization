package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type pair struct{ a, b uint64 }

const (
	sign = uint64(1 << 63)
	inf  = uint64(0x7ff0000000000000)
)

func main() {
	words := binary64.Boundaries()
	var edges []pair
	for _, a := range words {
		for _, b := range words {
			edges = append(edges, pair{a, b})
		}
	}
	check.CompareSamples("godiff-min", edges, sample, func(p pair) string {
		return fmt.Sprintf("%d %d", p.a, p.b)
	}, func(p pair) string {
		got := binary64.Load(math.Min(binary64.Store(p.a), binary64.Store(p.b)))
		if got%sign > inf {
			return "nan"
		}
		return strconv.FormatUint(got, 10)
	}, expected)
}

func sample(rng *rand.Rand, i int) pair {
	a, b := rng.Uint64(), rng.Uint64()
	switch i % 5 {
	case 1:
		b = a
	case 2:
		b = a ^ sign
	case 3:
		a = sign | inf
	case 4:
		b |= inf | 1 // Force a NaN, in either sign.
	}
	return pair{a, b}
}

// Key orders finite values and infinities, with negative zero first among zeros.
func key(word uint64) uint64 {
	if word&sign != 0 {
		return ^word
	}
	return word | sign
}

func expected(p pair) string {
	switch {
	case p.a == sign|inf || p.b == sign|inf:
		return strconv.FormatUint(sign|inf, 10)
	case p.a%sign > inf || p.b%sign > inf:
		return "nan"
	case key(p.a) <= key(p.b):
		return strconv.FormatUint(p.a, 10)
	default:
		return strconv.FormatUint(p.b, 10)
	}
}
