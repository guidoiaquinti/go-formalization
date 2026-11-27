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
	return math.Abs(got-test.want) <= 0x1p-49
}

func main() {
	negativeZero := math.Copysign(0, -1)
	// Finite golden values were obtained independently with platform libm.
	tests := []fixture{
		{0.125, 0x1.720a392c1d954p+0, false}, {-0.125, 0x1.b235315c680dcp+0, false},
		{0.5, 0x1.0c152382d7365p+0, false}, {-0.5, 0x1.0c152382d7366p+1, false},
		{0.99, 0x1.21df72882bfd8p-3, false}, {-0.99, 0x1.8001be1bc011bp+1, false},
		{0x1.fffffffffffffp-1, 0x1p-26, false}, {-0x1.fffffffffffffp-1, 0x1.921fb52442d18p+1, false},
		{0, math.Pi / 2, true}, {negativeZero, math.Pi / 2, true},
		{0x1.6666666666665p-1, 0x1.973e83f5d5c9dp-1, false},
		{0.7, 0x1.973e83f5d5c9bp-1, false},
		{0x1.6666666666667p-1, 0x1.973e83f5d5c9ap-1, false},
		{1, 0, true}, {-1, math.Pi, true},
		{math.Nextafter(1, math.Inf(1)), math.NaN(), true},
		{-math.Nextafter(1, math.Inf(1)), math.NaN(), true},
		{math.Inf(1), math.NaN(), true}, {math.Inf(-1), math.NaN(), true},
		{math.NaN(), math.NaN(), true},
	}
	const profiles = 3
	var input strings.Builder
	var reversedSign, clampedDomain int
	for _, test := range tests {
		for profile := range profiles {
			fmt.Fprintf(&input, "%d %d\n", profile, math.Float64bits(test.input))
		}
		got := math.Acos(test.input)
		if !agrees(got, test) {
			log.Fatalf("Go Acos(%g) = %.17g, want %.17g", test.input, got, test.want)
		}
		if !agrees(math.Pi-got, test) {
			reversedSign++
		}
		if math.Abs(test.input) > 1 && !agrees(math.Pi, test) {
			clampedDomain++
		}
	}
	rows := check.LeanRows("godiff-math-acos", profiles*len(tests), strings.NewReader(input.String()))
	values := make([]float64, len(rows))
	for i, row := range rows {
		got := math.NaN()
		if row != "NaN" {
			word, err := strconv.ParseUint(row, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			got = math.Float64frombits(word)
		}
		values[i] = got
		test := tests[i/profiles]
		if !agrees(got, test) {
			log.Fatalf("Lean profile %d Acos(%g) = %.17g, want %.17g", i%profiles, test.input, got, test.want)
		}
	}
	for i := 0; i < 8; i += 2 {
		for profile := range profiles {
			sum := values[i*profiles+profile] + values[(i+1)*profiles+profile]
			if math.Abs(sum-math.Pi) > 0x1p-49 {
				log.Fatalf("Lean profile %d failed complement identity at %g", profile, tests[i].input)
			}
		}
	}
	if reversedSign == 0 || clampedDomain == 0 {
		log.Fatal("controls did not distinguish both faults")
	}
	fmt.Printf("Acos passed %d native and %d Lean profile fixtures; reversing the Asin term failed %d, clamping out-of-domain inputs failed %d\n",
		len(tests), len(rows), reversedSign, clampedDomain)
}
