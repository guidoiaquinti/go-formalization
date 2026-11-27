package main

import (
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
	profile, err := nativeProfile()
	if err != nil {
		log.Fatal(err)
	}
	pi, _, err := big.ParseFloat(
		"3.141592653589793238462643383279502884197169399375105820974944592307816406286208998628034825342117068",
		10, 256, big.ToNearestEven)
	if err != nil {
		log.Fatal(err)
	}
	ln2 := atanhLog(precise(0).Quo(precise(1), precise(3)))
	sqrtTwoPi := precise(0).Sqrt(precise(0).SetMantExp(pi, 1))
	edges := binary64.Boundaries()
	for _, center := range []float64{math.SmallestNonzeroFloat64, 0x1p-1022, 1e-9, 0.5, 1, 2, 3, 32.5,
		33, 33.5, 143.01608, 171.5, 171.6243769563027, 172, 175.5, 180.5, 200, 0x1p52, math.MaxFloat64} {
		for _, x := range []float64{math.Nextafter(center, 0), center, math.Nextafter(center, math.Inf(1))} {
			edges = append(edges, binary64.Load(x), binary64.Load(-x))
		}
	}
	var maximumRelative float64
	var maximumSubnormal uint64
	check.CompareSamples("godiff-math-gamma", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return observe(math.Gamma(binary64.Store(word)))
	}, func(word uint64) string {
		x := binary64.Store(word)
		switch {
		case math.IsNaN(x) || math.IsInf(x, -1) || x < 0 && math.Trunc(x) == x:
			return "NaN"
		case x == 0 || math.IsInf(x, 1):
			return observe(math.Copysign(math.Inf(1), x))
		}
		got, want := math.Gamma(x), idealGamma(x, ln2, pi, sqrtTwoPi)
		if math.IsNaN(got) || math.Signbit(want) != math.Signbit(got) || math.IsInf(got, 0) != math.IsInf(want, 0) {
			log.Fatalf("Gamma(%#016x) class/sign mismatch Go=%g reference=%g", word, got, want)
		}
		if math.IsInf(got, 0) {
			return observe(got)
		}
		if math.Abs(want) < 0x1p-1022 {
			a, b := binary64.Load(got), binary64.Load(want)
			distance := max(a, b) - min(a, b)
			if distance > 128 {
				log.Fatalf("Gamma(%#016x) Go=%#016x reference=%#016x distance=%d", word, a, b, distance)
			}
			maximumSubnormal = max(maximumSubnormal, distance)
		} else {
			err := math.Abs(got/want - 1)
			if err > 0x1p-42 {
				log.Fatalf("Gamma(%#016x) Go=%g reference=%g relative error=%g", word, got, want, err)
			}
			maximumRelative = max(maximumRelative, err)
		}
		return observe(got)
	})
	fmt.Printf("Gamma 256-bit reference: maximum relative difference %g; subnormal distance %d encoded steps.\n",
		maximumRelative, maximumSubnormal)
}

func sample(r *rand.Rand, i int) uint64 {
	switch i % 6 {
	case 0:
		return r.Uint64()
	case 1:
		return binary64.Load(66*r.Float64() - 33)
	case 2:
		return binary64.Load(400*r.Float64() - 200)
	case 3:
		center := []float64{1e-9, 2, 33, 143.01608, 172, 200}[r.Intn(6)]
		word := binary64.Load(center) + uint64(r.Intn(129)-64)
		return word | (r.Uint64() & (1 << 63))
	case 4:
		return binary64.Load(math.Ldexp(r.Float64()-0.5, -r.Intn(1074)))
	default:
		n := 1 + r.Intn(200)
		return binary64.Load(-float64(n)) + uint64(r.Intn(129)-64)
	}
}

func precise(value float64) *big.Float {
	return new(big.Float).SetPrec(256).SetFloat64(value)
}

func atanhLog(x *big.Float) *big.Float {
	square := precise(0).Mul(x, x)
	power, sum := precise(0).Set(x), precise(0)
	term, divisor := precise(0), precise(0)
	for k := int64(1); k <= 159; k += 2 {
		sum.Add(sum, term.Quo(power, divisor.SetInt64(k)))
		power.Mul(power, square)
	}
	return sum.SetMantExp(sum, 1)
}

