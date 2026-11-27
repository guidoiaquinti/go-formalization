package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/bits"
	"math/rand"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const (
	fractionWeight = 1 << 52
	signWeight     = 1 << 63
)

type input struct {
	word uint64
	exp  int
}

func main() {
	contract := flag.Bool("contract-counterexample", false, "check the documented scaling law at native int extrema")
	flag.Parse()
	if *contract {
		// The mathematical product overflows at MaxInt and underflows at MinInt.
		cases := []struct {
			in   input
			want uint64
		}{
			{input{0x4000000000000000, math.MaxInt}, 0x7ff0000000000000},
			{input{0x3fe0000000000000, math.MinInt}, 0},
		}
		failed := false
		for _, c := range cases {
			if got := actual(c.in); got != c.want {
				failed = true
				log.Printf("word=%#016x exp=%d: Go=%#016x scaling law=%#016x", c.in.word, c.in.exp, got, c.want)
			}
		}
		if failed {
			log.Fatal("documented Ldexp scaling law disagrees at native integer extrema")
		}
		return
	}

	words := append(binary64.Boundaries(),
		0x3fe0000000000000, 0xbfe0000000000000,
		0x4000000000000000, 0xc000000000000000,
		0x3ff0000000000001, 0x3ff0000000000002,
		0x3ff0000000000003, 0x3fffffffffffffff)
	var boundaries []input
	for _, word := range words {
		_, _, highest := fields(word)
		for _, exp := range []int{math.MinInt, math.MinInt + 1, math.MinInt + 1074,
			-1076 - highest, -1075 - highest, -1074 - highest,
			-1024 - highest, -1023 - highest, -1022 - highest, -1021 - highest,
			1022 - highest, 1023 - highest, 1024 - highest,
			-1, 0, 1, math.MaxInt - 1023, math.MaxInt - 1, math.MaxInt} {
			boundaries = append(boundaries, input{word, exp})
		}
	}
	check.CompareSamples("godiff-ldexp", boundaries, sample, func(c input) string {
		return fmt.Sprintf("%d %d %d", strconv.IntSize, c.word, c.exp)
	}, func(c input) string { return decimal(actual(c)) },
		func(c input) string { return decimal(expected(c)) })
}

func sample(rng *rand.Rand, index int) input {
	word := rng.Uint64()
	if index%2 == 0 {
		word &= signWeight | (fractionWeight - 1)
	}
	_, _, highest := fields(word)
	var exp int
	switch index % 5 {
	case 0:
		exp = int(rng.Uint64())
	case 1:
		exp = math.MinInt + rng.Intn(1100)
	case 2:
		exp = math.MaxInt - rng.Intn(1100)
	case 3:
		exp = rng.Intn(65) - 1080 - highest
	default:
		exp = rng.Intn(2200) - 1100 - highest
	}
	return input{word, exp}
}

func actual(c input) uint64 {
	return binary64.Load(math.Ldexp(binary64.Store(c.word), c.exp))
}

func decimal(word uint64) string {
	return strconv.FormatUint(word, 10)
}

// fields also gives the unbiased exponent of the highest set value bit.
func fields(word uint64) (exponent, payload uint64, highest int) {
	exponent, payload = (word>>52)&2047, word%fractionWeight
	if exponent == 0 {
		return exponent, payload, bits.Len64(payload) - 1 - 1074
	}
	return exponent, payload, int(exponent) - 1023
}

// Compute the exact dyadic product at the wrapped exponent, then rely on
// big.Float's independent rounding implementation to convert to binary64.
func expected(c input) uint64 {
	e, payload, highest := fields(c.word)
	if e == 0 && payload == 0 || e == 2047 {
		return c.word
	}
	sign := c.word & signWeight
	wrapped := c.exp + highest
	if wrapped < -1075 {
		return sign
	}
	if wrapped > 1023 {
		return sign | 0x7ff0000000000000
	}
	coefficient := fractionWeight + payload
	if e == 0 {
		coefficient = payload << (53 - bits.Len64(payload))
	}
	f := new(big.Float).SetUint64(coefficient)
	f.SetMantExp(f, wrapped-52)
	rounded, _ := f.Float64()
	return binary64.Load(rounded) | sign
}
