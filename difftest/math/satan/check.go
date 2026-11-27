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
	profile, err := nativeProfile()
	if err != nil {
		return err
	}
	edges := []uint64{1 << 63, 0x7ff0000000000000}
	for exponent := uint64(0); exponent <= 2046; exponent++ {
		for _, mantissa := range []uint64{0, 1, 1 << 51, 1<<52 - 1} {
			edges = append(edges, exponent<<52|mantissa)
		}
	}
	for _, center := range []float64{0x1p-27, 0.66, 1, tan3PiOver8, 4, 0x1p52} {
		for _, x := range []float64{math.Nextafter(center, 0), center, math.Nextafter(center, math.Inf(1))} {
			edges = append(edges, binary64.Load(x))
		}
	}
	seeds := []int64{17, 2718, 65537}
	cases := append([]uint64(nil), edges...)
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			cases = append(cases, sample(rng, i))
		}
	}
	var input, fixture bytes.Buffer
	for _, word := range cases {
		fmt.Fprintln(&input, profile, word)
	}
	halfPi := precise(0).Set(trig.HalfPi())
	for i, model := range check.LeanRows("godiff-math-satan", len(cases), &input) {
		word := cases[i]
		fmt.Fprintln(&fixture, word, model, binary64.Load(idealAtan(binary64.Store(word), halfPi)))
	}
	dir, err := os.MkdirTemp("", "golean-math-satan-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_SATAN_CASES", path); err != nil {
		return err
	}
	const data = "difftest/math/satan/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanSatan", data+"bridge_test.go", data+"satan_test.go"); err != nil {
		return err
	}
	fmt.Printf("satan matched %d boundaries and %d samples with seeds %v.\n", len(edges), len(cases)-len(edges), seeds)
	return nil
}

func sample(rng *rand.Rand, i int) uint64 {
	switch i % 4 {
	case 0:
		return rng.Uint64() % 0x7ff0000000000000
	case 1:
		return binary64.Load(rng.Float64() * 4)
	case 2:
		center := []float64{0.66, 1, tan3PiOver8}[rng.Intn(3)]
		return binary64.Load(center) + uint64(rng.Intn(33)-16)
	default:
		return binary64.Load(math.Ldexp(rng.Float64(), rng.Intn(2098)-1074))
	}
}

func precise(value float64) *big.Float {
	return new(big.Float).SetPrec(256).SetFloat64(value)
}

func idealAtan(x float64, halfPi *big.Float) float64 {
	if x == 0 {
		return 0
	}
	if math.IsInf(x, 1) || x > 0x1p60 {
		return math.Pi / 2
	}
	// The omitted cubic correction is less than half an ulp here.
	if x < 0x1p-30 {
		return x
	}
	magnitude := precise(x)
	reciprocal := x > 1
	if reciprocal {
		magnitude.Quo(precise(1), magnitude)
	}
	// Halve the angle so that the alternating series converges faster.
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
	sum.SetMantExp(sum, 1)
	if reciprocal {
		sum.Sub(halfPi, sum)
	}
	result, _ := sum.Float64()
	return result
}

func nativeProfile() (int, error) {
	switch runtime.GOARCH {
	case "arm64":
		return 1, nil
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			return 0, fmt.Errorf("cannot determine satan build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1, nil
			}
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("unmodeled satan architecture %s", runtime.GOARCH)
	}
}
