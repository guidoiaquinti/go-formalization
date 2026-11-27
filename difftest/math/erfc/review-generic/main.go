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

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	x     float64
	want  uint64
	exact bool
}

func agrees(test fixture, got float64) bool {
	want := math.Float64frombits(test.want)
	if math.IsNaN(want) {
		return math.IsNaN(got)
	}
	if math.IsNaN(got) || math.Signbit(got) {
		return false
	}
	word := math.Float64bits(got)
	if test.exact {
		return word == test.want
	}
	return max(word, test.want)-min(word, test.want) <= 6
}

func observed(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(math.Float64bits(x), 10)
}

func main() {
	// The round-to-nearest golden words use a 384-bit positive-term
	// central series and alternating integration-by-parts tail expansion.
	tests := []fixture{
		{0, 0x3ff0000000000000, true},
		{math.Copysign(0, -1), 0x3ff0000000000000, true},
		{math.Inf(1), 0, true},
		{math.Inf(-1), 0x4000000000000000, true},
		{math.NaN(), 0x7ff8000000000000, true},
		{math.Float64frombits(0xfff0000000000001), 0x7ff8000000000000, true},
		{0x1.fffffffffffffp-57, 0x3ff0000000000000, false},
		{-0x1.fffffffffffffp-57, 0x3ff0000000000000, false},
		{0x1p-56, 0x3ff0000000000000, false},
		{-0x1p-56, 0x3ff0000000000000, false},
		{0x1.fffffffffffffp-03, 0x3fe728558ee694fc, false},
		{-0x1.fffffffffffffp-03, 0x3ff46bd5388cb582, false},
		{0x1p-02, 0x3fe728558ee694fc, false},
		{-0x1p-02, 0x3ff46bd5388cb582, false},
		{0x1.affffffffffffp-01, 0x3fcdcb8cae2d7481, false},
		{-0x1.affffffffffffp-01, 0x3ffc468e6a3a5170, false},
		{0x1.bp-01, 0x3fcdcb8cae2d747f, false},
		{-0x1.bp-01, 0x3ffc468e6a3a5170, false},
		{0x1.3ffffffffffffp+00, 0x3fb3bcd133aa1000, false},
		{-0x1.3ffffffffffffp+00, 0x3ffec432ecc55f00, false},
		{0x1.4p+00, 0x3fb3bcd133aa0ffc, false},
		{-0x1.4p+00, 0x3ffec432ecc55f00, false},
		{0x1.6db6db6db6db6p+01, 0x3f0bf37512335b03, false},
		{-0x1.6db6db6db6db6p+01, 0x3fffffc81915db99, false},
		{0x1.6db6db6db6db7p+01, 0x3f0bf37512335aee, false},
		{-0x1.6db6db6db6db7p+01, 0x3fffffc81915db99, false},
		{0x1.8p+02, 0x3c78cf81557d20b6, false},
		{-0x1.8p+02, 0x4000000000000000, true},
		{0x1p+03, 0x39ec74fc41217dfb, false},
		{-0x1p+03, 0x4000000000000000, false},
		{0x1.a8ccccccccccdp+04, 0x000b2ee03853bf84, false},
		{-0x1.a8ccccccccccdp+04, 0x4000000000000000, false},
		{0x1.bp+04, 0x0000000000019e0f, false},
		{-0x1.bp+04, 0x4000000000000000, false},
		{0x1.b19999999999ap+04, 0x00000000000001d8, false},
		{-0x1.b19999999999ap+04, 0x4000000000000000, false},
		{0x1.b333333333333p+04, 0x0000000000000002, true},
		{-0x1.b333333333333p+04, 0x4000000000000000, true},
		{0x1.b39999999999ap+04, 0x0000000000000001, true},
		{-0x1.b39999999999ap+04, 0x4000000000000000, true},
		{0x1.b3a1cac083127p+04, 0x0000000000000000, true},
		{-0x1.b3a1cac083127p+04, 0x4000000000000000, true},
		{0x1.cp+04, 0x0000000000000000, true},
		{-0x1.cp+04, 0x4000000000000000, true},
	}
	const profiles = 5
	var input, nativeFixture bytes.Buffer
	for _, test := range tests {
		for profile := range profiles {
			fmt.Fprintln(&input, profile, math.Float64bits(test.x))
		}
	}
	rows := check.LeanRows("godiff-math-erfc", profiles*len(tests), &input)
	var lostSign, earlyZero, fromErf int
	for i, test := range tests {
		for profile := range profiles {
			row := rows[i*profiles+profile]
			got := math.NaN()
			if row != "NaN" {
				word, err := strconv.ParseUint(row, 10, 64)
				must(err)
				got = math.Float64frombits(word)
			}
			if !agrees(test, got) {
				log.Fatalf("Lean profile %d erfc(%x): got=%#016x want=%#016x", profile,
					test.x, math.Float64bits(got), test.want)
			}
		}
		native := observed(math.Erfc(test.x))
		if runtime.GOARCH == "arm64" {
			native = rows[i*profiles+4]
		}
		kind := "numeric"
		if test.exact {
			kind = "exact"
		}
		fmt.Fprintln(&nativeFixture, math.Float64bits(test.x), native, observed(math.Float64frombits(test.want)), kind)
		if !agrees(test, math.Erfc(math.Abs(test.x))) {
			lostSign++
		}
		if test.x >= 26.5 && !agrees(test, 0) {
			earlyZero++
		}
		if !agrees(test, 1-math.Erf(test.x)) {
			fromErf++
		}
	}
	dir, err := os.MkdirTemp("", "golean-private-erfc-review-")
	must(err)
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	must(os.WriteFile(path, nativeFixture.Bytes(), 0600))
	must(os.Setenv("GOLEAN_MATH_ERFC_CASES", path))
	const overlay = "difftest/math/erfc/testdata/"
	must(check.TestInPackage("math", "TestGoLeanErfc", overlay+"bridge_test.go", overlay+"erfc_test.go"))
	if lostSign != 18 || earlyZero != 5 || fromErf < 7 {
		log.Fatalf("unexpected controls: lost sign=%d, early zero=%d, from Erf=%d", lostSign, earlyZero, fromErf)
	}
	fmt.Printf("private erfc: %d native and %d Lean profile rows; lost sign=%d, early zero=%d, from Erf=%d\n",
		len(tests), len(rows), lostSign, earlyZero, fromErf)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
