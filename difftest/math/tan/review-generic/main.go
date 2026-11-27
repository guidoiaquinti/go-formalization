package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	x, want, tolerance float64
}

func agrees(got float64, test fixture) bool {
	if math.IsNaN(test.want) {
		return math.IsNaN(got)
	}
	if test.want == 0 {
		return math.Float64bits(got) == math.Float64bits(test.want)
	}
	return math.Abs(got-test.want) <= test.tolerance*math.Abs(test.want)
}

func main() {
	negativeZero := math.Copysign(0, -1)
	nearPole := math.Nextafter(math.Pi/2, math.Inf(1))
	tests := []fixture{
		{0, 0, 0},
		{negativeZero, negativeZero, 0},
		{math.NaN(), math.NaN(), 0},
		{math.Inf(1), math.NaN(), 0},
		{math.Inf(-1), math.NaN(), 0},
		{math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, 0},
		{math.Pi / 4, 1, 1e-14},
		{3 * math.Pi / 4, -1, 1e-14},
		{-math.Pi / 4, -1, 1e-14},
		{math.Pi / 2, 16331239353195370, 1e-13},
		{nearPole, -6218431163823738, 1e-13},
		{-nearPole, 6218431163823738, 1e-13},
		// Cross-checks from the platform libm exercise Payne-Hanek reduction.
		{0x1p29, -0.34551069233430393, 1e-13},
		{0x1p128, 4.3759050220192828, 1e-13},
		{0x1p512, 7.2334072829059277, 1e-13},
		{math.MaxFloat64, -0.0049620158744448943, 1e-13},
	}
	var input strings.Builder
	for _, test := range tests {
		for profile := range 3 {
			fmt.Fprintf(&input, "%d %d\n", profile, math.Float64bits(test.x))
		}
	}
	results := check.LeanRows("godiff-math-tan", 3*len(tests), strings.NewReader(input.String()))
	profile := nativeProfile()
	var rows bytes.Buffer
	var modelFailures, dropOddReciprocal, loseNegativeSign int
	for i, test := range tests {
		for profile := range 3 {
			got := math.NaN()
			token := results[3*i+profile]
			if token != "NaN" {
				word, err := strconv.ParseUint(token, 10, 64)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				got = math.Float64frombits(word)
			}
			if !agrees(got, test) {
				fmt.Fprintf(os.Stderr, "Lean profile %d Tan(%g) = %g, wanted %g\n", profile, test.x, got, test.want)
				modelFailures++
			}
		}
		fmt.Fprintf(&rows, "%d %s 0 true\n", math.Float64bits(test.x), results[3*i+profile])
		reduced := math.Remainder(test.x, math.Pi/2)
		if !agrees(math.Tan(reduced), test) {
			dropOddReciprocal++
		}
		if !agrees(math.Tan(math.Abs(test.x)), test) {
			loseNegativeSign++
		}
	}
	dir, err := os.MkdirTemp("", "golean-private-tan-review-")
	must(err)
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	must(os.WriteFile(path, rows.Bytes(), 0600))
	must(os.Setenv("GOLEAN_MATH_TAN_CASES", path))
	const overlay = "difftest/math/tan/testdata/"
	must(check.TestInPackage("math", "TestGoLeanTan", overlay+"bridge_test.go", overlay+"tan_test.go"))
	fmt.Printf("Substituting rounded remainder for the odd-quadrant reciprocal disagreed on %d of %d fixtures\n", dropOddReciprocal, len(tests))
	fmt.Printf("Losing the input sign disagreed on %d of %d fixtures\n", loseNegativeSign, len(tests))
	fmt.Printf("Lean disagreed on %d of %d profile fixtures\n", modelFailures, len(results))
	if modelFailures > 0 || dropOddReciprocal == 0 || loseNegativeSign == 0 {
		os.Exit(1)
	}
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func nativeProfile() int {
	switch runtime.GOARCH {
	case "arm64":
		return 2
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			panic("cannot determine Tan build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		panic("unreviewed Tan architecture " + runtime.GOARCH)
	}
}
