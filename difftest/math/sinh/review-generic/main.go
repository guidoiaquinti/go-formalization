package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type observation struct {
	input, ideal float64
	exact        bool
}

func expected(x float64) observation {
	if math.IsNaN(x) || math.IsInf(x, 0) || x == 0 || math.Abs(x) == math.SmallestNonzeroFloat64 {
		return observation{x, x, true}
	}
	// Sum the unreduced exponential series at 384-bit precision. Even
	// near overflow, 3,000 terms leave the discarded tail far below
	// 384-bit rounding error. The reciprocal gives the negative Exp.
	const precision = 384
	magnitude := new(big.Float).SetPrec(precision).SetFloat64(math.Abs(x))
	term := new(big.Float).SetPrec(precision).SetInt64(1)
	sum := new(big.Float).Copy(term)
	divisor := new(big.Float).SetPrec(precision)
	for i := int64(1); i <= 3000; i++ {
		term.Mul(term, magnitude)
		term.Quo(term, divisor.SetInt64(i))
		sum.Add(sum, term)
	}
	sum.Sub(sum, divisor.Quo(divisor.SetInt64(1), sum))
	sum.SetMantExp(sum, -1)
	ideal, _ := sum.Float64()
	return observation{x, math.Copysign(ideal, x), false}
}

func stagedOverflow(x float64, arm64 bool) bool {
	abs := math.Abs(x)
	return abs == 710 || (!arm64 && abs == 709.75)
}

func agrees(test observation, got float64, arm64 bool) bool {
	if math.IsNaN(test.ideal) {
		return math.IsNaN(got)
	}
	if stagedOverflow(test.input, arm64) {
		return math.IsInf(got, int(math.Copysign(1, test.input)))
	}
	if test.exact {
		return math.Float64bits(got) == math.Float64bits(test.ideal)
	}
	a, b := math.Float64bits(got), math.Float64bits(test.ideal)
	return !math.IsInf(got, 0) && !math.IsNaN(got) && max(a, b)-min(a, b) <= 5
}

func bits(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(math.Float64bits(x), 10)
}

func decode(row string) float64 {
	if row == "NaN" {
		return math.NaN()
	}
	word, err := strconv.ParseUint(row, 10, 64)
	must(err)
	return math.Float64frombits(word)
}

func main() {
	values := []float64{0, math.Copysign(0, -1), math.SmallestNonzeroFloat64, -math.SmallestNonzeroFloat64,
		math.Inf(1), math.Inf(-1), math.NaN(), math.Float64frombits(0xfff0000000000001)}
	for _, x := range []float64{0.125, math.Nextafter(0.5, 0), 0.5, math.Nextafter(0.5, 1), 1, 2,
		math.Nextafter(21, 20), 21, math.Nextafter(21, 22), 709.75, 710} {
		values = append(values, x, -x)
	}
	const profiles = 5
	var input strings.Builder
	tests := make([]observation, len(values))
	for i, x := range values {
		tests[i] = expected(x)
		if stagedOverflow(x, false) && math.IsInf(tests[i].ideal, 0) {
			log.Fatalf("review reference at x=%g should remain finite", x)
		}
		for profile := range profiles {
			fmt.Fprintf(&input, "%d %d\n", profile, math.Float64bits(x))
		}
	}
	rows := check.LeanRows("godiff-math-sinh", profiles*len(tests), strings.NewReader(input.String()))
	var fixture bytes.Buffer
	var lostSign, noReciprocal, scaledOverflow int
	for i, test := range tests {
		for profile := range profiles {
			got := decode(rows[i*profiles+profile])
			if !agrees(test, got, profile == 4) {
				log.Fatalf("Lean profile %d private sinh(%g)=%g, mathematical reference=%g", profile, test.input, got, test.ideal)
			}
		}
		nativeARM := runtime.GOARCH == "arm64"
		observed := math.Sinh(test.input)
		if nativeARM {
			observed = decode(rows[i*profiles+4])
		}
		if !agrees(test, observed, nativeARM) {
			log.Fatalf("native private sinh(%g)=%g, mathematical reference=%g", test.input, observed, test.ideal)
		}
		kind := "numeric"
		if test.exact {
			kind = "exact"
		} else if stagedOverflow(test.input, nativeARM) {
			kind = "staged"
		}
		fmt.Fprintln(&fixture, math.Float64bits(test.input), bits(observed), bits(test.ideal), kind)
		if !agrees(test, math.Abs(observed), nativeARM) {
			lostSign++
		}
		magnitude := math.Abs(test.input)
		if 0.5 < magnitude && magnitude <= 21 &&
			!agrees(test, math.Copysign(0.5*math.Exp(magnitude), test.input), nativeARM) {
			noReciprocal++
		}
		if magnitude == 710 && !agrees(test, test.ideal, nativeARM) {
			scaledOverflow++
		}
	}
	dir, err := os.MkdirTemp("", "golean-private-sinh-review-")
	must(err)
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	must(os.WriteFile(path, fixture.Bytes(), 0600))
	must(os.Setenv("GOLEAN_MATH_SINH_CASES", path))
	const overlay = "difftest/math/sinh/testdata/"
	must(check.TestInPackage("math", "TestGoLeanSinh", overlay+"bridge_test.go", overlay+"sinh_test.go"))
	if lostSign != 14 || noReciprocal != 6 || scaledOverflow != 2 {
		log.Fatalf("unexpected discrimination: lost sign=%d, no reciprocal=%d, scaled overflow=%d",
			lostSign, noReciprocal, scaledOverflow)
	}
	fmt.Printf("private sinh: %d native and %d Lean profile rows; lost sign=%d, no reciprocal=%d, scaled overflow=%d\n",
		len(tests), len(rows), lostSign, noReciprocal, scaledOverflow)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
