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
	// Finite golden values were obtained independently with platform libm.
	tests := []fixture{
		{1.25, 0x1.62e42fefa39efp-1, false},
		{2.125, 0x1.62e42fefa39efp+0, false},
		{8.03125, 0x1.62e42fefa39efp+1, false},
		{0x1.0000000000010p+0, 0x1.6a09e667f3bcap-24, false},
		{1.125, 0x1.facfb2399e637p-2, false},
		{16, 0x1.bb7d2fe4dd78ep+1, false},
		{0x1p52, 0x1.25e4f7b2737fap+5, false},
		{0x1p1023, 0x1.62e42fefa39efp+9, false},
		{math.MaxFloat64, 0x1.633ce8fb9f87dp+9, false},
		{1, 0, true}, {math.Inf(1), math.Inf(1), true},
		{math.Nextafter(1, 0), math.NaN(), true},
		{0, math.NaN(), true}, {-1, math.NaN(), true},
		{math.Inf(-1), math.NaN(), true}, {math.NaN(), math.NaN(), true},
	}
	const profiles = 3
	var input strings.Builder
	var cases bytes.Buffer
	var clippedDomain, overflowingDouble int
	for _, test := range tests {
		for profile := range profiles {
			fmt.Fprintf(&input, "%d %d\n", profile, math.Float64bits(test.input))
		}
		got := math.Acosh(test.input)
		if !agrees(got, test) {
			log.Fatalf("Go Acosh(%g) = %.17g, want %.17g", test.input, got, test.want)
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
		if test.input < 1 && !agrees(0, test) {
			clippedDomain++
		}
		if test.input > math.MaxFloat64/2 && !agrees(math.Log(2*test.input), test) {
			overflowingDouble++
		}
	}
	rows := check.LeanRows("godiff-math-acosh", profiles*len(tests), strings.NewReader(input.String()))
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
			log.Fatalf("Lean profile %d Acosh(%g) = %.17g, want %.17g", i%profiles, test.input, got, test.want)
		}
	}
	for i := 0; i < 2; i++ {
		for profile := range profiles {
			x := values[i*profiles+profile]
			doubled := values[(i+1)*profiles+profile]
			if math.Abs(2*x-doubled) > 0x1p-50 {
				log.Fatalf("Lean profile %d failed the inverse hyperbolic double-angle identity", profile)
			}
		}
	}
	dir, err := os.MkdirTemp("", "golean-private-acosh-review-")
	must(err)
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	must(os.WriteFile(path, cases.Bytes(), 0600))
	must(os.Setenv("GOLEAN_MATH_ACOSH_CASES", path))
	const overlay = "difftest/math/acosh/testdata/"
	must(check.TestInPackage("math", "TestGoLeanAcosh", overlay+"bridge_test.go", overlay+"acosh_test.go"))
	if clippedDomain == 0 || overflowingDouble == 0 {
		log.Fatal("controls did not distinguish both faults")
	}
	fmt.Printf("Private acosh passed %d native and %d Lean profile fixtures; clipping the domain failed %d, doubling large finite magnitudes failed %d\n",
		len(tests), len(rows), clippedDomain, overflowingDouble)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
