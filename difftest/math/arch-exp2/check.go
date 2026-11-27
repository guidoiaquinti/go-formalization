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

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	if runtime.GOARCH == "loong64" || runtime.GOARCH == "riscv64" {
		fmt.Printf("The modeled archExp2 profiles do not cover %s assembly.\n", runtime.GOARCH)
		return
	}
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	profile := 0
	if runtime.GOARCH == "arm64" {
		profile = 2
	}
	edges := binary64.Boundaries()
	for _, center := range []float64{-1075, -1074.75, -1074.5, -1074.25, -1074, -1073.5,
		-1073, -1022, -1, -0.5, 0, 0.5, 1, 1023, 1023.5, 1024} {
		edges = append(edges, binary64.Load(math.Nextafter(center, math.Inf(-1))),
			binary64.Load(center), binary64.Load(math.Nextafter(center, math.Inf(1))))
	}
	for k := -1074; k <= 1023; k++ {
		edges = append(edges, binary64.Load(float64(k)))
	}
	var stubInput bytes.Buffer
	specials := binary64.Boundaries()
	for _, word := range specials {
		fmt.Fprintln(&stubInput, 0, word)
	}
	for i, model := range check.LeanRows("godiff-math-arch-exp2", len(specials), &stubInput) {
		if model != "panic:not implemented" {
			return fmt.Errorf("archExp2 stub(%#016x) returned %s", specials[i], model)
		}
	}

	dir, err := os.MkdirTemp("", "golean-math-arch-exp2-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	var fixture bytes.Buffer
	ln2 := highPrecisionLn2()
	compare := func(label string, cases []uint64) {
		var input bytes.Buffer
		for _, word := range cases {
			fmt.Fprintln(&input, profile, word)
		}
		for i, model := range check.LeanRows("godiff-math-arch-exp2", len(cases), &input) {
			word := cases[i]
			want, exact := uint64(0), false
			if profile == 2 {
				want, exact = reference(word, ln2)
				x := binary64.Store(word)
				if -1075 < x && x < -1074 && binary64.Load(idealExp2(x, ln2)) != 1 {
					log.Fatalf("independent result for early underflow at %#016x must round to the smallest subnormal", word)
				}
			}
			fmt.Fprintf(&fixture, "%d %q %d %t\n", word, model, want, exact)
		}
		fmt.Printf("archExp2 %s: %d Lean rows.\n", label, len(cases))
	}
	compare("boundaries", edges)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]uint64, 20_000)
		for i := range cases {
			cases[i] = sample(rng, i)
		}
		compare(fmt.Sprintf("seed=%d", seed), cases)
	}

	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_ARCH_EXP2_CASES", path); err != nil {
		return err
	}
	const data = "difftest/math/arch-exp2/testdata/"
	return check.TestInPackage("math", "TestGoLeanArchExp2", data+"bridge_test.go", data+"arch_exp2_test.go")
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
