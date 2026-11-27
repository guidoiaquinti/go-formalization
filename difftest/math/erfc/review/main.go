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

func precise(value float64) *big.Float {
	return new(big.Float).SetPrec(precision).SetFloat64(value)
}

func numerical(x float64, scale *big.Float) float64 {
	magnitude := math.Abs(x)
	if math.IsNaN(x) {
		return math.NaN()
	}
	if magnitude >= 28 {
		if math.Signbit(x) {
			return 2
		}
		return 0
	}
	abs := precise(magnitude)
	square := precise(0).Mul(abs, abs)
	term, exp, divisor := precise(1), precise(1), precise(0)
	for n := int64(1); n <= 3200; n++ {
		term.Mul(term, square)
		term.Quo(term, divisor.SetInt64(n))
		exp.Add(exp, term)
	}
	var value *big.Float
	if magnitude < 8 {
		// Build the positive-term series for the scaled Gaussian integral.
		// Division by exp(x²) and subtraction occur at high precision.
		twiceSquare := precise(0).Mul(precise(2), square)
		term.Set(abs)
		value = precise(0).Set(term)
		for n := int64(1); n <= 400; n++ {
			term.Mul(term, twiceSquare)
			term.Quo(term, divisor.SetInt64(2*n+1))
			value.Add(value, term)
		}
		value.Quo(value, exp)
		value.Mul(value, scale)
		value.Sub(precise(1), value)
	} else {
		// Repeated integration by parts gives an alternating expansion
		// for the scaled tail. For x >= 8 its 49th omitted term bounds
		// the relative truncation error below 2^-85.
		inverse := precise(0).Quo(precise(-0.5), square)
		term.SetInt64(1)
		value = precise(1)
		coefficient := precise(0)
		for n := int64(1); n <= 48; n++ {
			term.Mul(term, inverse)
			term.Mul(term, coefficient.SetInt64(2*n-1))
			value.Add(value, term)
		}
		value.Quo(value, abs)
		value.Quo(value, exp)
		value.Mul(value, scale)
		value.SetMantExp(value, -1)
	}
	if math.Signbit(x) {
		value.Sub(precise(2), value)
	}
	out, _ := value.Float64()
	return out
}

func agrees(x, got, want float64) bool {
	if math.IsNaN(want) {
		return math.IsNaN(got)
	}
	a, b := math.Float64bits(got), math.Float64bits(want)
	if x == 0 || math.Abs(x) >= 28 || x == 27.2 || x == 27.225 || x == 27.227 {
		return a == b
	}
	return !math.IsNaN(got) && !math.Signbit(got) && max(a, b)-min(a, b) <= 6
}

func main() {
	scale, _, err := big.ParseFloat(
		"1.128379167095512573896158903121545171688101258657997713688171443421284936882986828973487320404214726886056695812725",
		10, precision, big.ToNearestEven)
	if err != nil {
		log.Fatal(err)
	}
	tests := []float64{0, math.Copysign(0, -1), math.Inf(1), math.Inf(-1), math.NaN(),
		math.Float64frombits(0xfff0000000000001)}
	for _, center := range []float64{0x1p-56, 0.25, 0.84375, 1.25, 20.0 / 7, 6, 8, 26.543,
		27, 27.1, 27.2, 27.225, 27.227, 27.23, 28} {
		for _, x := range []float64{math.Nextafter(center, 0), center, math.Nextafter(center, math.Inf(1))} {
			tests = append(tests, x, -x)
		}
	}
	const profiles = 5
	var modelInput bytes.Buffer
	references := make([]float64, len(tests))
	var wrongSign, flush, viaErf int
	for i, x := range tests {
		got, want := math.Erfc(x), numerical(x, scale)
		references[i] = want
		if !agrees(x, got, want) {
			log.Fatalf("Go Erfc(%g)=%#016x, reference=%#016x", x, math.Float64bits(got), math.Float64bits(want))
		}
		if !agrees(x, math.Erfc(math.Abs(x)), want) {
			wrongSign++
		}
		if 0 < got && got < 0x1p-1022 && !agrees(x, 0, want) {
			flush++
		}
		if !agrees(x, 1-math.Erf(x), want) {
			viaErf++
		}
		for profile := range profiles {
			fmt.Fprintf(&modelInput, "%d %d\n", profile, math.Float64bits(x))
		}
	}
	rows := check.LeanRows("godiff-math-erfc", profiles*len(tests), &modelInput)
	for i, row := range rows {
		profile, x := i%profiles, tests[i/profiles]
		var got float64
		if row == "NaN" {
			got = math.NaN()
		} else {
			word, err := strconv.ParseUint(row, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			got = math.Float64frombits(word)
		}
		if !agrees(x, got, references[i/profiles]) {
			log.Fatalf("Lean profile %d Erfc(%g)=%#016x, reference=%#016x", profile, x,
				math.Float64bits(got), math.Float64bits(references[i/profiles]))
		}
		native := math.Erfc(x)
		if runtime.GOARCH == "arm64" && profile == 4 && !math.IsNaN(native) &&
			math.Float64bits(got) != math.Float64bits(native) {
			log.Fatalf("Erfc(%g) arm64 Lean=%#016x, Go=%#016x", x, math.Float64bits(got), math.Float64bits(native))
		}
	}
	if wrongSign != 43 || flush != 8 || viaErf != 20 {
		log.Fatalf("unexpected mutant failures: wrong sign=%d, flushed subnormals=%d, via Erf=%d",
			wrongSign, flush, viaErf)
	}
	fmt.Printf("Erfc: %d native and %d Lean profile rows; wrong sign=%d, flushed subnormals=%d, via Erf=%d\n",
		len(tests), len(rows), wrongSign, flush, viaErf)
}
