package main

import (
	"log"
	"math"
	"math/rand"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary32"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	check.CompareSamples("godiff-float32bits", binary32.Boundaries(),
		func(rng *rand.Rand, _ int) uint32 { return rng.Uint32() },
		decimal, func(word uint32) string {
			f := binary32.Store(word)
			got := math.Float32bits(f)
			if roundtrip := binary32.Load(math.Float32frombits(got)); roundtrip != word {
				log.Fatalf("Float32frombits(Float32bits(f)): input=%#010x got=%#010x", word, roundtrip)
			}
			return decimal(got)
		}, decimal)
}

func decimal(word uint32) string {
	return strconv.FormatUint(uint64(word), 10)
}
