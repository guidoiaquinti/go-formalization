package main

import (
	"bufio"
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

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	hi, lo float64
	k      int
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	profile := nativeProfile()
	var edges []input
	negativeZero := math.Copysign(0, -1)
	for _, k := range []int{-1075, -1074, -1073, -1023, -1022, -1, 0, 1, 1023, 1024} {
		for _, pair := range [][2]float64{{0, 0}, {negativeZero, 0}, {0, negativeZero}, {negativeZero, negativeZero},
			{math.Ln2 / 2, 0x1p-40}, {-math.Ln2 / 2, -0x1p-40}, {0x1p-53, 0}, {-0x1p-53, 0},
			{0x1p-10, -0x1p-27}, {-0x1p-10, 0x1p-27}} {
			edges = append(edges, input{pair[0], pair[1], k})
		}
	}
	dir, err := os.MkdirTemp("", "golean-math-expmulti-")
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
	compare := func(label string, cases []input) error {
		var model bytes.Buffer
		for _, x := range cases {
			if !(math.Abs(x.hi) <= 0.347 && math.Abs(x.lo) <= 2.1e-7 && math.Abs(x.hi-x.lo) <= 0.347 &&
				x.k >= -1075 && x.k <= 1024) {
				return fmt.Errorf("expmulti fixture outside reduced argument scope: %+v", x)
			}
			fmt.Fprintln(&model, profile, binary64.Load(x.hi), binary64.Load(x.lo), x.k)
		}
		for i, word := range check.LeanRows("godiff-math-expmulti", len(cases), &model) {
			x := cases[i]
			if _, err := fmt.Fprintln(fixture, binary64.Load(x.hi), binary64.Load(x.lo), x.k,
				word, binary64.Load(idealExpmulti(x))); err != nil {
				return err
			}
		}
		fmt.Printf("expmulti %s: %d Lean rows.\n", label, len(cases))
		return nil
	}
	if err := compare("boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20_000)
		for i := range cases {
			cases[i] = sample(rng, i)
		}
		if err := compare(fmt.Sprintf("seed=%d", seed), cases); err != nil {
			return err
		}
	}
	if err := fixture.Flush(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_EXPMULTI_CASES", path); err != nil {
		return err
	}
	const data = "difftest/math/expmulti/testdata/"
	return check.TestInPackage("math", "TestGoLeanExpmulti", data+"bridge_test.go", data+"expmulti_test.go")
}

func sample(rng *rand.Rand, i int) input {
	const ln2Hi = 6.93147180369123816490e-01
	const ln2Lo = 1.90821492927058770002e-10
	k := rng.Intn(2100) - 1075
	switch i % 5 {
	case 0:
		// Use the generic Exp reduction for the ordinary finite range.
		x := 1455*rng.Float64() - 745
		k = int(math.Round(x * math.Log2E))
		return input{x - float64(k)*ln2Hi, float64(k) * ln2Lo, k}
	case 1:
		// Use the generic Exp2 reduction for the ordinary finite range.
		x := 2098*rng.Float64() - 1074
		k = int(math.Round(x))
		t := x - float64(k)
		return input{t * ln2Hi, -t * ln2Lo, k}
	case 2:
		lo := 4.2e-7*rng.Float64() - 2.1e-7
		return input{(rng.Float64()-0.5)*math.Ln2 + lo, lo, k}
	case 3:
		edge := []int{-1075, -1074, -1023, -1022, 1023, 1024}[rng.Intn(6)]
		return input{(rng.Float64() - 0.5) * math.Ln2, 0, edge}
	default:
		lo := math.Ldexp(rng.Float64()-0.5, -rng.Intn(1060)-23)
		hi := math.Ldexp((rng.Float64()-0.5)*math.Ln2, -rng.Intn(1060))
		return input{hi, lo, k}
	}
}

func precise(value float64) *big.Float {
	return new(big.Float).SetPrec(256).SetFloat64(value)
}

func idealExpmulti(x input) float64 {
	reduced := precise(0).Sub(precise(x.hi), precise(x.lo))
	term, sum, divisor := precise(1), precise(1), precise(0)
	for n := int64(1); n <= 64; n++ {
		term.Mul(term, reduced)
		term.Quo(term, divisor.SetInt64(n))
		sum.Add(sum, term)
	}
	result, _ := sum.SetMantExp(sum, x.k).Float64()
	return result
}

func nativeProfile() int {
	switch runtime.GOARCH {
	case "arm64":
		return 2
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			log.Fatal("cannot determine expmulti build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		log.Fatalf("unmodeled expmulti architecture %s", runtime.GOARCH)
		return 0
	}
}
