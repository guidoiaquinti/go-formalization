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
		{0x1p-28, 0x1p-28, false},
		{0.125, 0x1.fead0b6996972p-4, false},
		{-0.5, -0x1.ecc2caec5160ap-2, false},
		{0x1.fffffffffffffp+0, 0x1.719218313d086p+0, false},
		{2, 0x1.719218313d087p+0, false},
		{0x1.0000000000001p+1, 0x1.719218313d088p+0, false},
		{0x1p28, 0x1.419ecb712c480p+4, false},
		{0x1.0000000000001p+28, 0x1.419ecb712c480p+4, false},
		{-0x1p1023, -0x1.62e42fefa39efp+9, false},
		{math.MaxFloat64, 0x1.633ce8fb9f87dp+9, false},
		{math.Inf(1), math.Inf(1), true}, {math.Inf(-1), math.Inf(-1), true},
		{math.NaN(), math.NaN(), true},
	}
	const profiles = 3
	var input strings.Builder
	var lostSign, overflowingDouble, linear int
	for _, test := range tests {
		for profile := range profiles {
			fmt.Fprintf(&input, "%d %d\n", profile, math.Float64bits(test.input))
		}
		got := math.Asinh(test.input)
		if !agrees(got, test) {
			log.Fatalf("Go Asinh(%g) = %.17g, want %.17g", test.input, got, test.want)
		}
		if !agrees(math.Abs(got), test) {
			lostSign++
		}
		if math.Abs(test.input) > math.MaxFloat64/2 && !agrees(math.Copysign(math.Log(2*math.Abs(test.input)), test.input), test) {
			overflowingDouble++
		}
		if math.Abs(test.input) <= 1 && !agrees(test.input, test) {
			linear++
		}
	}
	rows := check.LeanRows("godiff-math-asinh", profiles*len(tests), strings.NewReader(input.String()))
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
			log.Fatalf("Lean profile %d Asinh(%g) = %.17g, want %.17g", i%profiles, test.input, got, test.want)
		}
	}
	if lostSign == 0 || overflowingDouble == 0 || linear == 0 {
		log.Fatal("controls did not distinguish all three faults")
	}
	fmt.Printf("Asinh passed %d native and %d Lean profile fixtures; losing signs failed %d, doubling large finite magnitudes failed %d, substituting the identity failed %d\n",
		len(tests), len(rows), lostSign, overflowingDouble, linear)
}
