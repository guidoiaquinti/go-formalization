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
	y, x, want float64
	exact      bool
}

func agrees(got float64, test fixture) bool {
	if math.IsNaN(test.want) {
		return math.IsNaN(got)
	}
	if test.exact {
		return math.Float64bits(got) == math.Float64bits(test.want)
	}
	return math.Abs(got-test.want) <= 0x1p-50
}

func main() {
	negZero := math.Copysign(0, -1)
	inf := math.Inf(1)
	tiny := math.SmallestNonzeroFloat64
	tests := []fixture{
		{0, 0, 0, true}, {negZero, 0, negZero, true},
		{0, negZero, math.Pi, true}, {negZero, negZero, -math.Pi, true},
		{1, negZero, math.Pi / 2, true}, {-1, negZero, -math.Pi / 2, true},
		{inf, inf, math.Pi / 4, true}, {-inf, -inf, -3 * math.Pi / 4, true},
		{-1, inf, negZero, true}, {-1, -inf, -math.Pi, true},
		{math.NaN(), 1, math.NaN(), true}, {1, math.NaN(), math.NaN(), true},
		{1, 2, 0x1.dac670561bb50p-2, false},
		{1, -1, 3 * math.Pi / 4, false}, {-1, -1, -3 * math.Pi / 4, false},
		{tiny, -2, math.Pi, true},
		{-tiny, -math.Nextafter(2, 0), -math.Pi, true},
		{-tiny, -2, math.Pi, true},
		{-tiny, -math.Nextafter(2, inf), math.Pi, true},
		{-tiny, -inf, -math.Pi, true},
	}
	const profiles = 2
	var inputs strings.Builder
	var lostSign, mathematicalQuadrant int
	for _, test := range tests {
		for profile := range profiles {
			fmt.Fprintf(&inputs, "%d %d %d\n", profile, math.Float64bits(test.y), math.Float64bits(test.x))
		}
		got := math.Atan2(test.y, test.x)
		if !agrees(got, test) {
			log.Fatalf("Go Atan2(%g, %g) = %.17g, want %.17g", test.y, test.x, got, test.want)
		}
		if !agrees(math.Abs(got), test) {
			lostSign++
		}
		if test.y < 0 && test.x < 0 && test.y/test.x == 0 && !agrees(-math.Pi, test) {
			mathematicalQuadrant++
		}
	}
	rows := check.LeanRows("godiff-math-atan2", profiles*len(tests), strings.NewReader(inputs.String()))
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
			log.Fatalf("Lean profile %d Atan2(%g, %g) = %.17g, want %.17g", i%profiles, test.y, test.x, got, test.want)
		}
	}
	if lostSign == 0 || mathematicalQuadrant != 2 {
		log.Fatalf("expected sign and underflow differences, got %d and %d", lostSign, mathematicalQuadrant)
	}
	fmt.Printf("Atan2 passed %d native and %d Lean fixtures; lost sign failed %d, substituting the mathematical quadrant failed %d\n",
		len(tests), len(rows), lostSign, mathematicalQuadrant)
}
