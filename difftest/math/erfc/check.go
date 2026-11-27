package main

import (
	"bufio"
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
	dir, err := os.MkdirTemp("", "golean-math-erfc-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	fixture := bufio.NewWriter(file)
	if err := os.Setenv("GOLEAN_MATH_ERFC_CASES", path); err != nil {
		return err
	}
	exact := func(word uint64, expected string) string {
		fmt.Fprintln(fixture, word, expected, expected, "exact")
		return expected
	}

	profile, err := nativeProfile()
	if err != nil {
		return err
	}
	edges := binary64.Boundaries()
	for _, center := range []float64{0x1p-56, 0.25, 0.84375, 1, 1.25, 20.0 / 7, 4,
		6, 26.543, 27.2, 27.3, 28, math.MaxFloat64} {
		for _, value := range []float64{math.Nextafter(center, 0), center, math.Nextafter(center, math.Inf(1))} {
			edges = append(edges, binary64.Load(value), binary64.Load(-value))
		}
	}
	pi, _, err := big.ParseFloat(
		"3.141592653589793238462643383279502884197169399375105820974944592307816406286208998628034825342117068",
		10, 256, big.ToNearestEven)
	if err != nil {
		return err
	}
	scale := precise(2).Quo(precise(2), pi.Sqrt(pi))
	var maximum uint64
	check.CompareSamples("godiff-math-erfc", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return observe(math.Erfc(binary64.Store(word)))
	}, func(word uint64) string {
		x := binary64.Store(word)
		switch {
		case math.IsNaN(x):
			return exact(word, "NaN")
		case x == 0:
			return exact(word, observe(1))
		case x >= 28:
			return exact(word, observe(0))
		case x <= -28:
			return exact(word, observe(2))
		}
		got, want := math.Erfc(x), idealErfc(x, scale)
		a, b := binary64.Load(got), binary64.Load(want)
		distance := max(a, b) - min(a, b)
		if distance > 6 {
			log.Fatalf("Erfc(%#016x) = %#016x, reference = %#016x, distance = %d", word, a, b, distance)
		}
		maximum = max(maximum, distance)
		expected := observe(got)
		fmt.Fprintln(fixture, word, expected, observe(want), "numeric")
		return expected
	})
	fmt.Printf("Erfc 256-bit reference: maximum distance %d binary64 steps.\n", maximum)
	if err := fixture.Flush(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	const data = "difftest/math/erfc/testdata/"
	return check.TestInPackage("math", "TestGoLeanErfc", data+"bridge_test.go", data+"erfc_test.go")
}

func sample(rng *rand.Rand, i int) uint64 {
	switch i % 6 {
	case 0:
		return rng.Uint64()
	case 1:
		return binary64.Load(4*rng.Float64() - 2)
	case 2:
		return binary64.Load(32*rng.Float64() - 4)
	case 3:
		centers := []float64{0x1p-56, 0.25, 0.84375, 1.25, 20.0 / 7, 6, 28}
		center := binary64.Load(centers[rng.Intn(len(centers))])
		return (center + uint64(rng.Intn(129)-64)) | (rng.Uint64() & (1 << 63))
	case 4:
		return binary64.Load(math.Ldexp(rng.Float64()-0.5, -rng.Intn(1074)))
	default:
		return binary64.Load(26 + 2*rng.Float64())
	}
}

func precise(value float64) *big.Float {
	return new(big.Float).SetPrec(256).SetFloat64(value)
}

func integratedGaussian(x float64, scale *big.Float) *big.Float {
	power := precise(x)
	ratio := precise(0).Mul(power, power)
	ratio.Neg(ratio)
	sum, coefficient, term, divisor := precise(0).Set(power), precise(0), precise(0), precise(0)
	for n := int64(1); n <= 160; n++ {
		power.Mul(power, ratio)
		power.Quo(power, coefficient.SetInt64(n))
		term.Quo(power, divisor.SetInt64(2*n+1))
		sum.Add(sum, term)
		if n > 16 && term.MantExp(nil) < -220 {
			break
		}
	}
	return sum.Mul(sum, scale)
}

func gammaTail(x float64, scale *big.Float) *big.Float {
	argument := precise(0).Mul(precise(x), precise(x))
	// Q(1/2,x²) uses the upper incomplete gamma continued fraction.
	fraction, b, numerator, coefficient := precise(0), precise(0), precise(0), precise(0)
	for n := int64(160); n >= 1; n-- {
		b.Add(argument, coefficient.SetFloat64(float64(2*n)+0.5))
		b.Add(b, fraction)
		numerator.SetFloat64(-float64(n) * (float64(n) - 0.5))
		fraction.Quo(numerator, b)
	}
	fraction.Add(fraction, b.Add(argument, coefficient.SetFloat64(0.5)))
	// Compute exp(-x²) at high precision through range scaling.
	reduced := precise(0).SetMantExp(argument, -10)
	reduced.Neg(reduced)
	sum, term := precise(1), precise(1)
	for n := int64(1); n <= 85; n++ {
		term.Mul(term, reduced)
		term.Quo(term, coefficient.SetInt64(n))
		sum.Add(sum, term)
	}
	for range 10 {
		sum.Mul(sum, sum)
	}
	sum.Mul(sum, precise(x))
	sum.Mul(sum, scale)
	sum.SetMantExp(sum, -1)
	return sum.Quo(sum, fraction)
}

func idealErfc(x float64, scale *big.Float) float64 {
	magnitude := math.Abs(x)
	if magnitude >= 28 {
		if x < 0 {
			return 2
		}
		return 0
	}
	var result *big.Float
	if magnitude < 4 {
		value := integratedGaussian(magnitude, scale)
		if x < 0 {
			result = value.Add(precise(1), value)
		} else {
			result = value.Sub(precise(1), value)
		}
	} else {
		result = gammaTail(magnitude, scale)
		if x < 0 {
			result.Sub(precise(2), result)
		}
	}
	out, _ := result.Float64()
	return out
}

func observe(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(binary64.Load(x), 10)
}

func nativeProfile() (int, error) {
	if runtime.GOARCH == "arm64" {
		return 4, nil
	}
	if runtime.GOARCH != "amd64" {
		return 0, fmt.Errorf("unmodeled Erfc architecture %s", runtime.GOARCH)
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return 0, fmt.Errorf("cannot determine Erfc build profile")
	}
	profile := 0
	for _, setting := range info.Settings {
		if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
			profile = 2
		}
	}
	dir, err := os.MkdirTemp("", "golean-erf-exp-profile-")
	if err != nil {
		return 0, err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "profile")
	if err := os.Setenv("GOLEAN_MATH_EXP_PROFILE", path); err != nil {
		return 0, err
	}
	const fixtures = "difftest/math/exp/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanExpProfile",
		fixtures+"profile_amd64_test.go", fixtures+"observe_profile_amd64_test.go"); err != nil {
		return 0, err
	}
	value, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	fused, err := strconv.ParseBool(string(value))
	if fused {
		profile++
	}
	return profile, err
}
