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
	check.CompareSamples("godiff-float64frombits", binary64.Boundaries(),
		func(rng *rand.Rand, _ int) uint64 { return rng.Uint64() },
		decimal, func(word uint64) string {
			f := math.Float64frombits(word)
			if roundtrip := math.Float64bits(f); roundtrip != word {
				log.Fatalf("Float64bits(Float64frombits(word)): input=%#016x got=%#016x", word, roundtrip)
			}
			return decimal(binary64.Load(f))
		}, decimal)
}

func decimal(word uint64) string {
	return strconv.FormatUint(word, 10)
}
