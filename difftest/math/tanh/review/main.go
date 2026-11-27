package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	input, want float64
	exact       bool
}

func agrees(got float64, test fixture) bool {
	if math.IsNaN(test.want) {
		return math.IsNaN(got)
	}
	if test.exact {
		return math.Float64bits(got) == math.Float64bits(test.want)
	}
	return math.Abs(got-test.want) <= 0x1p-51
}

func main() {
	// The nontrivial finite results were generated with the platform libm.
	tests := []fixture{
		{0, 0, true},
		{math.Copysign(0, -1), math.Copysign(0, -1), true},
		{math.NaN(), math.NaN(), true},
		{math.Inf(1), 1, true},
		{math.Inf(-1), -1, true},
		{math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, true},
		{-math.SmallestNonzeroFloat64, -math.SmallestNonzeroFloat64, true},
		{0x1.3ffffffffffffp-1, 0x1.1bf47eabb8f95p-1, false},
		{0x1.4p-1, 0x1.1bf47eabb8f95p-1, false},
		{0x1.4000000000001p-1, 0x1.1bf47eabb8f96p-1, false},
		{-0x1.4000000000001p-1, -0x1.1bf47eabb8f96p-1, false},
		{1, 0x1.85efab514f394p-1, false},
		{2, 0x1.ed9505e1bc3d4p-1, false},
		{18, 0x1.ffffffffffffcp-1, true},
		{19, 0x1.fffffffffffffp-1, true},
		{20, 1, true},
		{0x1.601e678fc457ap+5, 1, true},
		{0x1.601e678fc457bp+5, 1, true},
		{0x1.601e678fc457cp+5, 1, true},
		{-0x1.601e678fc457cp+5, -1, true},
		{math.MaxFloat64, 1, true},
	}
	const profiles = 5
	var input strings.Builder
	var lostSign, linearSmall, earlySaturation int
	for _, test := range tests {
		for profile := range profiles {
			fmt.Fprintf(&input, "%d %d\n", profile, math.Float64bits(test.input))
		}
		got := math.Tanh(test.input)
		if !agrees(got, test) {
			log.Fatalf("Go Tanh(%g) = %.17g, want %.17g", test.input, got, test.want)
		}
		if !agrees(math.Abs(got), test) {
			lostSign++
		}
		if math.Abs(test.input) < 0.625 && !agrees(test.input, test) {
			linearSmall++
		}
		if math.Abs(test.input) >= 18 && !agrees(math.Copysign(1, test.input), test) {
			earlySaturation++
		}
	}
	rows := check.LeanRows("godiff-math-tanh", profiles*len(tests), strings.NewReader(input.String()))
	for i, row := range rows {
		got := math.NaN()
		if row != "NaN" {
			word, err := strconv.ParseUint(row, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			got = math.Float64frombits(word)
		}
		test := tests[i/profiles]
		if !agrees(got, test) {
			log.Fatalf("Lean profile %d Tanh(%g) = %.17g, want %.17g", i%profiles, test.input, got, test.want)
		}
	}
	if lostSign == 0 || linearSmall == 0 || earlySaturation == 0 {
		log.Fatalf("control misses a mutant, lost sign=%d linear=%d early saturation=%d", lostSign, linearSmall, earlySaturation)
	}
	fmt.Printf("Tanh passed %d native and %d Lean profile fixtures; lost sign failed %d, forcing linear approximation failed %d, early saturation failed %d\n",
		len(tests), len(rows), lostSign, linearSmall, earlySaturation)
}
