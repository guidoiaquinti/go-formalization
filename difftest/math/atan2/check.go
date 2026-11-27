package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/trig"
)

type input struct{ y, x uint64 }

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	profile := nativeProfile()
	words := binary64.Boundaries()
	for _, x := range []float64{0.5, 0.66, 1, 2, 2.41421356237309504880, 4} {
		words = append(words, binary64.Load(x), binary64.Load(-x))
	}
	var edges []input
	for _, y := range words {
		for _, x := range words {
			edges = append(edges, input{y, x})
		}
	}
	halfPi := precise(0).Set(trig.HalfPi())
	pi := precise(0).SetMantExp(halfPi, 1)
	for _, x := range []float64{math.Nextafter(-2, 0), -2, math.Nextafter(-2, math.Inf(-1))} {
		for _, y := range []float64{math.SmallestNonzeroFloat64, -math.SmallestNonzeroFloat64} {
			edges = append(edges, input{binary64.Load(y), binary64.Load(x)})
		}
	}
	checkUnderflow(halfPi, pi)

	var maxDistance uint64
	var lowerLeftUnderflows int
	var fixture bytes.Buffer
	check.CompareSamples("godiff-math-atan2", edges, sample, func(c input) string {
		return fmt.Sprintf("%d %d %d", profile, c.y, c.x)
	}, func(c input) string {
		return observe(math.Atan2(binary64.Store(c.y), binary64.Store(c.x)))
	}, func(c input) string {
		y, x := binary64.Store(c.y), binary64.Store(c.x)
		if want, exact := specialResult(y, x); exact {
			expected := observe(want)
			fmt.Fprintln(&fixture, c.y, c.x, expected, binary64.Load(want), "exact")
			return expected
		}
		got := math.Atan2(y, x)
		want, lowerLeftZero := idealAtan2(y, x, halfPi, pi)
		if lowerLeftZero {
			if y/x != 0 || want != -math.Pi || got != math.Pi {
				log.Fatalf("Atan2(%#016x, %#016x) unexpected lower-left zero case: ratio=%g Go=%g reference=%g",
					c.y, c.x, y/x, got, want)
			}
			lowerLeftUnderflows++
			expected := observe(math.Pi)
			fmt.Fprintln(&fixture, c.y, c.x, expected, binary64.Load(want), "underflow")
			return expected
		}
		bits, ref := binary64.Load(got), binary64.Load(want)
		distance := max(bits, ref) - min(bits, ref)
		if distance > 4 {
			log.Fatalf("Atan2(%#016x, %#016x) = %.17g, reference = %.17g, encoding distance %d",
				c.y, c.x, got, want, distance)
		}
		maxDistance = max(maxDistance, distance)
		expected := observe(got)
		// CompareSamples checks this public observation against the Lean word.
		fmt.Fprintln(&fixture, c.y, c.x, expected, ref, "numeric")
		return expected
	})
	fmt.Printf("Atan2 256-bit reference: maximum distance %d encoded steps outside %d lower-left zero-rounding cases.\n",
		maxDistance, lowerLeftUnderflows)
	dir, err := os.MkdirTemp("", "golean-math-atan2-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_ATAN2_CASES", path); err != nil {
		return err
	}
	const data = "difftest/math/atan2/testdata/"
	return check.TestInPackage("math", "TestGoLeanAtan2", data+"bridge_test.go", data+"atan2_test.go")
}

func sample(rng *rand.Rand, i int) input {
	switch i % 5 {
	case 0:
		return input{rng.Uint64(), rng.Uint64()}
	case 1:
		return input{binary64.Load(8*rng.Float64() - 4), binary64.Load(8*rng.Float64() - 4)}
	case 2:
		mantissa := rng.Uint64()&((1<<52)-1) | 1
		y := mantissa | (rng.Uint64() & (1 << 63))
		x := binary64.Load(math.Ldexp(0.5+rng.Float64()/2, rng.Intn(1023)))
		return input{y, x | (rng.Uint64() & (1 << 63))}
	case 3:
		return input{binary64.Load(2*rng.Float64() - 1), binary64.Load(math.Ldexp(rng.Float64()-0.5, -rng.Intn(1074)))}
	default:
		return input{rng.Uint64(), binary64.Load(math.Ldexp(rng.Float64()-0.5, rng.Intn(2098)-1074))}
	}
}

