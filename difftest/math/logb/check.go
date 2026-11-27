package main

import (
	"math"
	"math/bits"
	"math/rand"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	edges := binary64.Boundaries()
	for i := range 52 {
		high := uint64(1) << i
		edges = append(edges, high, high|(1<<63), (2*high)-1, ((2*high)-1)|(1<<63))
	}
	check.CompareSamples("godiff-logb", edges, func(rng *rand.Rand, i int) uint64 {
		word := rng.Uint64()
		if i%2 == 0 {
			return word &^ (0x7ff << 52)
		}
		return word
	}, func(word uint64) string {
		return strconv.FormatUint(word, 10)
	}, func(word uint64) string {
		out := math.Logb(binary64.Store(word))
		raw := binary64.Load(out)
		if math.IsInf(out, 0) || math.IsNaN(out) {
			return special(raw)
		}
		if out < -1074 || out > 1023 {
			return "unexpected finite word " + strconv.FormatUint(raw, 10)
		}
		exponent := int(out)
		if binary64.Load(float64(exponent)) != raw {
			return "noncanonical finite word " + strconv.FormatUint(raw, 10)
		}
		return finite(exponent)
	}, expected)
}

func expected(word uint64) string {
	magnitude := word % (1 << 63)
	switch {
	case magnitude == 0:
		return special(0xfff0000000000000)
	case magnitude == 0x7ff0000000000000:
		return special(0x7ff0000000000000)
	case magnitude > 0x7ff0000000000000:
		return special(word)
	case magnitude < 1<<52:
		return finite(bits.Len64(magnitude) - 1075)
	default:
		return finite(int(magnitude/(1<<52)) - 1023)
	}
}

func special(word uint64) string { return "special " + strconv.FormatUint(word, 10) }
func finite(exponent int) string { return "finite " + strconv.Itoa(exponent) }
