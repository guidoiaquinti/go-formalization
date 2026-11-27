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
	check.CompareSamples("godiff-float32frombits", binary32.Boundaries(),
		func(rng *rand.Rand, _ int) uint32 { return rng.Uint32() },
		decimal, func(word uint32) string {
			f := math.Float32frombits(word)
			if roundtrip := math.Float32bits(f); roundtrip != word {
				log.Fatalf("Float32bits(Float32frombits(word)): input=%#010x got=%#010x", word, roundtrip)
			}
			return decimal(binary32.Load(f))
		}, decimal)
}

func decimal(word uint32) string {
	return strconv.FormatUint(uint64(word), 10)
}
