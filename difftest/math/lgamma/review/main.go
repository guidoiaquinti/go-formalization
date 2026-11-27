package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

// Normalize before taking the logarithm so the amd64 public Log's
// behavior on an exponent-zero encoding cannot influence this oracle.
func logarithm(x float64, exponentAdjust int) float64 {
	mantissa, exponent := math.Frexp(x)
	return math.FMA(float64(exponent-1+exponentAdjust), math.Ln2, math.Log(2*mantissa))
}

func golden(x float64, arm64 bool) (float64, int) {
	switch {
	case math.IsNaN(x), math.IsInf(x, 0):
		return x, 1
	case x != 2 && math.Abs(x-2) < 1e-15:
		return (x - 2) * 0.42278433509846713, 1
	case x == 0, x < 0 && math.Trunc(x) == x:
		return math.Inf(1), 1
	}
	a := math.Abs(x)
	if a < 0x1p-69 {
		var result float64
		if !arm64 && a < 0x1p-1022 {
			// The pinned amd64 Log treats the input as having the
			// implicit leading bit at an exponent one below normal.
			withLeading := math.Float64frombits(math.Float64bits(a) | 1<<52)
			result = -logarithm(withLeading, -1)
		} else {
			result = -logarithm(a, 0)
		}
		if x < 0 {
			return result, -1
		}
		return result, 1
	}
	if a >= 0x1p+40 {
		// At this scale all remaining Stirling corrections are well
		// below the comparison tolerance, even close to overflow.
		logAbs := math.FMA(a, logarithm(a, 0)-1, 0)
		if x > 0 {
			return logAbs, 1
		}
		sign := -1
		if int64(a)&1 != 0 {
			sign = 1
		}
		return math.Log(math.Pi) - logarithm(a, 0) - logAbs, sign
	}
	if a-math.Floor(a) == 0.5 {
		logAbs, sign := 0.5*math.Log(math.Pi), 1
		for k := 0.5; k < x; k++ {
			logAbs += logarithm(k, 0)
		}
		for k := -0.5; k >= x; k-- {
			logAbs -= logarithm(-k, 0)
			sign = -sign
		}
		return logAbs, sign
	}
	if x < 0 {
		// The interior negative integer-neighbor probes reduce by
		// recurrence to 1+δ with |δ| <= 2^-49. Its log Gamma is
		// negligible compared with the selected tolerance.
		value, sign := 0.0, 1
		for k := x; k < 0.5; k++ {
			value -= logarithm(math.Abs(k), 0)
			if k < 0 {
				sign = -sign
			}
		}
		return value, sign
	}
	// Integer values and the immediate floating neighbors of 2 and 8.
	// The neighbors differ from the nearest factorial by less than 2^-45.
	var value float64
	for i := 2; i < int(math.Round(x)); i++ {
		value += logarithm(float64(i), 0)
	}
	return value, 1
}

func agrees(got float64, sign int, want float64, wantSign int) bool {
	if sign != wantSign {
		return false
	}
	if math.IsNaN(want) {
		return math.IsNaN(got)
	}
	if math.IsInf(want, 0) || want == 0 {
		return math.Float64bits(got) == math.Float64bits(want)
	}
	return !math.IsNaN(got) && math.Abs(got-want)/max(1, math.Abs(want)) <= 0x1p-42
}

func main() {
	tests := []float64{0, math.Copysign(0, -1), math.Inf(1), math.Inf(-1), math.NaN(),
		math.Float64frombits(0xfff0000000000001), -1, -2, -33, -0x1p52, -math.MaxFloat64,
		1, 2, 3, 7, 8, 20, 171, 172, 0x1p58, 1e250, 1e305, math.MaxFloat64,
		math.Nextafter(0x1p58, 0), math.Nextafter(0x1p58, math.Inf(1)),
		math.Nextafter(2, 0), math.Nextafter(2, 3), math.Nextafter(8, 0), math.Nextafter(8, 9),
		math.Nextafter(-8, -7), math.Nextafter(-8, -9), math.Nextafter(-2, -1), math.Nextafter(-2, -3),
		-math.Nextafter(0x1p52, 0)}
	for _, x := range []float64{math.SmallestNonzeroFloat64, math.Nextafter(0x1p-1022, 0), 0x1p-1022,
		0x1p-70, 0.5, 1.5, 2.5, 7.5, 8.5, 171.5, 172.5, 200.5} {
		tests = append(tests, x, -x)
	}
	const profiles = 3
	var model bytes.Buffer
	var alwaysPositive, roundedGamma, largeInfinite, subnormalStage int
	for _, x := range tests {
		nativeARM := runtime.GOARCH == "arm64"
		want, wantSign := golden(x, nativeARM)
		got, sign := math.Lgamma(x)
		if !agrees(got, sign, want, wantSign) {
			log.Fatalf("Go Lgamma(%x)=(%x,%d), want=(%x,%d)", x, got, sign, want, wantSign)
		}
		if !agrees(got, 1, want, wantSign) {
			alwaysPositive++
		}
		if math.Abs(x) > 171 && !agrees(math.Log(math.Abs(math.Gamma(x))), sign, want, wantSign) {
			roundedGamma++
		}
		if x >= 0x1p58 && !math.IsInf(want, 0) && !agrees(math.Inf(1), sign, want, wantSign) {
			largeInfinite++
		}
		if a := math.Abs(x); 0 < a && a < 0x1p-1022 {
			modelExpected, modelSign := golden(x, false)
			ideal, idealSign := golden(x, true)
			if modelSign == idealSign && !agrees(ideal, idealSign, modelExpected, modelSign) {
				subnormalStage++
			}
		}
		for profile := range profiles {
			fmt.Fprintf(&model, "%d %d\n", profile, math.Float64bits(x))
		}
	}
	rows := check.LeanRows("godiff-math-lgamma", profiles*len(tests), &model)
	for i, row := range rows {
		x, profile := tests[i/profiles], i%profiles
		word, signField, ok := strings.Cut(row, " ")
		if !ok {
			log.Fatalf("invalid Lean observation %q", row)
		}
		sign, err := strconv.Atoi(signField)
		if err != nil {
			log.Fatal(err)
		}
		got := math.NaN()
		if word != "NaN" {
			bits, err := strconv.ParseUint(word, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			got = math.Float64frombits(bits)
		}
		want, wantSign := golden(x, profile == 2)
		if !agrees(got, sign, want, wantSign) {
			log.Fatalf("Lean profile %d Lgamma(%x)=(%x,%d), want=(%x,%d)", profile, x, got, sign, want, wantSign)
		}
	}
	if alwaysPositive < 8 || roundedGamma < 8 || largeInfinite != 4 || subnormalStage != 2 {
		log.Fatalf("controls: sign=%d, rounded gamma=%d, premature infinity=%d, amd64 subnormal stage=%d",
			alwaysPositive, roundedGamma, largeInfinite, subnormalStage)
	}
	fmt.Printf("Lgamma: %d native and %d Lean rows; sign=%d, rounded gamma=%d, premature infinity=%d, amd64 subnormal stage=%d mutant failures\n",
		len(tests), len(rows), alwaysPositive, roundedGamma, largeInfinite, subnormalStage)
}
