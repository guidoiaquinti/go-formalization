package main

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"runtime"
	"runtime/debug"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	profile := nativeProfile()
	edges := binary64.Boundaries()
	for _, center := range []float64{math.SmallestNonzeroFloat64, 0x1p-1022, 0x1p-28, 0.5, 0.85,
		0.9, 1 - 2*math.Exp(-25), 1, math.MaxFloat64} {
		for _, value := range []float64{math.Nextafter(center, 0), center, math.Nextafter(center, math.Inf(1))} {
			edges = append(edges, binary64.Load(value), binary64.Load(-value))
		}
	}
	pi, _, err := big.ParseFloat(
		"3.141592653589793238462643383279502884197169399375105820974944592307816406286208998628034825342117068",
		10, 256, big.ToNearestEven)
	if err != nil {
		log.Fatal(err)
	}
	scale := precise(2).Quo(precise(2), pi.Sqrt(pi))
	var numericCases int
	const radius = 10
	check.CompareSamples("godiff-math-erfinv", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return observe(math.Erfinv(binary64.Store(word)))
	}, func(word uint64) string {
		x := binary64.Store(word)
		switch {
		case math.IsNaN(x) || math.Abs(x) > 1:
			return "NaN"
		case x == 0:
			return observe(x)
		case math.Abs(x) == 1:
			return observe(math.Copysign(math.Inf(1), x))
		}
		got := math.Erfinv(x)
		if math.IsNaN(got) || math.IsInf(got, 0) || math.Signbit(got) != math.Signbit(x) ||
			!bracket(x, got, radius, scale) {
			log.Fatalf("Erfinv(%#016x) = %#016x: independent Erf values do not bracket input within %d steps",
				word, binary64.Load(got), radius)
		}
		numericCases++
		return observe(got)
	})
	fmt.Printf("Erfinv 256-bit forward Erf brackets held within %d encoded steps on %d finite interior cases.\n",
		radius, numericCases)
}

func sample(rng *rand.Rand, i int) uint64 {
	const sign = uint64(1) << 63
	switch i % 6 {
	case 0:
		return rng.Uint64()
	case 1:
		return binary64.Load(2*rng.Float64() - 1)
	case 2:
		center := []float64{0.85, 1 - 2*math.Exp(-25), 1}[rng.Intn(3)]
		return (binary64.Load(center) + uint64(rng.Intn(129)-64)) | (rng.Uint64() & sign)
	case 3:
		return (binary64.Load(1) - 1 - uint64(rng.Intn(1<<21))) | (rng.Uint64() & sign)
	case 4:
		return binary64.Load(math.Ldexp(rng.Float64()-0.5, -rng.Intn(1074)))
	default:
		return binary64.Load(math.Copysign(1-math.Ldexp(rng.Float64(), -rng.Intn(55)), float64(rng.Intn(2)*2-1)))
	}
}

func precise(value float64) *big.Float {
	return new(big.Float).SetPrec(256).SetFloat64(value)
}

func exactErf(x float64, scale *big.Float) *big.Float {
	if x == 0 {
		return precise(0)
	}
	// Integrate the expansion of exp(-t*t) over [0, |x|].
	power := precise(math.Abs(x))
	ratio := precise(0).Mul(power, power)
	ratio.Neg(ratio)
	sum, coefficient, term, divisor := precise(0).Set(power), precise(0), precise(0), precise(0)
	exponent := power.MantExp(nil)
	for n := int64(1); n <= 224; n++ {
		power.Mul(power, ratio)
		power.Quo(power, coefficient.SetInt64(n))
		term.Quo(power, divisor.SetInt64(2*n+1))
		sum.Add(sum, term)
		// After n > 36 the terms decrease for the full |x| < 6 domain.
		if n > 36 && term.MantExp(nil) < exponent-210 {
			break
		}
	}
	sum.Mul(sum, scale)
	return sum
}

// The forward integral is strictly increasing on the nonnegative axis.
func bracket(x, candidate float64, steps uint64, scale *big.Float) bool {
	target := precise(math.Abs(x))
	word := math.Float64bits(math.Abs(candidate))
	lo := math.Float64frombits(word - min(word, steps))
	hi := math.Float64frombits(word + steps)
	return exactErf(lo, scale).Cmp(target) <= 0 && exactErf(hi, scale).Cmp(target) >= 0
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
			log.Fatal("cannot determine Erfinv build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		log.Fatalf("unmodeled Erfinv architecture %s", runtime.GOARCH)
		return 0
	}
}
