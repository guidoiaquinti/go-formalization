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
	profile, err := nativeProfile()
	if err != nil {
		return err
	}
	edges := binary64.Boundaries()
	for _, x := range []float64{-746, -745.1332191019411, -745, -744, -708.3964185322641,
		-1, -0x1p-28, -0x1p-53, 0, 0x1p-53, 0x1p-28, 1,
		1023.5 * math.Ln2, 709.5, 709.782712893384} {
		for _, value := range []float64{math.Nextafter(x, math.Inf(-1)), x, math.Nextafter(x, math.Inf(1))} {
			edges = append(edges, binary64.Load(value))
		}
	}
	ln2 := highPrecisionLn2()
	var maxDistance uint64
	var earlyOverflows int
	var cases []oracleRow
	check.CompareSamples("godiff-math-exp", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return strconv.FormatUint(binary64.Load(math.Exp(binary64.Store(word))), 10)
	}, func(word uint64) string {
		x := binary64.Store(word)
		want, exact := reference(word, ln2)
		cases = append(cases, oracleRow{word, want, exact})
		if exact {
			return strconv.FormatUint(want, 10)
		}
		got := math.Exp(x)
		if profile != 2 && math.IsInf(got, 1) && want != binary64.Load(math.Inf(1)) &&
			math.RoundToEven(x*math.Log2E) == 1024 {
			earlyOverflows++
			return strconv.FormatUint(binary64.Load(got), 10)
		}
		a, b := binary64.Load(got), want
		distance := max(a, b) - min(a, b)
		if distance > 2 {
			log.Fatalf("Exp(%#016x) bits %#016x differ from reference %#016x by %d steps", word, a, b, distance)
		}
		maxDistance = max(maxDistance, distance)
		return strconv.FormatUint(a, 10)
	})
	fmt.Printf("Exp high-precision reference: maximum distance %d encoded steps; %d amd64 early-overflow cases.\n", maxDistance, earlyOverflows)
	return compareGeneric(cases)
}

func reference(word uint64, ln2 *big.Float) (uint64, bool) {
	x := binary64.Store(word)
	switch {
	case math.IsNaN(x) || math.IsInf(x, 1):
		return word, true
	case x < -746:
		return 0, true
	case x > 710:
		return binary64.Load(math.Inf(1)), true
	case x == 0:
		return binary64.Load(1), true
	default:
		return binary64.Load(idealExp(x, ln2)), false
	}
}

func compareGeneric(cases []oracleRow) error {
	profile := 5
	if runtime.GOARCH == "amd64" {
		profile = 3
		info, ok := debug.ReadBuildInfo()
		if !ok {
			return fmt.Errorf("cannot determine generic Exp build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				profile = 4
			}
		}
	}

	var modelInput, fixture bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&modelInput, profile, c.word)
	}
	for i, model := range check.LeanRows("godiff-math-exp", len(cases), &modelInput) {
		c := cases[i]
		fmt.Fprintln(&fixture, c.word, model, c.reference, c.exact)
	}
	dir, err := os.MkdirTemp("", "golean-math-exp-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_EXP_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/exp/testdata/"
	return check.TestInPackage("math", "TestGoLeanExp", fixtures+"bridge_test.go", fixtures+"exp_test.go")

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
		return 0, fmt.Errorf("unmodeled Exp architecture %s", runtime.GOARCH)
	}
}

func sample(r *rand.Rand, i int) uint64 {
	switch i % 5 {
	case 0:
		return r.Uint64()
	case 1:
		return binary64.Load(r.Float64()*1457 - 746)
	case 2:
		k := float64(r.Intn(2100) - 1075)
		return binary64.Load((k+0.5)*math.Ln2) + uint64(r.Intn(9)-4)
	case 3:
		return binary64.Load(-746 + 38*r.Float64())
	default:
		return binary64.Load(708 + 3*r.Float64())
	}
}

const precision = 256

func highPrecisionLn2() *big.Float {
	value := func(x float64) *big.Float { return new(big.Float).SetPrec(precision).SetFloat64(x) }
	third := value(1).Quo(value(1), value(3))
	ninth := value(1).Mul(third, third)
	power, sum := value(1).Set(third), value(0)
	term, divisor := value(0), value(0)
	// ln(2) = 2 * atanh(1/3).
	for n := int64(1); n <= 255; n += 2 {
		sum.Add(sum, term.Quo(power, divisor.SetInt64(n)))
		power.Mul(power, ninth)
	}
	return sum.Mul(sum, value(2))
}

func idealExp(x float64, ln2 *big.Float) float64 {
	value := func(x float64) *big.Float { return new(big.Float).SetPrec(precision).SetFloat64(x) }
	k := math.Round(x * math.Log2E)
	r := value(x).Sub(value(x), value(k).Mul(value(k), ln2))
	term, sum, divisor := value(1), value(1), value(0)
	for n := int64(1); n <= 64; n++ {
		term.Mul(term, r).Quo(term, divisor.SetInt64(n))
		sum.Add(sum, term)
	}
	result, _ := sum.SetMantExp(sum, int(k)).Float64()
	return result
}
