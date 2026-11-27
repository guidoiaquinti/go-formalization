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
		{0.75, 0x1.62e42fefa39efp-1, false},
		{1.875, 0x1.62e42fefa39efp+0, false},
		{7.96875, 0x1.62e42fefa39efp+1, false},
		{0.375, 0x1.7786a7973a592p-2, false},
		{-0.9, -0x1.9e23ce96e38d4p-1, false},
		{0.99, 0x1.bfa2424777c63p-1, false},
		{1.99, 0x1.706c6bca62894p+0, false},
		{0x1.fffffffffffffp+27, 0x1.419ecb712c480p+4, false},
		{0x1.0000000000001p+28, 0x1.419ecb712c480p+4, false},
		{0x1p512, 0x1.6395a2079b70cp+8, false},
		{-math.MaxFloat64, -0x1.633ce8fb9f87dp+9, false},
		{0, 0, true}, {negativeZero, negativeZero, true},
		{math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, true},
		{-math.SmallestNonzeroFloat64, -math.SmallestNonzeroFloat64, true},
		{math.Inf(1), math.Inf(1), true}, {math.Inf(-1), math.Inf(-1), true},
		{math.NaN(), math.NaN(), true},
	}
	const profiles = 3
	var input strings.Builder
	var cases bytes.Buffer
	var lostSign, overflowingDouble, linear int
	for _, test := range tests {
		for profile := range profiles {
			fmt.Fprintf(&input, "%d %d\n", profile, math.Float64bits(test.input))
		}
		got := math.Asinh(test.input)
		if !agrees(got, test) {
			log.Fatalf("Go Asinh(%g) = %.17g, want %.17g", test.input, got, test.want)
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
		if math.Abs(test.input) > math.MaxFloat64/2 && !agrees(math.Copysign(math.Log(2*math.Abs(test.input)), test.input), test) {
			overflowingDouble++
		}
		if math.Abs(test.input) <= 1 && !agrees(test.input, test) {
			linear++
		}
	}
	rows := check.LeanRows("godiff-math-asinh", profiles*len(tests), strings.NewReader(input.String()))
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
			log.Fatalf("Lean profile %d Asinh(%g) = %.17g, want %.17g", i%profiles, test.input, got, test.want)
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
	dir, err := os.MkdirTemp("", "golean-private-asinh-review-")
	must(err)
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	must(os.WriteFile(path, cases.Bytes(), 0600))
	must(os.Setenv("GOLEAN_MATH_ASINH_CASES", path))
	const overlay = "difftest/math/asinh/testdata/"
	must(check.TestInPackage("math", "TestGoLeanAsinh", overlay+"bridge_test.go", overlay+"asinh_test.go"))
	if lostSign == 0 || overflowingDouble == 0 || linear == 0 {
		log.Fatal("controls did not distinguish all three faults")
	}
	fmt.Printf("Private asinh passed %d native and %d Lean profile fixtures; losing signs failed %d, doubling large finite magnitudes failed %d, substituting the identity failed %d\n",
		len(tests), len(rows), lostSign, overflowingDouble, linear)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
