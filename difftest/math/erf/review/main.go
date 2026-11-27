package main

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const precision = 384

func precise(value float64) *big.Float {
	return new(big.Float).SetPrec(precision).SetFloat64(value)
}

func numerical(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if math.Abs(x) >= 6 {
		return math.Copysign(1, x)
	}
	if x == 0 {
		return x
	}
	// For the normalized integral F = exp(x²) integral_0^x exp(-t²)dt,
	// F'=1+2xF gives a positive-term series. This avoids the cancellation
	// of the primary differential oracle's alternating expansion.
	abs := precise(math.Abs(x))
	square := precise(0).Mul(abs, abs)
	twiceSquare := precise(0).Mul(precise(2), square)
	integral, term := precise(0).Set(abs), precise(0).Set(abs)
	exp, expTerm := precise(1), precise(1)
	denominator := precise(0)
	for i := int64(1); i <= 400; i++ {
		term.Mul(term, twiceSquare)
		term.Quo(term, denominator.SetInt64(2*i+1))
		integral.Add(integral, term)
		expTerm.Mul(expTerm, square)
		expTerm.Quo(expTerm, denominator.SetInt64(i))
		exp.Add(exp, expTerm)
	}
	scale, _, err := big.ParseFloat(
		"1.128379167095512573896158903121545171688101258657997713688171443421284936882986828973487320404214726886056695812725",
		10, precision, big.ToNearestEven)
	if err != nil {
		log.Fatal(err)
	}
	integral.Quo(integral, exp)
	integral.Mul(integral, scale)
	result, _ := integral.Float64()
	return math.Copysign(result, x)
}

func agrees(x, got, want float64) bool {
	if math.IsNaN(want) {
		return math.IsNaN(got)
	}
	a, b := math.Float64bits(got), math.Float64bits(want)
	if x == 0 || math.Abs(x) >= 6 {
		return a == b
	}
	return !math.IsNaN(got) && math.Signbit(got) == math.Signbit(want) && max(a, b)-min(a, b) <= 4
}

func observed(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(math.Float64bits(x), 10)
}

func contraction(x float64, profile int) (float64, bool) {
	efx := math.Float64frombits(0x3fc06eba8214db69)
	efx8 := math.Float64frombits(0x3ff06eba8214db69)
	var value float64
	a := math.Abs(x)
	switch math.Float64bits(a) {
	case 0x007fffffffffffff:
		if profile == 4 {
			value = 0.125 * math.FMA(efx8, a, 8*a)
		} else if profile >= 2 {
			value = 0.125 * math.FMA(8, a, float64(efx8*a))
		} else {
			value = 0.125 * (8*a + float64(efx8*a))
		}
	case 0x0080000000000003, 0x3e2fffffffffffff:
		if profile >= 2 {
			value = math.FMA(efx, a, a)
		} else {
			value = a + float64(efx*a)
		}
	default:
		return 0, false
	}
	return math.Copysign(value, x), true
}

func main() {
	tests := []float64{0, math.Copysign(0, -1), math.Inf(1), math.Inf(-1), math.NaN(),
		math.Float64frombits(0xfff0000000000001)}
	for i := uint64(1); i <= 128; i++ {
		x := math.Float64frombits(i)
		tests = append(tests, x, -x)
	}
	for _, center := range []float64{0x1p-1015, 0x1p-28, 0.84375, 1.25, 20.0 / 7, 5.875, 6} {
		for _, x := range []float64{math.Nextafter(center, 0), center, math.Nextafter(center, math.Inf(1))} {
			tests = append(tests, x, -x)
		}
	}
	x := math.Float64frombits(0x0080000000000003)
	tests = append(tests, x, -x)
	const profiles = 5
	var input strings.Builder
	numericalResults := make([]float64, len(tests))
	var lostSign, droppedScaling, overEarly int
	for i, x := range tests {
		got, want := math.Erf(x), numerical(x)
		numericalResults[i] = want
		if !agrees(x, got, want) {
			log.Fatalf("Go Erf(%g)=%g, numerical reference=%g", x, got, want)
		}
		if !agrees(x, math.Abs(got), want) {
			lostSign++
		}
		if magnitude := math.Abs(x); 0 < magnitude && magnitude < 0x1p-1015 {
			const efx = 0.128379167095512586316
			unsplit := x + float64(efx*x)
			if math.Float64bits(unsplit) != math.Float64bits(got) {
				droppedScaling++
			}
		}
		// Saturation cannot be applied as early as the first tail branch.
		if magnitude := math.Abs(x); 1.25 <= magnitude && magnitude < 6 &&
			!agrees(x, math.Copysign(1, x), want) {
			overEarly++
		}
		for profile := range profiles {
			fmt.Fprintf(&input, "%d %d\n", profile, math.Float64bits(x))
		}
	}
	rows := check.LeanRows("godiff-math-erf", profiles*len(tests), strings.NewReader(input.String()))
	var contractionRows int
	for i, row := range rows {
		x, profile := tests[i/profiles], i%profiles
		got := math.NaN()
		if row != "NaN" {
			word, err := strconv.ParseUint(row, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			got = math.Float64frombits(word)
		}
		if !agrees(x, got, numericalResults[i/profiles]) {
			log.Fatalf("Lean profile %d Erf(%g)=%g, expected %g", profile, x, got, numericalResults[i/profiles])
		}
		if runtime.GOARCH == "arm64" && profile == 4 && row != observed(math.Erf(x)) {
			log.Fatalf("native ARM mismatch on Erf(%g): Lean=%s", x, row)
		}
		if exact, ok := contraction(x, profile); ok {
			if row != observed(exact) {
				log.Fatalf("profile %d contraction mismatch at %#016x: Lean=%s, expected=%s",
					profile, math.Float64bits(x), row, observed(exact))
			}
			contractionRows++
		}
	}
	if lostSign != 152 || droppedScaling < 16 || overEarly != 10 || contractionRows != 30 {
		log.Fatalf("unexpected controls: lost sign=%d, dropped tiny scaling=%d, premature saturation=%d, contraction rows=%d",
			lostSign, droppedScaling, overEarly, contractionRows)
	}
	fmt.Printf("Erf: %d native and %d Lean profile rows; lost sign=%d, dropped tiny scaling=%d, premature saturation=%d, contraction rows=%d\n",
		len(tests), len(rows), lostSign, droppedScaling, overEarly, contractionRows)
}
