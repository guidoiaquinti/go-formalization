package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
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
		{0.125, 0x1.015891c9eaef7p-3, false}, {0.25, 0x1.058aefa811451p-2, false},
		{4.0 / 11, 0x1.8637388d06bcep-2, false},
		{-0.75, -0x1.f2272ae325a57p-1, false}, {-20.0 / 29, -0x1.b1d10670aae9ap-1, false},
		{0.6, 0x1.62e42fefa39efp-1, false}, {-0.875, -0x1.5aa16394d481fp+0, false},
		{0.95, 0x1.d4ef968880dd4p+0, false}, {0x1.fffffffffffffp-1, 0x1.2b708872320e2p+4, false},
		{0, 0, true}, {negativeZero, negativeZero, true},
		{math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, true},
		{-math.SmallestNonzeroFloat64, -math.SmallestNonzeroFloat64, true},
		{1, math.Inf(1), true}, {-1, math.Inf(-1), true},
		{2, math.NaN(), true}, {-2, math.NaN(), true},
		{math.Inf(1), math.NaN(), true}, {math.NaN(), math.NaN(), true},
	}
	const profiles = 3
	var input strings.Builder
	var cases bytes.Buffer
	var lostSign, clampedDomain, linear int
	for _, test := range tests {
		for profile := range profiles {
			fmt.Fprintf(&input, "%d %d\n", profile, math.Float64bits(test.input))
		}
		got := math.Atanh(test.input)
		if !agrees(got, test) {
			log.Fatalf("Go Atanh(%g) = %.17g, want %.17g", test.input, got, test.want)
		}
		observation := "NaN"
		if !math.IsNaN(got) {
			observation = strconv.FormatUint(math.Float64bits(got), 10)
		}
		kind := "numeric"
		if test.exact {
			kind = "exact"
		}
		reference := "NaN"
		if !math.IsNaN(test.want) {
			reference = strconv.FormatUint(math.Float64bits(test.want), 10)
		}
		fmt.Fprintln(&cases, math.Float64bits(test.input), observation, reference, kind)
		if !agrees(math.Abs(got), test) {
			lostSign++
		}
		if math.Abs(test.input) > 1 && !agrees(math.Copysign(math.Inf(1), test.input), test) {
			clampedDomain++
		}
		if math.Abs(test.input) <= 1 && !agrees(test.input, test) {
			linear++
		}
	}
	rows := check.LeanRows("godiff-math-atanh", profiles*len(tests), strings.NewReader(input.String()))
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
			log.Fatalf("Lean profile %d Atanh(%g) = %.17g, want %.17g", i%profiles, test.input, got, test.want)
		}
	}
	for _, triple := range [][3]int{{0, 1, 2}, {0, 3, 4}} {
		for profile := range profiles {
			x := values[triple[0]*profiles+profile]
			y := values[triple[1]*profiles+profile]
			sum := values[triple[2]*profiles+profile]
			if math.Abs(x+y-sum) > 0x1p-50 {
				log.Fatalf("Lean profile %d failed the inverse hyperbolic addition identity", profile)
			}
		}
	}
	dir, err := os.MkdirTemp("", "golean-private-atanh-review-")
	must(err)
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	must(os.WriteFile(path, cases.Bytes(), 0600))
	must(os.Setenv("GOLEAN_MATH_ATANH_CASES", path))
	const overlay = "difftest/math/atanh/testdata/"
	must(check.TestInPackage("math", "TestGoLeanAtanh", overlay+"bridge_test.go", overlay+"atanh_test.go"))
	if lostSign == 0 || clampedDomain == 0 || linear == 0 {
		log.Fatal("controls did not distinguish all three faults")
	}
	fmt.Printf("Private atanh passed %d native and %d Lean profile fixtures; losing signs failed %d, clamping out-of-domain inputs failed %d, substituting the identity failed %d\n",
		len(tests), len(rows), lostSign, clampedDomain, linear)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
