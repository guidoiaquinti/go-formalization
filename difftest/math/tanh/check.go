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

const saturation = 0.5 * 8.8029691931113054295988e+01

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
	boundaries := binary64.Boundaries()
	for _, center := range []float64{0x1p-27, 0.125, 0.5, 0.625, 1, 16, 19, 20, 32, saturation} {
		for _, value := range []float64{math.Nextafter(center, 0), center, math.Nextafter(center, math.Inf(1))} {
			boundaries = append(boundaries, binary64.Load(value), binary64.Load(-value))
		}
	}
	var maxDistance uint64
	var fixture bytes.Buffer
	check.CompareSamples("godiff-math-tanh", boundaries, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return observe(math.Tanh(binary64.Store(word)))
	}, func(word uint64) string {
		x := binary64.Store(word)
		var special string
		switch {
		case math.IsNaN(x):
			special = "NaN"
		case x == 0:
			special = observe(x)
		case math.Abs(x) > saturation:
			special = observe(math.Copysign(1, x))
		}
		if special != "" {
			fmt.Fprintln(&fixture, word, special, 0, true)
			return special
		}

		got, want := math.Tanh(x), idealTanh(x)
		bits, ref := binary64.Load(got), binary64.Load(want)
		if binary64.Load(math.Tanh(-x)) != bits^(1<<63) {
			log.Fatalf("Tanh sign reflection failed for %#016x", word)
		}
		distance := max(bits, ref) - min(bits, ref)
		if distance > 4 {
			log.Fatalf("Tanh(%#016x) = %.17g, reference = %.17g, encoding distance %d", word, got, want, distance)
		}
		maxDistance = max(maxDistance, distance)
		expected := observe(got)
		// CompareSamples has verified that this native public result equals Lean.
		fmt.Fprintln(&fixture, word, expected, ref, false)
		return expected
	})
	fmt.Printf("Tanh 256-bit exponential reference: maximum distance %d encoded steps.\n", maxDistance)
	dir, err := os.MkdirTemp("", "golean-math-tanh-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_TANH_CASES", path); err != nil {
		return err
	}
	const data = "difftest/math/tanh/testdata/"
	return check.TestInPackage("math", "TestGoLeanTanh", data+"bridge_test.go", data+"tanh_test.go")
}

func sample(rng *rand.Rand, i int) uint64 {
	switch i % 5 {
	case 0:
		return rng.Uint64()
	case 1:
		return binary64.Load((2*rng.Float64() - 1) * 0.625)
	case 2:
		return binary64.Load((2*rng.Float64() - 1) * saturation)
	case 3:
		center := []float64{0.625, 19, saturation}[rng.Intn(3)]
		return (binary64.Load(center) + uint64(rng.Intn(33)-16)) | (rng.Uint64() & (1 << 63))
	default:
		return binary64.Load(math.Ldexp(rng.Float64()-0.5, -rng.Intn(1074)))
	}
}

func observe(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(binary64.Load(x), 10)
}

// Compute expm1(2|x|) from its series at a scaled argument and double
// the argument repeatedly with expm1(2a) = expm1(a) * (expm1(a)+2).
func idealTanh(x float64) float64 {
	if math.Abs(x) > 32 {
		return math.Copysign(1, x)
	}
	precise := func(value float64) *big.Float {
		return new(big.Float).SetPrec(256).SetFloat64(value)
	}
	argument := precise(math.Abs(x))
	argument.SetMantExp(argument, -6)
	term, sum, divisor := precise(0).Set(argument), precise(0).Set(argument), precise(0)
	for k := int64(2); k <= 64; k++ {
		term.Mul(term, argument).Quo(term, divisor.SetInt64(k))
		sum.Add(sum, term)
	}
	two, factor := precise(2), precise(0)
	for range 7 {
		sum.Mul(sum, factor.Add(sum, two))
	}
	result, _ := sum.Quo(sum, factor.Add(sum, two)).Float64()
	return math.Copysign(result, x)
}

func nativeProfile() (int, error) {
	if runtime.GOARCH == "arm64" {
		return 4, nil
	}
	if runtime.GOARCH != "amd64" {
		return 0, fmt.Errorf("unmodeled Tanh architecture %s", runtime.GOARCH)
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return 0, fmt.Errorf("cannot determine Tanh build profile")
	}
	profile := 0
	for _, setting := range info.Settings {
		if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
			profile = 2
		}
	}
	dir, err := os.MkdirTemp("", "golean-tanh-exp-profile-")
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
