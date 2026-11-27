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
	"strings"

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
		return fmt.Errorf("unmodeled Log10 architecture %s", runtime.GOARCH)
	}
	edges := binary64.Boundaries()
	for k := -323; k <= 308; k++ {
		x := math.Pow10(k)
		edges = append(edges, binary64.Load(math.Nextafter(x, 0)), binary64.Load(x), binary64.Load(math.Nextafter(x, math.Inf(1))))
	}
	ln2 := atanhLog(precise(1).Quo(precise(1), precise(3)))
	// 10 = 8*1.25; (1.25 - 1)/(1.25 + 1) = 1/9.
	ln10 := atanhLog(precise(1).Quo(precise(1), precise(9)))
	ln10.Add(ln10, precise(3).Mul(precise(3), ln2))
	checkAMD64SubnormalOracle(ln2, ln10)

	var maxDistance, subnormalDifferences uint64
	var fixture bytes.Buffer
	check.CompareSamples("godiff-math-log10", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return observe(math.Log10(binary64.Store(word)))
	}, func(word uint64) string {
		x := binary64.Store(word)
		if expected, exact := specialResult(x); exact {
			fmt.Fprintln(&fixture, word, expected, 0, true)
			return expected
		}
		got, want := math.Log10(x), idealLog10(x, 0, ln2, ln10)
		if profile == 0 && word < 1<<52 {
			if binary64.Load(got) != binary64.Load(want) {
				subnormalDifferences++
			}
			want = idealLog10(binary64.Store(word|1<<52), -1, ln2, ln10)
		}
		a, b := binary64.Load(got), binary64.Load(want)
		distance := max(a, b) - min(a, b)
		if distance > 2 {
			log.Fatalf("Log10(%#016x) bits %#016x differ from reference %#016x by %d steps", word, a, b, distance)
		}
		maxDistance = max(maxDistance, distance)
		expected := observe(got)
		// On these profiles the public wrapper calls log10, and the public
		// comparison guarantees that this word also equals the Lean result.
		fmt.Fprintln(&fixture, word, expected, b, false)
		return expected
	})
	fmt.Printf("Log10 high-precision reference: maximum distance %d encoded steps; %d amd64 subnormal differences from mathematical inputs.\n",
		maxDistance, subnormalDifferences)
	dir, err := os.MkdirTemp("", "golean-math-log10-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_LOG10_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/log10/testdata/"
	return check.TestInPackage("math", "TestGoLeanLog10", fixtures+"bridge_test.go", fixtures+"log10_test.go")
}

func specialResult(x float64) (string, bool) {
	switch {
	case math.IsNaN(x) || x < 0:
		return "NaN", true
	case x == 0:
		return observe(math.Inf(-1)), true
	case math.IsInf(x, 1):
		return observe(x), true
	case x == 1:
		return "0", true
	default:
		return "", false
	}
}

func checkAMD64SubnormalOracle(ln2, ln10 *big.Float) {
	rows := check.LeanRows("godiff-math-log10", 1, strings.NewReader("0 1\n"))
	got, err := strconv.ParseUint(rows[0], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	want := binary64.Load(idealLog10(binary64.Store(0x0010000000000001), -1, ln2, ln10))
	if distance := max(got, want) - min(got, want); distance > 2 {
		log.Fatalf("simulated amd64 Log10 of smallest subnormal: Lean=%#016x reference=%#016x distance=%d", got, want, distance)
	}
	fmt.Printf("Simulated amd64 Log10 of smallest subnormal: Lean=%#016x reference=%#016x.\n", got, want)
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
		return binary64.Load(math.Pow10(r.Intn(632)-323)) + uint64(r.Intn(9)-4)
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

func idealLog10(x float64, adjustExponent int, ln2, ln10 *big.Float) float64 {
	u, mantissa := precise(x), precise(0)
	exponent := u.MantExp(mantissa)
	// Center the mantissa around 1 to avoid cancellation in the sum near 1.
	if mantissa.Cmp(precise(0.75)) < 0 {
		exponent--
		mantissa.SetMantExp(mantissa, 1)
	}
	ratio := precise(0).Quo(precise(0).Sub(mantissa, precise(1)), precise(0).Add(mantissa, precise(1)))
	result := atanhLog(ratio)
	exponent += adjustExponent
	result.Add(result, precise(float64(exponent)).Mul(precise(float64(exponent)), ln2))
	rounded, _ := result.Quo(result, ln10).Float64()
	return rounded
}
