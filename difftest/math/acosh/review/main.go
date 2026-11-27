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
	// Finite golden values were obtained independently with platform libm.
	tests := []fixture{
		{1, 0, true},
		{0x1.0000000000001p+0, 0x1.6a09e667f3bcdp-26, false},
		{0x1.0000000001000p+0, 0x1.6a09e667f39eap-20, false},
		{1.25, 0x1.62e42fefa39efp-1, false}, {1.5, 0x1.ecc2caec5160ap-1, false},
		{0x1.fffffffffffffp+0, 0x1.5124271980434p+0, false},
		{2, 0x1.5124271980434p+0, false},
		{0x1.0000000000001p+1, 0x1.5124271980436p+0, false},
		{10, 0x1.7f21ed1ce05d4p+1, false},
		{0x1.fffffffffffffp+27, 0x1.419ecb712c480p+4, false},
		{0x1p28, 0x1.419ecb712c480p+4, false},
		{0x1.0000000000001p+28, 0x1.419ecb712c480p+4, false},
		{0x1p512, 0x1.6395a2079b70cp+8, false},
		{math.MaxFloat64, 0x1.633ce8fb9f87dp+9, false},
		{math.Inf(1), math.Inf(1), true},
		{math.Nextafter(1, 0), math.NaN(), true},
		{0, math.NaN(), true}, {-1, math.NaN(), true},
		{math.Inf(-1), math.NaN(), true}, {math.NaN(), math.NaN(), true},
	}
	const profiles = 3
	var input strings.Builder
	var clippedDomain, overflowingDouble int
	for _, test := range tests {
		for profile := range profiles {
			fmt.Fprintf(&input, "%d %d\n", profile, math.Float64bits(test.input))
		}
		got := math.Acosh(test.input)
		if !agrees(got, test) {
			log.Fatalf("Go Acosh(%g) = %.17g, want %.17g", test.input, got, test.want)
		}
		if test.input < 1 && !agrees(0, test) {
			clippedDomain++
		}
		if test.input > math.MaxFloat64/2 && !agrees(math.Log(2*test.input), test) {
			overflowingDouble++
		}
	}
	rows := check.LeanRows("godiff-math-acosh", profiles*len(tests), strings.NewReader(input.String()))
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
			log.Fatalf("Lean profile %d Acosh(%g) = %.17g, want %.17g", i%profiles, test.input, got, test.want)
		}
	}
	if clippedDomain == 0 || overflowingDouble == 0 {
		log.Fatal("controls did not distinguish both faults")
	}
	fmt.Printf("Acosh passed %d native and %d Lean profile fixtures; clipping the domain failed %d, doubling large finite magnitudes failed %d\n",
		len(tests), len(rows), clippedDomain, overflowingDouble)
}
