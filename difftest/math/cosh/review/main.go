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

func numerical(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if math.IsInf(x, 0) {
		return math.Inf(1)
	}
	// Direct even power series, without Exp or range reduction. The omitted
	// terms after degree 3200 are negligible at this precision for |x| ≤ 710.
	const precision = 384
	abs := new(big.Float).SetPrec(precision).SetFloat64(math.Abs(x))
	square := new(big.Float).SetPrec(precision).Mul(abs, abs)
	term := new(big.Float).SetPrec(precision).SetInt64(1)
	sum := new(big.Float).Copy(term)
	divisor := new(big.Float).SetPrec(precision)
	for k := int64(2); k <= 3200; k += 2 {
		term.Mul(term, square)
		term.Quo(term, divisor.SetInt64(k*(k-1)))
		sum.Add(sum, term)
	}
	out, _ := sum.Float64()
	return out
}

func agrees(x, got, want float64, arm64 bool) bool {
	if math.IsNaN(want) {
		return math.IsNaN(got)
	}
	magnitude := math.Abs(x)
	if math.IsInf(x, 0) || magnitude == 710 || (!arm64 && magnitude == 709.75) {
		return math.IsInf(got, 1)
	}
	if x == 0 || magnitude == math.SmallestNonzeroFloat64 {
		return math.Float64bits(got) == math.Float64bits(want)
	}
	a, b := math.Float64bits(got), math.Float64bits(want)
	return !math.IsNaN(got) && !math.IsInf(got, 0) && max(a, b)-min(a, b) <= 5
}

func main() {
	inputs := []float64{0, math.Copysign(0, -1), math.SmallestNonzeroFloat64, -math.SmallestNonzeroFloat64,
		math.Inf(1), math.Inf(-1), math.NaN(), math.Float64frombits(0xfff0000000000001)}
	for _, x := range []float64{0.125, 0.5, 1, 2, math.Nextafter(21, 20), 21, math.Nextafter(21, 22),
		100, 709.75, 710} {
		inputs = append(inputs, x, -x)
	}
	const profiles = 3
	var input strings.Builder
	var forcedSign, subtractReciprocal, scaledOverflow int
	nativeARM := runtime.GOARCH == "arm64"
	references := make([]float64, len(inputs))
	for i, x := range inputs {
		want := numerical(x)
		references[i] = want
		got := math.Cosh(x)
		if !agrees(x, got, want, nativeARM) {
			log.Fatalf("native Cosh(%g)=%g, numerical reference=%g", x, got, want)
		}
		if !math.IsNaN(x) && math.Float64bits(got) != math.Float64bits(math.Cosh(-x)) {
			log.Fatalf("native Cosh(%g) failed sign symmetry", x)
		}
		if !agrees(x, math.Copysign(got, x), want, nativeARM) {
			forcedSign++
		}
		magnitude := math.Abs(x)
		if magnitude <= 21 {
			ex := math.Exp(magnitude)
			if !agrees(x, 0.5*(ex-1/ex), want, nativeARM) {
				subtractReciprocal++
			}
		}
		if magnitude == 710 {
			if math.IsInf(want, 0) {
				log.Fatalf("expected reference at %g to remain finite", x)
			}
			if !agrees(x, want, want, nativeARM) {
				scaledOverflow++
			}
		}
		for profile := range profiles {
			fmt.Fprintf(&input, "%d %d\n", profile, math.Float64bits(x))
		}
	}
	rows := check.LeanRows("godiff-math-cosh", profiles*len(inputs), strings.NewReader(input.String()))
	for i, row := range rows {
		profile := i % profiles
		x, want := inputs[i/profiles], references[i/profiles]
		got := math.NaN()
		if row != "NaN" {
			word, err := strconv.ParseUint(row, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			got = math.Float64frombits(word)
		}
		if !agrees(x, got, want, profile == 2) {
			log.Fatalf("Lean profile %d Cosh(%g)=%g, numerical reference=%g", profile, x, got, want)
		}
		if i/profiles%2 == 1 && row != rows[i-profiles] {
			log.Fatalf("Lean profile %d Cosh(%g) failed sign symmetry", profile, x)
		}
	}
	if forcedSign != 13 || subtractReciprocal != 12 || scaledOverflow != 2 {
		log.Fatalf("unexpected discrimination: forced sign=%d, subtraction=%d, scaled overflow=%d",
			forcedSign, subtractReciprocal, scaledOverflow)
	}
	fmt.Printf("Cosh: %d native and %d Lean profile rows; forced sign=%d, subtraction=%d, scaled overflow=%d\n",
		len(inputs), len(rows), forcedSign, subtractReciprocal, scaledOverflow)
}
