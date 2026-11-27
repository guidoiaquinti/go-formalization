package main

import (
	"log"
	"math"
	"math/rand"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	check.CompareSamples("godiff-float64bits", binary64.Boundaries(),
		func(rng *rand.Rand, _ int) uint64 { return rng.Uint64() },
		decimal, func(word uint64) string {
			f := binary64.Store(word)
			got := math.Float64bits(f)
			if roundtrip := binary64.Load(math.Float64frombits(got)); roundtrip != word {
				log.Fatalf("Float64frombits(Float64bits(f)): input=%#016x got=%#016x", word, roundtrip)
			}
			return decimal(got)
		}, decimal)
}

func decimal(word uint64) string {
	return strconv.FormatUint(word, 10)
}