func logarithm(x, ln2 *big.Float) *big.Float {
	fraction := precise(0)
	k := x.MantExp(fraction) - 1
	fraction.SetMantExp(fraction, 1)
	ratio := precise(0).Quo(precise(0).Sub(fraction, precise(1)), precise(0).Add(fraction, precise(1)))
	fractionLog := atanhLog(ratio)
	return fractionLog.Add(fractionLog, precise(0).Mul(precise(float64(k)), ln2))
}

func exponential(x, ln2 *big.Float) *big.Float {
	approximation, _ := x.Float64()
	k := math.Round(approximation * math.Log2E)
	remainder := precise(0).Sub(x, precise(0).Mul(precise(k), ln2))
	term, sum, divisor := precise(1), precise(1), precise(0)
	for n := int64(1); n <= 64; n++ {
		term.Mul(term, remainder)
		term.Quo(term, divisor.SetInt64(n))
		sum.Add(sum, term)
	}
	return sum.SetMantExp(sum, int(k))
}

func correction(x *big.Float) *big.Float {
	reciprocal := precise(0).Quo(precise(1), x)
	square := precise(0).Mul(reciprocal, reciprocal)
	coef := func(n, d int64) *big.Float {
		return precise(0).Quo(precise(0).SetInt64(n), precise(0).SetInt64(d))
	}
	sum := coef(-691, 360360)
	for _, c := range [][2]int64{{1, 1188}, {-1, 1680}, {1, 1260}, {-1, 360}, {1, 12}} {
		sum.Mul(sum, square)
		sum.Add(sum, coef(c[0], c[1]))
	}
	return sum.Mul(sum, reciprocal)
}

func gammaPositive(x float64, ln2, sqrtTwoPi *big.Float) *big.Float {
	argument, prod, one := precise(x), precise(1), precise(1)
	for argument.Cmp(precise(36)) < 0 {
		prod.Mul(prod, argument)
		argument.Add(argument, one)
	}
	exponent := precise(0).Sub(argument, precise(0.5))
	exponent.Mul(exponent, logarithm(argument, ln2))
	exponent.Sub(exponent, argument)
	exponent.Add(exponent, correction(argument))
	value := exponential(exponent, ln2)
	value.Mul(value, sqrtTwoPi)
	return value.Quo(value, prod)
}

func sinPiMagnitude(x, floor float64, pi *big.Float) *big.Float {
	fraction := precise(0).Sub(precise(x), precise(floor))
	if fraction.Cmp(precise(0.5)) > 0 {
		fraction.Sub(precise(1), fraction)
	}
	angle := fraction.Mul(fraction, pi)
	square := precise(0).Mul(angle, angle)
	square.Neg(square)
	sum, term, divisor := precise(0).Set(angle), precise(0).Set(angle), precise(0)
	for n := int64(3); n <= 79; n += 2 {
		term.Mul(term, square)
		term.Quo(term, divisor.SetInt64(n*(n-1)))
		sum.Add(sum, term)
	}
	return sum
}

func idealGamma(x float64, ln2, pi, sqrtTwoPi *big.Float) float64 {
	magnitude := math.Abs(x)
	if x > 172 {
		return math.Inf(1)
	}
	if x < 0 {
		floor := math.Floor(magnitude)
		sign := float64(1)
		if int64(floor)&1 == 0 {
			sign = -1
		}
		if magnitude > 200 {
			return math.Copysign(0, sign)
		}
		value := gammaPositive(magnitude, ln2, sqrtTwoPi)
		value.Mul(value, precise(magnitude))
		value.Mul(value, sinPiMagnitude(magnitude, floor, pi))
		value.Quo(pi, value)
		rounded, _ := value.Float64()
		return math.Copysign(rounded, sign)
	}
	value := gammaPositive(magnitude, ln2, sqrtTwoPi)
	rounded, _ := value.Float64()
	return rounded
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
		return 0, fmt.Errorf("unmodeled Gamma architecture %s", runtime.GOARCH)
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return 0, fmt.Errorf("cannot determine Gamma build profile")
	}
	profile := 0
	for _, setting := range info.Settings {
		if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
			profile = 2
		}
	}
	dir, err := os.MkdirTemp("", "golean-gamma-exp-profile-")
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
