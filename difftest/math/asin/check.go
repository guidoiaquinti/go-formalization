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
	var maxDistance uint64
	var maxUpperError float64
	var fixture bytes.Buffer
	check.CompareSamples("godiff-math-asin", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return observe(math.Asin(binary64.Store(word)))
	}, func(word uint64) string {
		x := binary64.Store(word)
		magnitude := math.Abs(x)
		switch {
		case math.IsNaN(x) || magnitude > 1:
			fmt.Fprintln(&fixture, word, "NaN", 0, "exact")
			return "NaN"
		case x == 0:
			expected := observe(x)
			fmt.Fprintln(&fixture, word, expected, word, "exact")
			return expected
		case magnitude == 1:
			want := math.Copysign(math.Pi/2, x)
			expected := observe(want)
			fmt.Fprintln(&fixture, word, expected, binary64.Load(want), "exact")
			return expected
		}
		got, want := math.Asin(x), idealAsin(x)
		bits, ref := binary64.Load(got), binary64.Load(want)
		if binary64.Load(math.Asin(-x)) != bits^(1<<63) {
			log.Fatalf("Asin sign reflection failed for %#016x", word)
		}
		kind := "numeric"
		if profile != 2 && magnitude > 0.7 {
			kind = "upper"
			// Rounding x*x before subtraction amplifies the radicand
			// error through sqrt near one. Check an absolute bound here.
			delta := math.Abs(got - want)
			if math.IsNaN(delta) || delta > 0x1p-38 {
				log.Fatalf("Asin(%#016x) = %.17g, reference = %.17g, absolute error %g", word, got, want, delta)
			}
			maxUpperError = max(maxUpperError, delta)
		} else {
			distance := max(bits, ref) - min(bits, ref)
			if distance > 4 {
				log.Fatalf("Asin(%#016x) = %.17g, reference = %.17g, encoding distance %d", word, got, want, distance)
			}
			maxDistance = max(maxDistance, distance)
		}
		expected := observe(got)
		// CompareSamples checks this public observation against the Lean word.
		fmt.Fprintln(&fixture, word, expected, ref, kind)
		return expected
	})
	fmt.Printf("Asin 256-bit reference: maximum encoding distance %d; maximum separately rounded upper-branch absolute error %g.\n",
		maxDistance, maxUpperError)
	dir, err := os.MkdirTemp("", "golean-math-asin-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_ASIN_CASES", path); err != nil {
		return err
	}
	const data = "difftest/math/asin/testdata/"
	return check.TestInPackage("math", "TestGoLeanAsin", data+"bridge_test.go", data+"asin_test.go")
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

func idealAsin(x float64) float64 {
	if math.Abs(x) < 0x1p-30 {
		return x
	}
	precise := func(value float64) *big.Float {
		return new(big.Float).SetPrec(256).SetFloat64(value)
	}
	magnitude := precise(math.Abs(x))
	square := precise(0).Mul(magnitude, magnitude)
	root := precise(0).Sqrt(precise(0).Sub(precise(1), square))
	// asin(x) = 2 atan(x/(1+sqrt(1-x*x))). Halve the atan
	// angle once more to keep the series argument at most sqrt(2)-1.
	argument := magnitude.Quo(magnitude, root.Add(root, precise(1)))
	root.Mul(argument, argument)
	root.Sqrt(root.Add(root, precise(1)))
	argument.Quo(argument, root.Add(root, precise(1)))
	factor := precise(0).Mul(argument, argument)
	factor.Neg(factor)
	power, sum, term, divisor := precise(0).Set(argument), precise(0), precise(0), precise(0)
	for k := int64(1); k <= 255; k += 2 {
		sum.Add(sum, term.Quo(power, divisor.SetInt64(k)))
		power.Mul(power, factor)
	}
	sum.SetMantExp(sum, 2)
	result, _ := sum.Float64()
	return math.Copysign(result, x)
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
			log.Fatal("cannot determine Asin build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		log.Fatalf("unmodeled Asin architecture %s", runtime.GOARCH)
		return 0
	}
}
