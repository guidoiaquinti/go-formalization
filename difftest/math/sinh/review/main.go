package main

import (
	"fmt"
	"log"
	"math"
	"runtime"
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
		{0.125, 0x1.00aaccd00d2f1p-3, false},
		{0x1.fffffffffffffp-2, 0x1.0acd00fe63b96p-1, false},
		{0.5, 0x1.0acd00fe63b97p-1, false},
		{0x1.0000000000001p-1, 0x1.0acd00fe63b98p-1, false},
		{-1, -0x1.2cd9fc44eb982p+0, false},
		{5, 0x1.28d0166f07374p+6, false},
		{0x1.4ffffffffffffp+4, 0x1.3a6e1fd9eeceap+29, false},
		{21, 0x1.3a6e1fd9eecfdp+29, false},
		{0x1.5000000000001p+4, 0x1.3a6e1fd9eed11p+29, false},
		{100, 0x1.3494a9b171bf5p+143, false},
		{709, 0x1.d422d2be5dc9bp+1021, false},
		{709.75, 0x1.ef85a11e73f2dp+1022, false},
		{709.8, math.Inf(1), true}, {710, math.Inf(1), true},
		{710.475, math.Inf(1), true}, {-710, math.Inf(-1), true},
		{710.5, math.Inf(1), true},
		{math.Inf(1), math.Inf(1), true}, {math.Inf(-1), math.Inf(-1), true},
		{math.NaN(), math.NaN(), true},
	}
	const profiles = 5
	var input strings.Builder
	var lostSign, lostHalf, correctedOverflow int
	// These exact rounded finite values come from 110-digit decimal Exp.
	reference := map[float64]float64{
		709.8:   0x1.0476c6ea64149p+1023,
		710:     0x1.3e21a464507f9p+1023,
		710.475: 0x1.ff8f511d45587p+1023,
		-710:    -0x1.3e21a464507f9p+1023,
	}
	for _, test := range tests {
		for profile := range profiles {
			fmt.Fprintf(&input, "%d %d\n", profile, math.Float64bits(test.input))
		}
		got := math.Sinh(test.input)
		if test.input == 709.75 && runtime.GOARCH == "amd64" {
			test.want, test.exact = math.Inf(1), true
		}
		if !agrees(got, test) {
			log.Fatalf("Go Sinh(%g) = %.17g, want %.17g", test.input, got, test.want)
		}
		if !agrees(math.Abs(got), test) {
			lostSign++
		}
		if test.input > 21 && !agrees(2*got, test) {
			lostHalf++
		}

		if finite, found := reference[test.input]; found && !agrees(finite, test) {
			correctedOverflow++
		}
	}
	rows := check.LeanRows("godiff-math-sinh", profiles*len(tests), strings.NewReader(input.String()))
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
		// On amd64, the rounded reduction exponent is already 1024.
		if test.input == 709.75 && i%profiles < 4 {
			test.want, test.exact = math.Inf(1), true
		}
		if !agrees(got, test) {
			log.Fatalf("Lean profile %d Sinh(%g) = %.17g, want %.17g", i%profiles, test.input, got, test.want)
		}
	}
	if lostSign == 0 || lostHalf == 0 || correctedOverflow != 4 {
		log.Fatal("controls did not distinguish sign, half factor and pinned overflow")
	}
	fmt.Printf("Sinh passed %d native and %d Lean profile fixtures; losing signs failed %d, omitting the half factor failed %d, finite reference differed on %d overflow fixtures\n",
		len(tests), len(rows), lostSign, lostHalf, correctedOverflow)
}
