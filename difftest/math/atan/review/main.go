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
		{0.125, 0x1.fd5ba9aac2f6ep-4, false}, {8, 0x1.7249faa996a21p+0, false},
		{0.25, 0x1.f5b75f92c80ddp-3, false}, {4, 0x1.5368c951e9cfdp+0, false},
		{-0.5, -0x1.dac670561bb50p-2, false}, {-2, -0x1.1b6e192ebbe45p+0, false},
		{0x1.51eb851eb851ep-1, 0x1.2aafdde4d0c9ep-1, false},
		{0x1.51eb851eb851fp-1, 0x1.2aafdde4d0c9fp-1, false},
		{0x1.51eb851eb8520p-1, 0x1.2aafdde4d0ca0p-1, false},
		{1, math.Pi / 4, false},
		{0x1.3504f333f9de5p+1, 0x1.2d97c7f3321d2p+0, false},
		{0x1.3504f333f9de6p+1, 0x1.2d97c7f3321d2p+0, false},
		{0x1.3504f333f9de7p+1, 0x1.2d97c7f3321d2p+0, false},
		{16, 0x1.82250768ac529p+0, false},
		{1e25, math.Pi / 2, true},
		{0, 0, true}, {negativeZero, negativeZero, true},
		{math.Inf(1), math.Pi / 2, true}, {math.Inf(-1), -math.Pi / 2, true},
		{math.NaN(), math.NaN(), true},
	}
	const profiles = 2
	var input strings.Builder
	var lostSign, lostShift int
	for _, test := range tests {
		for profile := range profiles {
			fmt.Fprintf(&input, "%d %d\n", profile, math.Float64bits(test.input))
		}
		got := math.Atan(test.input)
		if !agrees(got, test) {
			log.Fatalf("Go Atan(%g) = %.17g, want %.17g", test.input, got, test.want)
		}
		if !agrees(math.Abs(got), test) {
			lostSign++
		}
		x := math.Abs(test.input)
		if x > 2.414213562373095 && !agrees(-math.Copysign(math.Atan(1/x), test.input), test) {
			lostShift++
		}
	}
	rows := check.LeanRows("godiff-math-atan", profiles*len(tests), strings.NewReader(input.String()))
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
			log.Fatalf("Lean profile %d Atan(%g) = %.17g, want %.17g", i%profiles, test.input, got, test.want)
		}
	}
	for i := 0; i < 6; i += 2 {
		for profile := range profiles {
			sum := values[profiles*i+profile] + values[profiles*(i+1)+profile]
			if math.Abs(sum-math.Copysign(math.Pi/2, tests[i].input)) > 0x1p-51 {
				log.Fatalf("Lean profile %d failed reciprocal identity at %g", profile, tests[i].input)
			}
		}
	}
	if lostSign == 0 || lostShift == 0 {
		log.Fatal("controls did not distinguish both faults")
	}
	fmt.Printf("Atan passed %d native and %d Lean profile fixtures; losing signs failed %d, omitting the large-argument offset failed %d\n",
		len(tests), len(rows), lostSign, lostShift)
}
