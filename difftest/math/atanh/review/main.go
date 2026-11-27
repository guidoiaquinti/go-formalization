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
	return math.Abs(got-test.want) <= 0x1p-50*max(1, math.Abs(test.want))
}

func main() {
	negativeZero := math.Copysign(0, -1)
	// Finite golden values were obtained independently with platform libm.
	tests := []fixture{
		{0, 0, true}, {negativeZero, negativeZero, true},
		{math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, true},
		{-math.SmallestNonzeroFloat64, -math.SmallestNonzeroFloat64, true},
		{0x1.fffffffffffffp-29, 0x1.fffffffffffffp-29, true},
		{0x1p-28, 0x1p-28, false}, {0x1.0000000000001p-28, 0x1.0000000000001p-28, false},
		{0.125, 0x1.015891c9eaef7p-3, false}, {-0.25, -0x1.058aefa811451p-2, false},
		{0x1.fffffffffffffp-2, 0x1.193ea7aad030ap-1, false},
		{0.5, 0x1.193ea7aad030bp-1, false},
		{0x1.0000000000001p-1, 0x1.193ea7aad030cp-1, false},
		{-0.75, -0x1.f2272ae325a57p-1, false}, {0.99, 0x1.52c581997cd85p+1, false},
		{0x1.fffffffffffffp-1, 0x1.2b708872320e2p+4, false},
		{-0x1.fffffffffffffp-1, -0x1.2b708872320e2p+4, false},
		{1, math.Inf(1), true}, {-1, math.Inf(-1), true},
		{math.Nextafter(1, math.Inf(1)), math.NaN(), true},
		{-math.Nextafter(1, math.Inf(1)), math.NaN(), true},
		{math.Inf(1), math.NaN(), true}, {math.Inf(-1), math.NaN(), true},
		{math.NaN(), math.NaN(), true},
	}
	const profiles = 3
	var input strings.Builder
	var lostSign, clampedDomain, linear int
	for _, test := range tests {
		for profile := range profiles {
			fmt.Fprintf(&input, "%d %d\n", profile, math.Float64bits(test.input))
		}
		got := math.Atanh(test.input)
		if !agrees(got, test) {
			log.Fatalf("Go Atanh(%g) = %.17g, want %.17g", test.input, got, test.want)
		}
		if !agrees(math.Abs(got), test) {
			lostSign++
		}
		if math.Abs(test.input) > 1 && !agrees(math.Copysign(math.Pi/2, test.input), test) {
			clampedDomain++
		}
		if math.Abs(test.input) <= 1 && !agrees(test.input, test) {
			linear++
		}
	}
	rows := check.LeanRows("godiff-math-atanh", profiles*len(tests), strings.NewReader(input.String()))
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
			log.Fatalf("Lean profile %d Atanh(%g) = %.17g, want %.17g", i%profiles, test.input, got, test.want)
		}
	}
	if lostSign == 0 || clampedDomain == 0 || linear == 0 {
		log.Fatal("controls did not distinguish all three faults")
	}
	fmt.Printf("Atanh passed %d native and %d Lean profile fixtures; losing signs failed %d, clamping out-of-domain inputs failed %d, substituting the identity failed %d\n",
		len(tests), len(rows), lostSign, clampedDomain, linear)
}
