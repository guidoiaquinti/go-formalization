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

func agrees(got float64, test fixture, separatelyRounded bool) bool {
	if math.IsNaN(test.want) {
		return math.IsNaN(got)
	}
	if test.exact {
		return math.Float64bits(got) == math.Float64bits(test.want)
	}
	if separatelyRounded && math.Abs(test.input) > 0.7 {
		return math.Abs(got-test.want) <= 0x1p-38
	}
	return math.Abs(got-test.want) <= 0x1p-49
}

func main() {
	negativeZero := math.Copysign(0, -1)
	// Finite golden values were obtained independently with platform libm.
	tests := []fixture{
		{0.25, 0x1.51700e0c14b25p+0, false}, {-0.25, 0x1.d2cf5c7c70f0bp+0, false},
		{0.6, 0x1.dac670561bb50p-1, false}, {-0.6, 0x1.1b6e192ebbe44p+1, false},
		{0.8, 0x1.4978fa3269ee0p-1, false}, {-0.8, 0x1.3fc176b7a8560p+1, false},
		{0x1.fffffff800000p-1, 0x1.6a09e6686c6acp-15, false},
		{-0x1.fffffff800000p-1, 0x1.921e4b3a5c691p+1, false},
		{0.95, 0x1.452e8fa93e43cp-2, false},
		{0x1.fffffffffe000p-1, 0x1.6a09e667f3dafp-20, false},
		{-0x1.ffffffffffffep-1, 0x1.921fb5170194bp+1, false},
		{0, math.Pi / 2, true}, {negativeZero, math.Pi / 2, true},
		{1, 0, true}, {-1, math.Pi, true},
		{2, math.NaN(), true}, {-2, math.NaN(), true},
		{math.Inf(1), math.NaN(), true}, {math.Inf(-1), math.NaN(), true},
		{math.NaN(), math.NaN(), true},
	}
	const profiles = 3
	var input strings.Builder
	var cases bytes.Buffer
	separateNative := runtime.GOARCH == "amd64"
	var reversedSign, clampedDomain int
	for _, test := range tests {
		for profile := range profiles {
			fmt.Fprintf(&input, "%d %d\n", profile, math.Float64bits(test.input))
		}
		got := math.Acos(test.input)
		if !agrees(got, test, separateNative) {
			log.Fatalf("Go Acos(%g) = %.17g, want %.17g", test.input, got, test.want)
		}
		observation := "NaN"
		if !math.IsNaN(got) {
			observation = strconv.FormatUint(math.Float64bits(got), 10)
		}
		kind := "numeric"
		if test.exact {
			kind = "exact"
		} else if separateNative && math.Abs(test.input) > 0.7 {
			kind = "upper"
		}
		fmt.Fprintln(&cases, math.Float64bits(test.input), observation, math.Float64bits(test.want), kind)
		if !agrees(math.Pi-got, test, separateNative) {
			reversedSign++
		}
		if math.Abs(test.input) > 1 && !agrees(math.Pi, test, separateNative) {
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
		if !agrees(got, test, i%profiles != 2) {
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
	dir, err := os.MkdirTemp("", "golean-private-acos-review-")
	must(err)
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	must(os.WriteFile(path, cases.Bytes(), 0600))
	must(os.Setenv("GOLEAN_MATH_ACOS_CASES", path))
	const overlay = "difftest/math/acos/testdata/"
	must(check.TestInPackage("math", "TestGoLeanAcos", overlay+"bridge_test.go", overlay+"acos_test.go"))
	if reversedSign == 0 || clampedDomain == 0 {
		log.Fatal("controls did not distinguish both faults")
	}
	fmt.Printf("Private acos passed %d native and %d Lean profile fixtures; reversing the Asin term failed %d, clamping out-of-domain inputs failed %d\n",
		len(tests), len(rows), reversedSign, clampedDomain)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
