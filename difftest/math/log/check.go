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
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type oracleRow struct {
	word, reference uint64
	exact           bool
}

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
		return fmt.Errorf("unmodeled Log architecture %s", runtime.GOARCH)
	}
	edges := binary64.Boundaries()
	for _, x := range []float64{0.5, math.Sqrt2 / 2, 1, math.Sqrt2, 2, 0x1p-1022, 0x1p-1021} {
		for _, word := range []float64{math.Nextafter(x, 0), x, math.Nextafter(x, math.Inf(1))} {
			edges = append(edges, binary64.Load(word), binary64.Load(-word))
		}
	}
	for k := -1074; k <= 1023; k++ {
		edges = append(edges, binary64.Load(math.Ldexp(1, k)))
	}
	ln2 := atanhLog(precise(1).Quo(precise(1), precise(3)))
	checkAMD64SubnormalOracle(ln2)
	var maxDistance, subnormalDifferences uint64
	var cases []oracleRow
	check.CompareSamples("godiff-math-log", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return observe(math.Log(binary64.Store(word)))
	}, func(word uint64) string {
		x := binary64.Store(word)
		if value, exact := specialLog(x); exact {
			cases = append(cases, oracleRow{word: word, exact: true})
			return value
		}
		got, want := math.Log(x), idealLog(x, 0, ln2)
		cases = append(cases, oracleRow{word: word, reference: binary64.Load(want)})
		if profile == 0 && word < 1<<52 {
			if binary64.Load(got) != binary64.Load(want) {
				subnormalDifferences++
			}
			// This assembly treats a zero exponent field as if a leading
			// bit were present and reads it with exponent -1023. Test
			// the effective input without rounding the division by two.
			want = idealLog(binary64.Store(word|1<<52), -1, ln2)
		}
		a, b := binary64.Load(got), binary64.Load(want)
		distance := max(a, b) - min(a, b)
		if distance > 1 {
			log.Fatalf("Log(%#016x) bits %#016x differ from reference %#016x by %d steps", word, a, b, distance)
		}
		maxDistance = max(maxDistance, distance)
		return observe(got)
	})
	fmt.Printf("Log high-precision reference: maximum distance %d encoded steps; %d amd64 subnormal differences from mathematical inputs.\n",
		maxDistance, subnormalDifferences)
	return compareGeneric(cases)
}

func specialLog(x float64) (string, bool) {
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

func compareGeneric(cases []oracleRow) error {
	profile := 4
	if runtime.GOARCH == "amd64" {
		profile = 2
		info, ok := debug.ReadBuildInfo()
		if !ok {
			return fmt.Errorf("cannot determine generic Log build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				profile = 3
			}
		}
	}

	var modelInput, fixture bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&modelInput, profile, c.word)
	}
	for i, model := range check.LeanRows("godiff-math-log", len(cases), &modelInput) {
		c := cases[i]
		if expected, exact := specialLog(binary64.Store(c.word)); exact && model != expected {
			return fmt.Errorf("input %#016x: generic Lean=%s expected=%s", c.word, model, expected)
		}
		fmt.Fprintln(&fixture, c.word, model, c.reference, c.exact)
	}
	dir, err := os.MkdirTemp("", "golean-math-log-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_LOG_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/log/testdata/"
	return check.TestInPackage("math", "TestGoLeanLog", fixtures+"bridge_test.go", fixtures+"log_test.go")

}

// Run this check on every platform so the effective-operand reference is
// exercised even when native amd64 Go execution is unavailable.
func checkAMD64SubnormalOracle(ln2 *big.Float) {
	rows := check.LeanRows("godiff-math-log", 1, strings.NewReader("0 1\n"))
	got, err := strconv.ParseUint(rows[0], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	want := binary64.Load(idealLog(binary64.Store(0x0010000000000001), -1, ln2))
	if distance := max(got, want) - min(got, want); distance > 1 {
		log.Fatalf("simulated amd64 Log of smallest subnormal: Lean=%#016x reference=%#016x distance=%d", got, want, distance)
	}
	fmt.Printf("Simulated amd64 Log of smallest subnormal: Lean=%#016x reference=%#016x.\n", got, want)
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

func idealLog(x float64, adjustExponent int, ln2 *big.Float) float64 {
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
	rounded, _ := result.Float64()
	return rounded
}