func specialResult(y, x float64) (float64, bool) {
	switch {
	case math.IsNaN(y) || math.IsNaN(x):
		return math.NaN(), true
	case y == 0:
		if math.Signbit(x) {
			return math.Copysign(math.Pi, y), true
		}
		return y, true
	case x == 0:
		return math.Copysign(math.Pi/2, y), true
	case math.IsInf(x, 1):
		if math.IsInf(y, 0) {
			return math.Copysign(math.Pi/4, y), true
		}
		return math.Copysign(0, y), true
	case math.IsInf(x, -1):
		if math.IsInf(y, 0) {
			return math.Copysign(3*math.Pi/4, y), true
		}
		return math.Copysign(math.Pi, y), true
	case math.IsInf(y, 0):
		return math.Copysign(math.Pi/2, y), true
	default:
		return 0, false
	}
}

func checkUnderflow(halfPi, pi *big.Float) {
	y := -math.SmallestNonzeroFloat64
	cases := []struct {
		x                float64
		quotient, output uint64
	}{
		{math.Nextafter(-2, 0), 1, 0xc00921fb54442d18},
		{-2, 0, 0x400921fb54442d18},
		{math.Nextafter(-2, math.Inf(-1)), 0, 0x400921fb54442d18},
	}
	for _, c := range cases {
		q, got := binary64.Load(y/c.x), binary64.Load(math.Atan2(y, c.x))
		want, underflows := idealAtan2(y, c.x, halfPi, pi)
		if q != c.quotient || got != c.output || want != -math.Pi || underflows != (c.quotient == 0) {
			log.Fatalf("Atan2 lower-left boundary failed: y=%#016x x=%#016x quotient=%#016x Go=%#016x reference=%#016x",
				binary64.Load(y), binary64.Load(c.x), q, got, binary64.Load(want))
		}
		fmt.Printf("Atan2 lower-left boundary: y=%#016x x=%#016x quotient=%#016x Go=%#016x reference=%#016x.\n",
			binary64.Load(y), binary64.Load(c.x), q, got, binary64.Load(want))
	}
}

func observe(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(binary64.Load(x), 10)
}

func precise(value float64) *big.Float {
	return new(big.Float).SetPrec(256).SetFloat64(value)
}

func idealAtan2(y, x float64, halfPi, pi *big.Float) (float64, bool) {
	numerator, denominator := precise(math.Abs(y)), precise(math.Abs(x))
	// Rounding at or below half the smallest subnormal gives positive
	// zero for the quotient of two negative inputs.
	zeroLimit := precise(0).SetMantExp(denominator, -1075)
	lowerLeftZero := x < 0 && y < 0 && numerator.Cmp(zeroLimit) <= 0
	reciprocal := numerator.Cmp(denominator) > 0
	if reciprocal {
		numerator, denominator = denominator, numerator
	}
	magnitude := numerator.Quo(numerator, denominator)
	angle := magnitude
	if magnitude.Cmp(precise(0x1p-30)) >= 0 {
		root := precise(0).Mul(magnitude, magnitude)
		root.Sqrt(root.Add(root, precise(1)))
		argument := magnitude.Quo(magnitude, root.Add(root, precise(1)))
		factor := precise(0).Mul(argument, argument)
		factor.Neg(factor)
		power, sum, term, divisor := precise(0).Set(argument), precise(0), precise(0), precise(0)
		for k := int64(1); k <= 255; k += 2 {
			sum.Add(sum, term.Quo(power, divisor.SetInt64(k)))
			power.Mul(power, factor)
		}
		angle = sum.SetMantExp(sum, 1)
	}
	if reciprocal {
		angle.Sub(halfPi, angle)
	}
	if x < 0 {
		angle.Sub(pi, angle)
	}
	result, _ := angle.Float64()
	return math.Copysign(result, y), lowerLeftZero
}

func nativeProfile() int {
	switch runtime.GOARCH {
	case "arm64":
		return 1
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			log.Fatal("cannot determine Atan2 build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		log.Fatalf("unmodeled Atan2 architecture %s", runtime.GOARCH)
		return 0
	}
}
