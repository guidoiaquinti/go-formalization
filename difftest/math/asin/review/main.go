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
	negativeZero := math.Copysign(0, -1)
	// Finite golden values were obtained independently with platform libm.
	tests := []fixture{
		{0, 0, true}, {negativeZero, negativeZero, true},
		{math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, true},
		{-math.SmallestNonzeroFloat64, -math.SmallestNonzeroFloat64, true},
		{0x1p-27, 0x1p-27, true},
		{0.125, 0x1.00abe0c129e1ep-3, false}, {-0.3, -0x1.380159e14f6fep-2, false},
		{0.5, 0x1.0c152382d7365p-1, false},
		{0x1.6666666666665p-1, 0x1.8d00e692afd94p-1, false},
		{0.7, 0x1.8d00e692afd95p-1, false},
		{0x1.6666666666667p-1, 0x1.8d00e692afd97p-1, false},
		{0.75, 0x1.b235315c680dcp-1, false}, {-0.875, -0x1.10c066d3e6932p+0, false},
		{0x1.fffffffffffffp-1, 0x1.921fb50442d18p+0, false},
		{-0x1.fffffffffffffp-1, -0x1.921fb50442d18p+0, false},
		{1, math.Pi / 2, true}, {-1, -math.Pi / 2, true},
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
		got := math.Asin(test.input)
		if !agrees(got, test) {
			log.Fatalf("Go Asin(%g) = %.17g, want %.17g", test.input, got, test.want)
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
	rows := check.LeanRows("godiff-math-asin", profiles*len(tests), strings.NewReader(input.String()))
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
			log.Fatalf("Lean profile %d Asin(%g) = %.17g, want %.17g", i%profiles, test.input, got, test.want)
		}
	}
	if lostSign == 0 || clampedDomain == 0 || linear == 0 {
		log.Fatal("controls did not distinguish all three faults")
	}
	fmt.Printf("Asin passed %d native and %d Lean profile fixtures; losing signs failed %d, clamping out-of-domain inputs failed %d, substituting the identity failed %d\n",
		len(tests), len(rows), lostSign, clampedDomain, linear)
}
