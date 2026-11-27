package main

import (
	"fmt"
	"math"
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

func main() {
	edges := binary64.Boundaries()
	for k := range 52 {
		power := uint64(1) << k
		for _, sign := range []uint64{0, signWeight} {
			edges = append(edges, sign+power-1, sign+power, sign+power+1)
		}
	}
	check.CompareSamples("godiff-frexp", edges, func(rng *rand.Rand, i int) uint64 {
		word := rng.Uint64()
		if i%2 == 0 {
			return word & (signWeight | (fractionWeight - 1))
		}
		return word
	}, func(word uint64) string {
		return strconv.FormatUint(word, 10)
	}, func(word uint64) string {
		frac, exp := math.Frexp(binary64.Store(word))
		return result(binary64.Load(frac), exp)
	}, expected)
}

func expected(word uint64) string {
	magnitude := word % signWeight
	sign := word - magnitude
	e, f := magnitude/fractionWeight, magnitude%fractionWeight
	switch {
	case magnitude == 0 || e == 2047:
		return result(word, 0)
	case e == 0:
		top := bits.Len64(f) - 1
		payload := (f - (uint64(1) << top)) << (52 - top)
		return result(sign+1022*fractionWeight+payload, top-1073)
	default:
		return result(sign+1022*fractionWeight+f, int(e)-1022)
	}
}

func result(fraction uint64, exponent int) string {
	return fmt.Sprintf("%d %d", fraction, exponent)
}
