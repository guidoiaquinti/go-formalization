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
	profile := 0
	switch runtime.GOARCH {
	case "arm64":
		profile = 2
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			return fmt.Errorf("cannot determine Exp2 build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				profile = 1
			}
		}
	default:
		return fmt.Errorf("unmodeled Exp2 architecture %s", runtime.GOARCH)
	}

	edges := binary64.Boundaries()
	for _, x := range []float64{-1075, -1074.75, -1074.5, -1074.25, -1074, -1073.5, -1073, -1022, -1, -0.5, 0, 0.5, 1, 1023, 1023.5, 1024} {
		edges = append(edges, binary64.Load(math.Nextafter(x, math.Inf(-1))),
			binary64.Load(x), binary64.Load(math.Nextafter(x, math.Inf(1))))
	}
	for k := -1074; k <= 1023; k++ {
		edges = append(edges, binary64.Load(float64(k)))
	}

	ln2 := highPrecisionLn2()
	var maxDistance uint64
	var earlyUnderflows int
	var cases []oracleRow
	check.CompareSamples("godiff-math-exp2", edges, sample, func(word uint64) string {
		return fmt.Sprintf("%d %d", profile, word)
	}, func(word uint64) string {
		return strconv.FormatUint(binary64.Load(math.Exp2(binary64.Store(word))), 10)
	}, func(word uint64) string {
		x := binary64.Store(word)
		if x < -1074 && x > -1075 {
			if ideal := binary64.Load(idealExp2(x, ln2)); ideal != 1 {
				log.Fatalf("Exp2 early-underflow reference at %#016x: got %#016x, want smallest subnormal", word, ideal)
			}
			earlyUnderflows++
		}
		want, exact := reference(word, ln2)
		cases = append(cases, oracleRow{word, want, exact})
		if exact {
			return strconv.FormatUint(want, 10)
		}
		got := binary64.Load(math.Exp2(x))
		distance := max(got, want) - min(got, want)
		if distance > 1 {
			log.Fatalf("Exp2(%#016x) bits %#016x differ from reference %#016x by %d steps", word, got, want, distance)
		}
		maxDistance = max(maxDistance, distance)
		return strconv.FormatUint(got, 10)
	})
	fmt.Printf("Exp2 high-precision reference: maximum distance %d encoded steps in the Go evaluation interval; %d inputs in the early-underflow interval.\n",
		maxDistance, earlyUnderflows)
	if profile == 2 {
		profile = 3 // Generic arm64 contracts the polynomial but uses Ldexp scaling.
	}
	return compareGeneric(profile, cases)
}

func reference(word uint64, ln2 *big.Float) (uint64, bool) {
	x := binary64.Store(word)
	switch {
	case math.IsNaN(x):
		return word, true
	case x < -1074:
		return 0, true
	case x >= 1024:
		return binary64.Load(math.Inf(1)), true
	case x == math.Trunc(x):
		k := int(x)
		if k < -1022 {
			return 1 << (k + 1074), true
		}
		return uint64(k+1023) << 52, true
	default:
		return binary64.Load(idealExp2(x, ln2)), false
	}
}

func compareGeneric(profile int, cases []oracleRow) error {
	var modelInput, fixture bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&modelInput, profile, c.word)
	}
	for i, model := range check.LeanRows("godiff-math-exp2", len(cases), &modelInput) {
		c := cases[i]
		fmt.Fprintln(&fixture, c.word, model, c.reference, c.exact)
	}
	dir, err := os.MkdirTemp("", "golean-math-exp2-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_EXP2_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/exp2/testdata/"
	return check.TestInPackage("math", "TestGoLeanExp2", fixtures+"bridge_test.go", fixtures+"exp2_test.go")

}

func sample(r *rand.Rand, i int) uint64 {
	switch i % 5 {
	case 0:
		return r.Uint64()
	case 1:
		return binary64.Load(r.Float64()*2098 - 1074)
	case 2:
		k := float64(r.Intn(2099) - 1074)
		return binary64.Load(k+0.5) + uint64(r.Intn(9)-4)
	case 3:
		return binary64.Load(-1075 + 54*r.Float64())
	default:
		return binary64.Load(1022 + 3*r.Float64())
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

func idealExp2(x float64, ln2 *big.Float) float64 {
	value := func(x float64) *big.Float { return new(big.Float).SetPrec(precision).SetFloat64(x) }
	k := math.Trunc(x)
	r := value(x).Sub(value(x), value(k))
	r.Mul(r, ln2)
	term, sum, divisor := value(1), value(1), value(0)
	for n := int64(1); n <= 64; n++ {
		term.Mul(term, r).Quo(term, divisor.SetInt64(n))
		sum.Add(sum, term)
	}
	result, _ := sum.SetMantExp(sum, int(k)).Float64()
	return result
}
