package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	bits uint64
	sign int
}

func main() {
	signs := []int{math.MinInt, -2, -1, 0, 1, 2, math.MaxInt}
	var edges []input
	for _, sign := range signs {
		for _, signBit := range []uint64{0, 1 << 63} {
			for _, bits := range []uint64{
				0, 1, 0x7feffffffffffffe, 0x7fefffffffffffff,
				0x7ff0000000000000, 0x7ff0000000000001, 0x7ff8000000000000, 0x7fffffffffffffff,
			} {
				edges = append(edges, input{bits | signBit, sign})
			}
		}
	}
	check.CompareSamples("godiff-isinf", edges, func(rng *rand.Rand, i int) input {
		bits := rng.Uint64()
		if i%2 == 0 {
			bits = (bits & (1 << 63)) | 0x7ff0000000000000
		}
		return input{bits, signs[rng.Intn(len(signs))]}
	}, func(in input) string {
		return fmt.Sprintf("%d %d", in.bits, in.sign)
	}, func(in input) string {
		return strconv.FormatBool(math.IsInf(math.Float64frombits(in.bits), in.sign))
	}, expected)
}

func expected(in input) string {
	positive := in.bits == 0x7ff0000000000000
	negative := in.bits == 0xfff0000000000000
	return strconv.FormatBool((in.sign >= 0 && positive) || (in.sign <= 0 && negative))
}
