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
	dir, err := os.MkdirTemp("", "golean-math-sinh-")
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
	if err := os.Setenv("GOLEAN_MATH_SINH_CASES", path); err != nil {
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
	for _, center := range []float64{math.SmallestNonzeroFloat64, 0x1p-1022, 0x1p-28, 0.5, 21, 709,
		0x1.62e42fefa39efp+9, 710, 710.475860073944, 711, math.MaxFloat64} {
		for _, value := range []float64{math.Nextafter(center, 0), center, math.Nextafter(center, math.Inf(1))} {
			edges = append(edges, binary64.Load(value), binary64.Load(-value))
		}
	}
	var maximum uint64
	var premature int
	check.CompareSamples("godiff-math-sinh", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return observe(math.Sinh(binary64.Store(word)))
	}, func(word uint64) string {
		x := binary64.Store(word)
		switch {
		case math.IsNaN(x):
			return exact(word, "NaN")
		case x == 0 || math.IsInf(x, 0):
			return exact(word, observe(x))
		}
		got, want := math.Sinh(x), idealSinh(x)
		if math.Abs(x) > 21 {
			staged := math.Copysign(0.5*math.Exp(math.Abs(x)), x)
			if binary64.Load(got) != binary64.Load(staged) {
				log.Fatalf("Sinh(%#016x) = %g, pinned staged Exp = %g", word, got, staged)
			}
			if math.IsInf(got, 0) && !math.IsInf(want, 0) {
				premature++
				expected := observe(staged)
				fmt.Fprintln(fixture, word, expected, observe(want), "staged")
				return expected
			}
		}
		a, b := binary64.Load(got), binary64.Load(want)
		distance := max(a, b) - min(a, b)
		if math.IsInf(got, 0) != math.IsInf(want, 0) || distance > 5 {
			log.Fatalf("Sinh(%#016x) = %#016x, reference = %#016x, distance = %d", word, a, b, distance)
		}
		maximum = max(maximum, distance)
		expected := observe(got)
		fmt.Fprintln(fixture, word, expected, observe(want), "numeric")
		return expected
	})
	for _, x := range []float64{-710, 710} {
		if got, ideal := math.Sinh(x), idealSinh(x); !math.IsInf(got, 0) || math.IsInf(ideal, 0) {
			log.Fatalf("expected staged overflow discrepancy at %g: Go=%g, reference=%g", x, got, ideal)
		} else {
			fmt.Printf("Sinh(%g): pinned=%g, independent reference=%.17g (%#016x).\n", x, got, ideal, binary64.Load(ideal))
		}
	}
	fmt.Printf("Sinh 256-bit reference: maximum distance %d steps; %d sampled or boundary staged overflows with finite reference.\n",
		maximum, premature)
	if err := fixture.Flush(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	const data = "difftest/math/sinh/testdata/"
	return check.TestInPackage("math", "TestGoLeanSinh", data+"bridge_test.go", data+"sinh_test.go")
}

func sample(rng *rand.Rand, i int) uint64 {
	switch i % 6 {
	case 0:
		return rng.Uint64()
	case 1:
		return binary64.Load(rng.Float64() - 0.5)
	case 2:
		return binary64.Load(1422*rng.Float64() - 711)
	case 3:
		centers := []float64{0.5, 21, 0x1.62e42fefa39efp+9, 710.475860073944}
		word := binary64.Load(centers[rng.Intn(len(centers))])
		return (word + uint64(rng.Intn(129)-64)) | (rng.Uint64() & (1 << 63))
	case 4:
		return binary64.Load(math.Ldexp(rng.Float64()-0.5, -rng.Intn(1074)))
	default:
		return binary64.Load(math.Copysign(709+2*rng.Float64(), float64(rng.Intn(2)*2-1)))
	}
}

func precise(value float64) *big.Float {
	return new(big.Float).SetPrec(256).SetFloat64(value)
}

func idealSinh(x float64) float64 {
	magnitude := math.Abs(x)
	if magnitude >= 711 {
		return math.Copysign(math.Inf(1), x)
	}
	sum := precise(0)
	if magnitude <= 0.5 {
		// Odd Taylor series preserves tiny inputs without cancellation.
		power := precise(magnitude)
		square := precise(0).Mul(power, power)
		sum.Set(power)
		divisor := precise(0)
		for n := int64(3); n <= 65; n += 2 {
			power.Mul(power, square)
			power.Quo(power, divisor.SetInt64(n*(n-1)))
			sum.Add(sum, power)
		}
	} else {
		// Compute e^(|x|/1024) as a Taylor series and square ten times.
		reduced := precise(magnitude / 1024)
		term, divisor := precise(1), precise(0)
		sum.SetInt64(1)
		for n := int64(1); n <= 85; n++ {
			term.Mul(term, reduced)
			term.Quo(term, divisor.SetInt64(n))
			sum.Add(sum, term)
		}
		for range 10 {
			sum.Mul(sum, sum)
		}
		sum.Sub(sum, precise(0).Quo(precise(1), sum))
		sum.SetMantExp(sum, -1)
	}
	result, _ := sum.Float64()
	return math.Copysign(result, x)
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
		return 0, fmt.Errorf("unmodeled Sinh architecture %s", runtime.GOARCH)
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return 0, fmt.Errorf("cannot determine Sinh build profile")
	}
	profile := 0
	for _, setting := range info.Settings {
		if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
			profile = 2
		}
	}
	dir, err := os.MkdirTemp("", "golean-sinh-exp-profile-")
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
