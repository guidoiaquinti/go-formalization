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

const tan3PiOver8 = 2.41421356237309504880

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	profile := nativeProfile()
	edges := binary64.Boundaries()
	for _, center := range []float64{0x1p-27, 0.25, 0.66, 1, tan3PiOver8, 4, 0x1p26, 0x1p52, 0x1p1023} {
		for _, value := range []float64{math.Nextafter(center, 0), center, math.Nextafter(center, math.Inf(1))} {
			edges = append(edges, binary64.Load(value), binary64.Load(-value))
		}
	}
	halfPi := precise(0).Set(trig.HalfPi())
	var maxDistance uint64
	var fixture bytes.Buffer
	check.CompareSamples("godiff-math-atan", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return observe(math.Atan(binary64.Store(word)))
	}, func(word uint64) string {
		x := binary64.Store(word)
		var special string
		switch {
		case math.IsNaN(x):
			special = "NaN"
		case x == 0:
			special = observe(x)
		case math.IsInf(x, 0):
			special = observe(math.Copysign(math.Pi/2, x))
		}
		if special != "" {
			fmt.Fprintln(&fixture, word, special, 0, true)
			return special
		}
		got, want := math.Atan(x), idealAtan(x, halfPi)
		bits, ref := binary64.Load(got), binary64.Load(want)
		if binary64.Load(math.Atan(-x)) != bits^(1<<63) {
			log.Fatalf("Atan sign reflection failed for %#016x", word)
		}
		distance := max(bits, ref) - min(bits, ref)
		if distance > 4 {
			log.Fatalf("Atan(%#016x) = %.17g, reference = %.17g, encoding distance %d", word, got, want, distance)
		}
		maxDistance = max(maxDistance, distance)
		expected := observe(got)
		// CompareSamples checks that the public observation equals the Lean output.
		fmt.Fprintln(&fixture, word, expected, ref, false)
		return expected
	})
	fmt.Printf("Atan 256-bit series reference: maximum distance %d encoded steps.\n", maxDistance)
	dir, err := os.MkdirTemp("", "golean-math-atan-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_ATAN_CASES", path); err != nil {
		return err
	}
	const data = "difftest/math/atan/testdata/"
	return check.TestInPackage("math", "TestGoLeanAtan", data+"bridge_test.go", data+"atan_test.go")
}

func sample(rng *rand.Rand, i int) uint64 {
	switch i % 5 {
	case 0:
		return rng.Uint64()
	case 1:
		return binary64.Load((2*rng.Float64() - 1) * 0.66)
	case 2:
		return binary64.Load((2*rng.Float64() - 1) * 4)
	case 3:
		center := []float64{0.66, 1, tan3PiOver8}[rng.Intn(3)]
		return (binary64.Load(center) + uint64(rng.Intn(33)-16)) | (rng.Uint64() & (1 << 63))
	default:
		return binary64.Load(math.Ldexp(rng.Float64()-0.5, rng.Intn(2098)-1074))
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

// Use reciprocal reflection above one, then the half-angle identity to
// put the argument of the alternating series in [0, sqrt(2)-1].
func idealAtan(x float64, halfPi *big.Float) float64 {
	magnitude := precise(math.Abs(x))
	reciprocal := math.Abs(x) > 1
	if reciprocal {
		magnitude.Quo(precise(1), magnitude)
	}
	square := precise(0).Mul(magnitude, magnitude)
	root := square.Sqrt(square.Add(square, precise(1)))
	argument := magnitude.Quo(magnitude, root.Add(root, precise(1)))
	factor := precise(0).Mul(argument, argument)
	factor.Neg(factor)
	power, sum, term, divisor := precise(0).Set(argument), precise(0), precise(0), precise(0)
	for k := int64(1); k <= 255; k += 2 {
		sum.Add(sum, term.Quo(power, divisor.SetInt64(k)))
		power.Mul(power, factor)
	}
	sum.SetMantExp(sum, 1)
	if reciprocal {
		sum.Sub(halfPi, sum)
	}
	result, _ := sum.Float64()
	return math.Copysign(result, x)
}

func nativeProfile() int {
	switch runtime.GOARCH {
	case "arm64":
		return 1
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			log.Fatal("cannot determine Atan build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		log.Fatalf("unmodeled Atan architecture %s", runtime.GOARCH)
		return 0
	}
}
