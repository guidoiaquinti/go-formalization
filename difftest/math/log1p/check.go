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
	profile := 0
	switch runtime.GOARCH {
	case "arm64":
		profile = 2
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			return fmt.Errorf("cannot determine Log1p build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				profile = 1
			}
		}
	default:
		return fmt.Errorf("unmodeled Log1p architecture %s", runtime.GOARCH)
	}

	edges := binary64.Boundaries()
	for _, value := range []float64{-1, -0x1p-54, 0x1p-54, -0x1p-29, 0x1p-29,
		-0.2928932188134524756, 0.4142135623730950488, 0.5, 1, 3, 0x1p53} {
		for _, x := range []float64{math.Nextafter(value, math.Inf(-1)), value, math.Nextafter(value, math.Inf(1))} {
			edges = append(edges, binary64.Load(x))
		}
	}
	for k := -52; k <= 1023; k++ {
		value := math.Ldexp(1, k) - 1
		edges = append(edges, binary64.Load(value))
	}
	ln2 := atanhLog(precise(1).Quo(precise(1), precise(3)))
	var maxDistance uint64
	var fixture bytes.Buffer
	check.CompareSamples("godiff-math-log1p", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return observe(math.Log1p(binary64.Store(word)))
	}, func(word uint64) string {
		x := binary64.Store(word)
		if expected, exact := specialResult(x); exact {
			fmt.Fprintln(&fixture, word, expected, 0, true)
			return expected
		}
		got, want := math.Log1p(x), idealLog1p(x, ln2)
		a, b := binary64.Load(got), binary64.Load(want)
		distance := max(a, b) - min(a, b)
		if distance > 1 {
			log.Fatalf("Log1p(%#016x) bits %#016x differ from reference %#016x by %d steps", word, a, b, distance)
		}
		maxDistance = max(maxDistance, distance)
		expected := observe(got)
		// The public wrapper routes to log1p on both modeled profiles, and
		// CompareSamples confirms that this row also equals the Lean output.
		fmt.Fprintln(&fixture, word, expected, b, false)
		return expected
	})
	fmt.Printf("Log1p high-precision reference: maximum encoded-step distance %d.\n", maxDistance)
	dir, err := os.MkdirTemp("", "golean-math-log1p-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_LOG1P_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/log1p/testdata/"
	return check.TestInPackage("math", "TestGoLeanLog1p", fixtures+"bridge_test.go", fixtures+"log1p_test.go")
}

func specialResult(x float64) (string, bool) {
	switch {
	case math.IsNaN(x) || x < -1:
		return "NaN", true
	case x == -1:
		return observe(math.Inf(-1)), true
	case math.IsInf(x, 1) || x == 0:
		return observe(x), true
	default:
		return "", false
	}
}

func observe(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(binary64.Load(x), 10)
}

func sample(r *rand.Rand, i int) uint64 {
	switch i % 5 {
	case 0:
		return r.Uint64()
	case 1:
		return binary64.Load(2*r.Float64() - 1)
	case 2:
		k := r.Intn(1076) - 52
		value := math.Ldexp(1, k) - 1
		return binary64.Load(value) + uint64(r.Intn(9)-4)
	case 3:
		return binary64.Load(math.Ldexp(r.Float64()-0.5, -r.Intn(1074)))
	default:
		return binary64.Load(-1 + math.Ldexp(r.Float64(), -r.Intn(54)))
	}
}

const precision = 256

func precise(value float64) *big.Float {
	return new(big.Float).SetPrec(precision).SetFloat64(value)
}

// Compute 2*atanh(s) with |s| at most 1/3.
func atanhLog(s *big.Float) *big.Float {
	square := precise(0).Mul(s, s)
	power, sum := precise(0).Set(s), precise(0)
	term, divisor := precise(0), precise(0)
	for n := int64(1); n <= 255; n += 2 {
		sum.Add(sum, term.Quo(power, divisor.SetInt64(n)))
		power.Mul(power, square)
	}
	return sum.Mul(sum, precise(2))
}

func idealLog1p(x float64, ln2 *big.Float) float64 {
	var exponent int
	var ratio *big.Float
	if math.Abs(x) <= 0.5 {
		// log(1+x) = 2*atanh(x/(2+x)) avoids rounding away tiny x.
		ratio = precise(x).Quo(precise(x), precise(2).Add(precise(2), precise(x)))
	} else {
		u, mantissa := precise(1).Add(precise(1), precise(x)), precise(0)
		exponent = u.MantExp(mantissa) - 1
		mantissa.SetMantExp(mantissa, 1)
		ratio = precise(0).Quo(precise(0).Sub(mantissa, precise(1)), precise(0).Add(mantissa, precise(1)))
	}
	value := atanhLog(ratio)
	if exponent != 0 {
		value.Add(value, precise(float64(exponent)).Mul(precise(float64(exponent)), ln2))
	}
	rounded, _ := value.Float64()
	return rounded
}
