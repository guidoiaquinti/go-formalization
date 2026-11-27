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

// The scaled integral e^(x²) integral_0^x e^(-t²)dt solves F'=1+2xF.
// Its positive terms avoid cancellation while locating the inverse.
func forward(x float64, scale *big.Float) *big.Float {
	p := precise(x)
	square := precise(0).Mul(p, p)
	twice := precise(0).Mul(precise(2), square)
	integral, term := precise(x), precise(x)
	exponential, expTerm, divisor := precise(1), precise(1), precise(0)
	for n := int64(1); n <= 220; n++ {
		term.Mul(term, twice)
		term.Quo(term, divisor.SetInt64(2*n+1))
		integral.Add(integral, term)
		expTerm.Mul(expTerm, square)
		expTerm.Quo(expTerm, divisor.SetInt64(n))
		exponential.Add(exponential, expTerm)
	}
	integral.Quo(integral, exponential)
	return integral.Mul(integral, scale)
}

type reference struct {
	x      float64
	lo, hi uint64
}

func newReference(x float64, scale *big.Float) reference {
	r := reference{x: x}
	if x == 0 || math.Abs(x) >= 1 || math.IsNaN(x) {
		return r
	}
	target := precise(math.Abs(x))
	r.hi = math.Float64bits(6)
	if forward(6, scale).Cmp(target) <= 0 {
		log.Fatalf("erf inverse exceeds the initial magnitude bracket for %x", x)
	}
	for r.hi-r.lo > 1 {
		mid := r.lo + (r.hi-r.lo)/2
		if forward(math.Float64frombits(mid), scale).Cmp(target) < 0 {
			r.lo = mid
		} else {
			r.hi = mid
		}
	}
	return r
}

func (r reference) agrees(y float64) bool {
	x := r.x
	if math.IsNaN(x) || math.Abs(x) > 1 {
		return math.IsNaN(y)
	}
	if x == 0 {
		return math.Float64bits(x) == math.Float64bits(y)
	}
	if math.Abs(x) == 1 {
		return math.IsInf(y, int(math.Copysign(1, x)))
	}
	if math.IsNaN(y) || math.IsInf(y, 0) || math.Signbit(y) != math.Signbit(x) {
		return false
	}
	const radius = 10
	word := math.Float64bits(math.Abs(y))
	return word+radius >= r.lo && word <= r.hi+radius
}

func main() {
	scale, _, err := big.ParseFloat(
		"1.128379167095512573896158903121545171688101258657997713688171443421284936882986828973487320404214726886056695812725",
		10, precision, big.ToNearestEven)
	if err != nil {
		log.Fatal(err)
	}
	values := []float64{0, math.Copysign(0, -1), math.Inf(1), math.Inf(-1), math.NaN(),
		math.Float64frombits(0xfff0000000000001)}
	for _, center := range []float64{math.SmallestNonzeroFloat64, 0x1p-1022, 0x1p-28,
		0.5, 0.85, 0.9, 1 - 2*math.Exp(-25), 1} {
		for _, x := range []float64{math.Nextafter(center, 0), center, math.Nextafter(center, math.Inf(1))} {
			values = append(values, x, -x)
		}
	}
	const profiles = 3
	var requests bytes.Buffer
	tests := make([]reference, len(values))
	var lostSign, linear, earlyEndpoint int
	for i, x := range values {
		test := newReference(x, scale)
		tests[i] = test
		got := math.Erfinv(x)
		if !test.agrees(got) {
			log.Fatalf("Go Erfinv(%x)=%x outside inverse bracket [%#x,%#x]", x, got, test.lo, test.hi)
		}
		if !test.agrees(math.Abs(got)) {
			lostSign++
		}
		if 0 < math.Abs(x) && math.Abs(x) <= 0.85 && !test.agrees((math.SqrtPi/2)*x) {
			linear++
		}
		if 0.99 < math.Abs(x) && math.Abs(x) < 1 && !test.agrees(math.Copysign(math.Inf(1), x)) {
			earlyEndpoint++
		}
		for profile := range profiles {
			fmt.Fprintln(&requests, profile, math.Float64bits(x))
		}
	}
	rows := check.LeanRows("godiff-math-erfinv", profiles*len(tests), &requests)
	for i, row := range rows {
		test, profile := tests[i/profiles], i%profiles
		got := math.NaN()
		if row != "NaN" {
			word, err := strconv.ParseUint(row, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			got = math.Float64frombits(word)
		}
		if !test.agrees(got) {
			log.Fatalf("Lean profile %d Erfinv(%x)=%x outside [%#x,%#x]", profile, test.x, got, test.lo, test.hi)
		}
		native := math.Erfinv(test.x)
		if runtime.GOARCH == "arm64" && profile == 2 && !math.IsNaN(native) &&
			math.Float64bits(got) != math.Float64bits(native) {
			log.Fatalf("arm64 Erfinv(%x)=%x, Lean=%x", test.x, native, got)
		}
	}
	if lostSign != 24 || linear != 10 || earlyEndpoint != 8 {
		log.Fatalf("unexpected controls: sign=%d, linear central=%d, early endpoint=%d", lostSign, linear, earlyEndpoint)
	}
	fmt.Printf("Erfinv: %d native and %d Lean profile rows; sign=%d, linear central=%d, early endpoint=%d mutant failures\n",
		len(tests), len(rows), lostSign, linear, earlyEndpoint)
}
