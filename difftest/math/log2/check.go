package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/bits"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
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
		profile = 1
	case "amd64":
	default:
		return fmt.Errorf("unmodeled Log2 architecture %s", runtime.GOARCH)
	}
	edges := binary64.Boundaries()
	for k := -1074; k <= 1023; k++ {
		x := math.Ldexp(1, k)
		edges = append(edges, binary64.Load(math.Nextafter(x, 0)), binary64.Load(x), binary64.Load(math.Nextafter(x, math.Inf(1))))
	}
	ln2 := atanhLog(precise(1).Quo(precise(1), precise(3)))
	var maxScaledError float64
	var maxStepsAwayFromZero uint64
	var fixture bytes.Buffer
	check.CompareSamples("godiff-math-log2", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return observe(math.Log2(binary64.Store(word)))
	}, func(word uint64) string {
		x := binary64.Store(word)
		if expected, exact := specialResult(word, x); exact {
			fmt.Fprintln(&fixture, word, expected, 0, true)
			return expected
		}
		got, want := math.Log2(x), idealLog2(x, ln2)
		scaledError := math.Abs(got-want) / max(1, math.Abs(want))
		if math.IsNaN(scaledError) || scaledError > 0x1p-50 {
			log.Fatalf("Log2(%#016x) = %.17g, reference = %.17g, scaled error = %g", word, got, want, scaledError)
		}
		maxScaledError = max(maxScaledError, scaledError)
		if math.Abs(want) >= 1 {
			a, b := binary64.Load(got), binary64.Load(want)
			distance := max(a, b) - min(a, b)
			if distance > 2 {
				log.Fatalf("Log2(%#016x) bits %#016x differ from reference %#016x by %d steps", word, a, b, distance)
			}
			maxStepsAwayFromZero = max(maxStepsAwayFromZero, distance)
		}
		expected := observe(got)
		// On both modeled profiles, the public wrapper selects log2, and
		// CompareSamples checks that this result also equals the Lean row.
		fmt.Fprintln(&fixture, word, expected, binary64.Load(want), false)
		return expected
	})
	fmt.Printf("Log2 high-precision reference: max scaled absolute error %g; max encoded steps where |result| >= 1 is %d.\n",
		maxScaledError, maxStepsAwayFromZero)
	dir, err := os.MkdirTemp("", "golean-math-log2-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_LOG2_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/log2/testdata/"
	return check.TestInPackage("math", "TestGoLeanLog2", fixtures+"bridge_test.go", fixtures+"log2_test.go")
}

func specialResult(word uint64, x float64) (string, bool) {
	switch {
	case math.IsNaN(x) || x < 0:
		return "NaN", true
	case x == 0:
		return observe(math.Inf(-1)), true
	case math.IsInf(x, 1):
		return observe(x), true
	}
	if power, ok := powerOfTwoExponent(word); ok {
		return observe(float64(power)), true
	}
	return "", false
}

func powerOfTwoExponent(word uint64) (int, bool) {
	exponent := int(word >> 52)
	fraction := word & (1<<52 - 1)
	if exponent != 0 {
		return exponent - 1023, fraction == 0
	}
	return bits.TrailingZeros64(fraction) - 1074, fraction&(fraction-1) == 0
}

func sample(r *rand.Rand, i int) uint64 {
	switch i % 5 {
	case 0:
		return r.Uint64()
	case 1:
		return r.Uint64() & (1<<52 - 1)
	case 2:
		return uint64(int64(binary64.Load(1)) + r.Int63n(8193) - 4096)
	case 3:
		k := r.Intn(2098) - 1074
		return binary64.Load(math.Ldexp(math.Sqrt2, k)) + uint64(r.Intn(9)-4)
	default:
		return r.Uint64() &^ (1 << 63)
	}
}

func observe(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(binary64.Load(x), 10)
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

func idealLog2(x float64, ln2 *big.Float) float64 {
	u, mantissa := precise(x), precise(0)
	exponent := u.MantExp(mantissa)
	// Center the mantissa around 1 to avoid cancellation in the sum near 1.
	if mantissa.Cmp(precise(0.75)) < 0 {
		exponent--
		mantissa.SetMantExp(mantissa, 1)
	}
	ratio := precise(0).Quo(precise(0).Sub(mantissa, precise(1)), precise(0).Add(mantissa, precise(1)))
	result := atanhLog(ratio)
	result.Quo(result, ln2).Add(result, precise(float64(exponent)))
	rounded, _ := result.Float64()
	return rounded
}
