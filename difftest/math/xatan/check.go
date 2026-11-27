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
)

const limit = 0.66

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
	var edges []uint64
	const sign = uint64(1) << 63
	for exponent := uint64(0); exponent <= 1022; exponent++ {
		for _, mantissa := range []uint64{0, 1, 1 << 51, 1<<52 - 1} {
			word := exponent<<52 | mantissa
			if word <= binary64.Load(limit) {
				edges = append(edges, word, word|sign)
			}
		}
	}
	for _, center := range []float64{0x1p-27, 0.125, 0.25, 0.5, limit} {
		for _, x := range []float64{math.Nextafter(center, 0), center, math.Nextafter(center, math.Inf(1))} {
			if x <= limit {
				word := binary64.Load(x)
				edges = append(edges, word, word|sign)
			}
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
	for i, model := range check.LeanRows("godiff-math-xatan", len(cases), &input) {
		word := cases[i]
		fmt.Fprintln(&fixture, word, model, binary64.Load(idealAtan(binary64.Store(word))))
	}
	dir, err := os.MkdirTemp("", "golean-math-xatan-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_XATAN_CASES", path); err != nil {
		return err
	}
	const data = "difftest/math/xatan/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanXatan", data+"bridge_test.go", data+"xatan_test.go"); err != nil {
		return err
	}
	fmt.Printf("xatan matched %d boundaries and %d samples with seeds %v.\n", len(edges), len(cases)-len(edges), seeds)
	return nil
}

func sample(rng *rand.Rand, i int) uint64 {
	switch i % 4 {
	case 0:
		return (rng.Uint64() % (binary64.Load(limit) + 1)) | (rng.Uint64() & (1 << 63))
	case 1:
		return binary64.Load((2*rng.Float64() - 1) * limit)
	case 2:
		return (binary64.Load(limit) - uint64(rng.Intn(129))) | (rng.Uint64() & (1 << 63))
	default:
		return binary64.Load(math.Ldexp(rng.Float64()-0.5, -rng.Intn(1074)))
	}
}

func idealAtan(x float64) float64 {
	if x == 0 {
		return 0
	}
	// The omitted cubic correction is less than half an ulp here.
	if math.Abs(x) < 0x1p-30 {
		return x
	}
	precise := func(value float64) *big.Float {
		return new(big.Float).SetPrec(256).SetFloat64(value)
	}
	// Halve the angle so that the alternating series converges faster.
	argument := precise(math.Abs(x))
	root := precise(0).Mul(argument, argument)
	root.Sqrt(root.Add(root, precise(1)))
	argument.Quo(argument, root.Add(root, precise(1)))
	factor := precise(0).Mul(argument, argument)
	factor.Neg(factor)
	power, sum, term, divisor := precise(0).Set(argument), precise(0), precise(0), precise(0)
	for k := int64(1); k <= 255; k += 2 {
		sum.Add(sum, term.Quo(power, divisor.SetInt64(k)))
		power.Mul(power, factor)
	}
	sum.SetMantExp(sum, 1)
	result, _ := sum.Float64()
	return math.Copysign(result, x)
}

func nativeProfile() (int, error) {
	switch runtime.GOARCH {
	case "arm64":
		return 1, nil
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			return 0, fmt.Errorf("cannot determine xatan build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1, nil
			}
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("unmodeled xatan architecture %s", runtime.GOARCH)
	}
}
