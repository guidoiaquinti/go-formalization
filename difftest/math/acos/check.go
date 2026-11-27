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
	profile := nativeProfile()
	edges := binary64.Boundaries()
	for _, center := range []float64{0x1p-27, 0.25, 0.5, 0.66, 0.7, 0.75, 0.875, 1, 2} {
		for _, value := range []float64{math.Nextafter(center, 0), center, math.Nextafter(center, math.Inf(1))} {
			edges = append(edges, binary64.Load(value), binary64.Load(-value))
		}
	}
	pi := precise(0).Set(trig.HalfPi())
	pi.SetMantExp(pi, 1)
	var maxPositiveError, maxNegativeError float64
	var fixture bytes.Buffer
	check.CompareSamples("godiff-math-acos", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return observe(math.Acos(binary64.Store(word)))
	}, func(word uint64) string {
		x := binary64.Store(word)
		magnitude := math.Abs(x)
		var special float64
		exact := true
		switch {
		case math.IsNaN(x) || magnitude > 1:
			special = math.NaN()
		case x == 0:
			special = math.Pi / 2
		case x == 1:
			special = 0
		case x == -1:
			special = math.Pi
		default:
			exact = false
		}
		if exact {
			expected := observe(special)
			fmt.Fprintln(&fixture, word, expected, binary64.Load(special), "exact")
			return expected
		}
		got, want := math.Acos(x), idealAcos(x, pi)
		if staged := math.Pi/2 - math.Asin(x); binary64.Load(staged) != binary64.Load(got) {
			log.Fatalf("Acos(%#016x): Go=%#016x, rounded subtraction=%#016x",
				word, binary64.Load(got), binary64.Load(staged))
		}
		delta := math.Abs(got - want)
		bound := float64(0x1p-49)
		kind := "numeric"
		if profile != 2 && magnitude > 0.7 {
			kind = "upper"
			// The separately rounded radicand amplifies error near one.
			bound = 0x1p-38
		}
		if math.IsNaN(delta) || delta > bound {
			log.Fatalf("Acos(%#016x) = %.17g, reference = %.17g, absolute difference %g", word, got, want, delta)
		}
		if x < 0 {
			maxNegativeError = max(maxNegativeError, delta)
		} else {
			maxPositiveError = max(maxPositiveError, delta)
		}
		expected := observe(got)
		// CompareSamples checks this public observation against the Lean word.
		fmt.Fprintln(&fixture, word, expected, binary64.Load(want), kind)
		return expected
	})
	fmt.Printf("Acos 256-bit reference: maximum absolute difference positive inputs=%g negative inputs=%g.\n",
		maxPositiveError, maxNegativeError)
	dir, err := os.MkdirTemp("", "golean-math-acos-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_ACOS_CASES", path); err != nil {
		return err
	}
	const data = "difftest/math/acos/testdata/"
	return check.TestInPackage("math", "TestGoLeanAcos", data+"bridge_test.go", data+"acos_test.go")
}

func sample(rng *rand.Rand, i int) uint64 {
	switch i % 6 {
	case 0:
		return rng.Uint64()
	case 1:
		return binary64.Load(2*rng.Float64() - 1)
	case 2:
		center := []float64{0.7, 1}[rng.Intn(2)]
		return (binary64.Load(center) + uint64(rng.Intn(129)-64)) | (rng.Uint64() & (1 << 63))
	case 3:
		return (binary64.Load(1) - uint64(rng.Intn(1<<24))) | (rng.Uint64() & (1 << 63))
	case 4:
		return binary64.Load(math.Copysign(1-math.Ldexp(rng.Float64(), -rng.Intn(55)), float64(rng.Int63n(2)*2-1)))
	default:
		return binary64.Load(math.Ldexp(rng.Float64()-0.5, -rng.Intn(1074)))
	}
}

func precise(value float64) *big.Float {
	return new(big.Float).SetPrec(256).SetFloat64(value)
}

func idealAcos(x float64, pi *big.Float) float64 {
	magnitude := precise(math.Abs(x))
	// acos(|x|) = 2 atan(sqrt((1-|x|)/(1+|x|))).
	ratio := precise(0).Sub(precise(1), magnitude)
	ratio.Quo(ratio, precise(0).Add(precise(1), magnitude))
	ratio.Sqrt(ratio)
	// Halve the atan angle to put the series argument in [0, sqrt(2)-1].
	root := precise(0).Mul(ratio, ratio)
	root.Sqrt(root.Add(root, precise(1)))
	argument := ratio.Quo(ratio, root.Add(root, precise(1)))
	factor := precise(0).Mul(argument, argument)
	factor.Neg(factor)
	power, sum, term, divisor := precise(0).Set(argument), precise(0), precise(0), precise(0)
	for k := int64(1); k <= 255; k += 2 {
		sum.Add(sum, term.Quo(power, divisor.SetInt64(k)))
		power.Mul(power, factor)
	}
	sum.SetMantExp(sum, 2)
	if x < 0 {
		sum.Sub(pi, sum)
	}
	result, _ := sum.Float64()
	return result
}

func observe(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(binary64.Load(x), 10)
}

func nativeProfile() int {
	switch runtime.GOARCH {
	case "arm64":
		return 2
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			log.Fatal("cannot determine Acos build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		log.Fatalf("unmodeled Acos architecture %s", runtime.GOARCH)
		return 0
	}
}
