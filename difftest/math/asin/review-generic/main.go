package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
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
	return math.Abs(got-test.want) <= 0x1p-51
}

func main() {
	negativeZero := math.Copysign(0, -1)
	// Finite golden values were obtained independently with platform libm.
	tests := []fixture{
		{0.6, 0x1.4978fa3269ee1p-1, false}, {0.8, 0x1.dac670561bb50p-1, false},
		{-0.6, -0x1.4978fa3269ee1p-1, false}, {-0.8, -0x1.dac670561bb50p-1, false},
		{0x1.6a09e667f3bcdp-1, 0x1.921fb54442d19p-1, false},
		{0.03125, 0x1.000aabde0b9c8p-5, false}, {-0.25, -0x1.02be9ce0b87cdp-2, false},
		{0.625, 0x1.59aad71ced00fp-1, false}, {0.9, 0x1.1ea93705fa172p+0, false},
		{-0.95, -0x1.40d41159f3409p+0, false}, {0.99, 0x1.6de3c6f33d51dp+0, false},
		{0x1.fffffffffffffp-1, 0x1.921fb50442d18p+0, false},
		{0, 0, true}, {negativeZero, negativeZero, true},
		{1, math.Pi / 2, true}, {-1, -math.Pi / 2, true},
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
		got := math.Asin(test.input)
		if !agrees(got, test) {
			log.Fatalf("Go Asin(%g) = %.17g, want %.17g", test.input, got, test.want)
		}
		observation := "NaN"
		if !math.IsNaN(got) {
			observation = strconv.FormatUint(math.Float64bits(got), 10)
		}
		kind := "numeric"
		if test.exact {
			kind = "exact"
		} else if runtime.GOARCH == "amd64" && math.Abs(test.input) > 0.7 {
			kind = "upper"
		}
		fmt.Fprintln(&cases, math.Float64bits(test.input), observation, math.Float64bits(test.want), kind)
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
			log.Fatalf("Lean profile %d Asin(%g) = %.17g, want %.17g", i%profiles, test.input, got, test.want)
		}
	}
	for i := 0; i < 4; i += 2 {
		for profile := range profiles {
			got := values[i*profiles+profile] + values[(i+1)*profiles+profile]
			if math.Abs(got-math.Copysign(math.Pi/2, tests[i].input)) > 0x1p-50 {
				log.Fatalf("Lean profile %d failed the unit-circle complement identity", profile)
			}
		}
	}
	dir, err := os.MkdirTemp("", "golean-private-asin-review-")
	must(err)
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	must(os.WriteFile(path, cases.Bytes(), 0600))
	must(os.Setenv("GOLEAN_MATH_ASIN_CASES", path))
	const overlay = "difftest/math/asin/testdata/"
	must(check.TestInPackage("math", "TestGoLeanAsin", overlay+"bridge_test.go", overlay+"asin_test.go"))
	if lostSign == 0 || clampedDomain == 0 || linear == 0 {
		log.Fatal("controls did not distinguish all three faults")
	}
	fmt.Printf("Private asin passed %d native and %d Lean profile fixtures; losing signs failed %d, clamping out-of-domain inputs failed %d, substituting the identity failed %d\n",
		len(tests), len(rows), lostSign, clampedDomain, linear)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
