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

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	profile := 0
	switch runtime.GOARCH {
	case "arm64":
		profile = 2
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			return fmt.Errorf("cannot determine Tan build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				profile = 1
			}
		}
	default:
		return fmt.Errorf("unmodeled Tan architecture %s", runtime.GOARCH)
	}

	edges := binary64.Boundaries()
	for _, x := range []float64{0x1p-27, 1e-7, math.Pi / 4, math.Pi / 2, math.Pi, 3 * math.Pi / 2, 2 * math.Pi,
		0x1p29, 0x1p30, 0x1p64, 0x1p128, 0x1p512, 0x1p1023} {
		for _, value := range []float64{math.Nextafter(x, 0), x, math.Nextafter(x, math.Inf(1))} {
			edges = append(edges, binary64.Load(value), binary64.Load(-value))
		}
	}
	halfPi := trig.HalfPi()
	var maxError float64
	var fixture bytes.Buffer
	check.CompareSamples("godiff-math-tan", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return observe(math.Tan(binary64.Store(word)))
	}, func(word uint64) string {
		x := binary64.Store(word)
		switch {
		case math.IsNaN(x) || math.IsInf(x, 0):
			fmt.Fprintln(&fixture, word, "NaN", 0, true)
			return "NaN"
		case x == 0:
			expected := observe(x)
			fmt.Fprintln(&fixture, word, expected, 0, true)
			return expected
		}
		got, want := math.Tan(x), idealTan(x, halfPi)
		if opposite := binary64.Load(math.Tan(-x)); opposite != binary64.Load(got)^(1<<63) {
			log.Fatalf("Tan sign reflection failed for %#016x", word)
		}
		delta := math.Abs(got-want) / max(1, math.Abs(want))
		if math.IsNaN(delta) || delta > 0x1p-35 || math.Abs(want) <= 2 && delta > 0x1p-48 {
			log.Fatalf("Tan(%#016x) = %.17g, reference = %.17g, scaled absolute error = %g", word, got, want, delta)
		}
		maxError = max(maxError, delta)
		expected := observe(got)
		// CompareSamples checks that the public result equals the Lean output.
		fmt.Fprintln(&fixture, word, expected, binary64.Load(want), false)
		return expected
	})
	fmt.Printf("Tan 1,408-bit reduction reference: maximum scaled absolute error %g.\n", maxError)
	dir, err := os.MkdirTemp("", "golean-math-tan-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_TAN_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/tan/testdata/"
	return check.TestInPackage("math", "TestGoLeanTan", fixtures+"bridge_test.go", fixtures+"tan_test.go")
}

func sample(r *rand.Rand, i int) uint64 {
	switch i % 5 {
	case 0:
		return r.Uint64()
	case 1:
		return binary64.Load((2*r.Float64() - 1) * 0x1p29)
	case 2:
		return binary64.Load(math.Ldexp(r.Float64()+0.5, 29+r.Intn(995))) | r.Uint64()&(1<<63)
	case 3:
		multiple := (r.Float64() - 0.5) * 0x1p29
		return binary64.Load(math.Round(multiple)*math.Pi/4) + uint64(r.Intn(9)-4)
	default:
		return binary64.Load(math.Ldexp(r.Float64()-0.5, -r.Intn(1074)))
	}
}

func observe(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(binary64.Load(x), 10)
}

// idealTan evaluates a ratio of reduced Taylor series before rounding.
func idealTan(x float64, halfPi *big.Float) float64 {
	angle, q := reduceQuadrant(x, halfPi)
	sine, cosine := series(angle, false), series(angle, true)
	var result *big.Float
	if q&1 == 0 {
		result = sine.Quo(sine, cosine)
	} else {
		result = cosine.Neg(cosine).Quo(cosine, sine)
	}
	if math.Signbit(x) {
		result.Neg(result)
	}
	rounded, _ := result.Float64()
	return rounded
}

func reduceQuadrant(x float64, halfPi *big.Float) (*big.Float, uint64) {
	value := func(x float64) *big.Float { return new(big.Float).SetPrec(halfPi.Prec()).SetFloat64(x) }
	absolute := value(math.Abs(x))
	quotient := value(0).Quo(absolute, halfPi)
	quadrant, _ := quotient.Add(quotient, value(0.5)).Int(nil)
	angle := value(0).Sub(absolute, value(0).Mul(value(0).SetInt(quadrant), halfPi))
	q := new(big.Int).And(quadrant, big.NewInt(3)).Uint64()

	return angle, q
}

func series(angle *big.Float, cosine bool) *big.Float {
	arg := new(big.Float).SetPrec(256).Set(angle)
	negativeSquare := new(big.Float).SetPrec(256).Mul(arg, arg)
	negativeSquare.Neg(negativeSquare)
	term := new(big.Float).SetPrec(256).Set(arg)
	power := int64(1)
	if cosine {
		term.SetInt64(1)
		power = 0
	}
	sum, divisor := new(big.Float).SetPrec(256).Set(term), new(big.Float).SetPrec(256)
	for n := power + 2; n <= power+50; n += 2 {
		term.Mul(term, negativeSquare).Quo(term, divisor.SetInt64(n*(n-1)))
		sum.Add(sum, term)
	}
	return sum
}
