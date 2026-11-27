package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/big"
	"runtime"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const precision = 384

func precise(x float64) *big.Float {
	return new(big.Float).SetPrec(precision).SetFloat64(x)
}

func golden(x float64, sqrtPi, euler *big.Float) float64 {
	switch {
	case math.IsNaN(x), math.IsInf(x, -1), x < 0 && math.Trunc(x) == x:
		return math.NaN()
	case x == 0, math.IsInf(x, 1):
		return math.Copysign(math.Inf(1), x)
	}
	var answer *big.Float
	switch {
	case math.Abs(x) <= 1.01e-9:
		// At these inputs the terms omitted from 1/x − Euler affect
		// the result by less than 2^-58 in relative magnitude.
		answer = precise(0).Quo(precise(1), precise(x))
		answer.Sub(answer, euler)
	case x-math.Floor(x) == 0.5:
		// Exact binary half-integer factors, reduced to Gamma(1/2).
		answer = precise(0).Set(sqrtPi)
		for z := 0.5; z < x; z += 1 {
			answer.Mul(answer, precise(z))
		}
		for z := -0.5; z >= x; z -= 1 {
			answer.Quo(answer, precise(z))
		}
	case x < 0:
		// These probes are adjacent to -33. The exact recurrence
		// reaches 1+δ, where replacing Gamma(1+δ) by one loses less
		// than 2^-46 in relative magnitude.
		z, denominator := precise(x), precise(1)
		for z.Cmp(precise(0.5)) < 0 {
			denominator.Mul(denominator, z)
			z.Add(z, precise(1))
		}
		answer = precise(0).Quo(precise(1), denominator)
	default:
		// The only positive noninteger probes here are adjacent to 33.
		// Comparing them with exact 32! differs by less than 2^-44.
		factorial := big.NewInt(1)
		for i := int64(2); i < int64(math.Round(x)); i++ {
			factorial.Mul(factorial, big.NewInt(i))
		}
		answer = precise(0).SetInt(factorial)
	}
	result, _ := answer.Float64()
	return result
}

func agrees(got, want float64) bool {
	if math.IsNaN(want) {
		return math.IsNaN(got)
	}
	if math.IsNaN(got) || math.Signbit(got) != math.Signbit(want) ||
		math.IsInf(got, 0) != math.IsInf(want, 0) {
		return false
	}
	if math.IsInf(want, 0) || want == 0 {
		return math.Float64bits(got) == math.Float64bits(want)
	}
	if math.Abs(want) < 0x1p-1022 {
		a, b := math.Float64bits(got), math.Float64bits(want)
		return max(a, b)-min(a, b) <= 128
	}
	return math.Abs(got/want-1) <= 0x1p-42
}

func main() {
	pi, _, err := big.ParseFloat(
		"3.141592653589793238462643383279502884197169399375105820974944592307816406286208998628034825342117068",
		10, precision, big.ToNearestEven)
	if err != nil {
		log.Fatal(err)
	}
	sqrtPi := precise(0).Sqrt(pi)
	euler, _, err := big.ParseFloat("0.57721566490153286060651209008240243104215933593992",
		10, precision, big.ToNearestEven)
	if err != nil {
		log.Fatal(err)
	}
	tests := []float64{0, math.Copysign(0, -1), math.Inf(1), math.Inf(-1), math.NaN(), -1, -2, -33, -200,
		math.SmallestNonzeroFloat64, -math.SmallestNonzeroFloat64, 1, 2, 3, 8, 33, 34, 171, 172}
	for _, x := range []float64{0.5, 1.5, 2.5, 32.5, 33.5, 142.5, 143.5, 171.5, 172.5,
		175.5, 179.5, 180.5, 181.5, 199.5, 200.5, math.Nextafter(33, 0),
		math.Nextafter(33, math.Inf(1)), math.Nextafter(1e-9, 0), 1e-9, math.Nextafter(1e-9, 1)} {
		tests = append(tests, x, -x)
	}
	const profiles = 5
	var model bytes.Buffer
	references := make([]float64, len(tests))
	var lostSign, noRecurrence, overflowedDenominator int
	for i, x := range tests {
		want := golden(x, sqrtPi, euler)
		references[i] = want
		got := math.Gamma(x)
		if !agrees(got, want) {
			log.Fatalf("Go Gamma(%x)=%x, recurrence reference=%x", x, got, want)
		}
		if !agrees(math.Abs(got), want) {
			lostSign++
		}
		if 3 <= x && x <= 33 && x == math.Trunc(x) && !agrees(1, want) {
			noRecurrence++
		}
		if x < -172 && got != 0 && !agrees(math.Copysign(0, got), want) {
			overflowedDenominator++
		}
		for profile := range profiles {
			fmt.Fprintln(&model, profile, math.Float64bits(x))
		}
	}
	rows := check.LeanRows("godiff-math-gamma", profiles*len(tests), &model)
	for i, row := range rows {
		profile, x := i%profiles, tests[i/profiles]
		got := math.NaN()
		if row != "NaN" {
			bits, err := strconv.ParseUint(row, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			got = math.Float64frombits(bits)
		}
		if !agrees(got, references[i/profiles]) {
			log.Fatalf("Lean profile %d Gamma(%x)=%x, recurrence reference=%x",
				profile, x, got, references[i/profiles])
		}
		native := math.Gamma(x)
		if runtime.GOARCH == "arm64" && profile == 4 && !math.IsNaN(native) &&
			math.Float64bits(got) != math.Float64bits(native) {
			log.Fatalf("arm64 Gamma(%x): Go=%x Lean=%x", x, native, got)
		}
	}
	if lostSign != 13 || noRecurrence != 3 || overflowedDenominator != 4 {
		log.Fatalf("mutants: sign=%d, recurrence=%d, overflowed denominator=%d",
			lostSign, noRecurrence, overflowedDenominator)
	}
	fmt.Printf("Gamma: %d native and %d Lean profile rows; sign=%d, recurrence=%d, overflowed denominator=%d mutant failures\n",
		len(tests), len(rows), lostSign, noRecurrence, overflowedDenominator)
}
