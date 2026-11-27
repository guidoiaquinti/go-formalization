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
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct{ x, y uint64 }

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
	words := binary64.Boundaries()
	for _, x := range []float64{0.5, 0.75, 1.5, 2, 3, 31, 32, 33, 4096, 0x1p52 + 1, 0x1p63} {
		words = append(words, binary64.Load(x), binary64.Load(-x))
	}
	var edges []input
	for _, x := range words {
		for _, y := range words {
			edges = append(edges, input{x, y})
		}
	}
	ln2 := atanhLog(precise(1).Quo(precise(1), precise(3)))
	var exactCases, numericalCases int
	var maxDistance uint64
	var fixture bytes.Buffer
	check.CompareSamples("godiff-math-pow", edges, sample, func(c input) string {
		return fmt.Sprintf("%d %d %d", profile, c.x, c.y)
	}, func(c input) string {
		return observe(math.Pow(binary64.Store(c.x), binary64.Store(c.y)))
	}, func(c input) string {
		x, y := binary64.Store(c.x), binary64.Store(c.y)
		if want, exact := exactResult(x, y); exact {
			exactCases++
			expected := observe(want)
			fmt.Fprintln(&fixture, c.x, c.y, expected, 0, false)
			return expected
		}
		got := math.Pow(x, y)
		var reference uint64
		bounded := math.Abs(x) >= 0.5 && math.Abs(x) <= 2 && math.Abs(y) <= 32
		if bounded {
			want := idealPow(x, y, ln2)
			a, b := binary64.Load(got), binary64.Load(want)
			distance := max(a, b) - min(a, b)
			if distance > 64 {
				log.Fatalf("Pow(%#016x, %#016x) bits %#016x differ from reference %#016x by %d steps", c.x, c.y, a, b, distance)
			}
			reference = b
			numericalCases++
			maxDistance = max(maxDistance, distance)
		}
		expected := observe(got)
		// The public wrapper selects pow on both modeled profiles, and
		// CompareSamples verifies that this word is also the Lean output.
		fmt.Fprintln(&fixture, c.x, c.y, expected, reference, bounded)
		return expected
	})
	fmt.Printf("Pow reference: %d exact special/identity cases and %d bounded numerical cases; maximum distance %d encoded steps.\n",
		exactCases, numericalCases, maxDistance)
	dir, err := os.MkdirTemp("", "golean-math-pow-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_POW_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/pow/testdata/"
	return check.TestInPackage("math", "TestGoLeanPow", fixtures+"bridge_test.go", fixtures+"pow_test.go")
}

func exactResult(x, y float64) (float64, bool) {
	switch {
	case y == 0 || x == 1:
		return 1, true
	case y == 1:
		return x, true
	case math.IsNaN(x) || math.IsNaN(y):
		return math.NaN(), true
	}
	odd := math.Mod(math.Abs(y), 2) == 1
	negative := math.Signbit(x) && odd
	switch {
	case x == 0 || math.IsInf(x, 0):
		value := 0.0
		if (x == 0) == (y < 0) {
			value = math.Inf(1)
		}
		if negative {
			value = -value
		}
		return value, true
	case math.IsInf(y, 0) || math.Abs(y) >= 0x1p63:
		magnitude := math.Abs(x)
		if magnitude == 1 {
			return 1, true
		}
		if (magnitude < 1) == (y > 0) {
			return 0, true
		}
		return math.Inf(1), true
	case x < 0 && y != math.Trunc(y):
		return math.NaN(), true
	}
	return 0, false
}

func sample(r *rand.Rand, i int) input {
	x, y := r.Uint64(), r.Uint64()
	switch i % 5 {
	case 0:
		x = binary64.Load(0.5 + 1.5*r.Float64())
		y = binary64.Load(64*r.Float64() - 32)
	case 1:
		x = binary64.Load(0.5+1.5*r.Float64()) | x&(1<<63)
		y = binary64.Load(float64(r.Intn(65) - 32))
	case 2:
		y = binary64.Load(float64(r.Intn(131073) - 65536))
	case 3:
		x = uint64(int64(binary64.Load(1)) + r.Int63n(8193) - 4096)
	}
	return input{x, y}
}

func observe(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(binary64.Load(x), 10)
}

func nativeProfile() (int, error) {
	switch runtime.GOARCH {
	case "arm64":
		return 2, nil
	case "amd64":
		dir, err := os.MkdirTemp("", "golean-math-exp-")
		if err != nil {
			return 0, err
		}
		defer os.RemoveAll(dir)
		path := filepath.Join(dir, "profile")
		if err := os.Setenv("GOLEAN_MATH_EXP_PROFILE", path); err != nil {
			return 0, err
		}
		const fixtures = "difftest/math/exp/testdata/"
		if err := check.TestInPackage("math", "TestGoLeanExpProfile", fixtures+"profile_amd64_test.go", fixtures+"observe_profile_amd64_test.go"); err != nil {
			return 0, err
		}
		value, err := os.ReadFile(path)
		if err != nil {
			return 0, err
		}
		fused, err := strconv.ParseBool(string(value))
		if fused {
			return 1, err
		}
		return 0, err
	default:
		return 0, fmt.Errorf("unmodeled Pow architecture %s", runtime.GOARCH)
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

func idealPow(x, y float64, ln2 *big.Float) float64 {
	if y == math.Trunc(y) {
		base := new(big.Rat).SetFloat64(x)
		power := new(big.Rat).SetInt64(1)
		for i := int(math.Abs(y)); i != 0; i >>= 1 {
			if i&1 != 0 {
				power.Mul(power, base)
			}
			base.Mul(base, base)
		}
		if y < 0 {
			power.Inv(power)
		}
		result, _ := power.Float64()
		return result
	}
	// Within [0.5, 2], the atanh transform has magnitude at most 1/3.
	ratio := precise(0).Quo(precise(0).Sub(precise(x), precise(1)), precise(0).Add(precise(x), precise(1)))
	power := atanhLog(ratio)
	power.Mul(power, precise(y))
	approximation, _ := power.Float64()
	k := math.Round(approximation * math.Log2E)
	r := power.Sub(power, precise(k).Mul(precise(k), ln2))
	term, sum, divisor := precise(1), precise(1), precise(0)
	for n := int64(1); n <= 64; n++ {
		term.Mul(term, r).Quo(term, divisor.SetInt64(n))
		sum.Add(sum, term)
	}
	result, _ := sum.SetMantExp(sum, int(k)).Float64()
	return result
}
