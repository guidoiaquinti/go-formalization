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
	words := binary64.Boundaries()
	words = append(words, 0x3ff8000000000000, 0x4008000000000000, 0xc008000000000000,
		0x4000000000000000, 0xc000000000000000, 0x7fefffffffffffff)
	var edges []input
	for _, x := range words {
		for _, y := range words {
			edges = append(edges, input{x, y})
		}
	}
	dir, err := os.MkdirTemp("", "golean-math-mod-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.Setenv("GOLEAN_MATH_MOD_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/mod/testdata/"
	compare := func(label string, cases []input) error {
		var leanInput, fixture bytes.Buffer
		want := make([]string, len(cases))
		for i, c := range cases {
			want[i] = expected(c)
			if got := observe(math.Mod(binary64.Store(c.x), binary64.Store(c.y))); got != want[i] {
				return fmt.Errorf("%s input %#016x %% %#016x: public Go=%s want=%s", label, c.x, c.y, got, want[i])
			}
			fmt.Fprintln(&leanInput, c.x, c.y)
			fmt.Fprintln(&fixture, c.x, c.y, want[i])
		}
		for i, row := range check.LeanRows("godiff-math-mod", len(cases), &leanInput) {
			if row != want[i] {
				c := cases[i]
				return fmt.Errorf("%s input %#016x %% %#016x: Lean=%s want=%s", label, c.x, c.y, row, want[i])
			}
		}
		if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
			return err
		}
		if err := check.TestInPackage("math", "TestGoLeanMod", fixtures+"bridge_test.go", fixtures+"mod_test.go"); err != nil {
			return fmt.Errorf("%s: %w", label, err)
		}
		fmt.Printf("math.Mod and math.mod %s: %d inputs agree.\n", label, len(cases))
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
	return nil
}

func sample(r *rand.Rand, i int) input {
	x, y := r.Uint64(), r.Uint64()
	switch i % 4 {
	case 0:
		y &= 1<<52 - 1
	case 1:
		y = y&^(uint64(2047)<<52) | x&(uint64(2047)<<52)
	case 2:
		y = x
	}
	return input{x, y}
}

func expected(c input) string {
	x, y := binary64.Store(c.x), binary64.Store(c.y)
	if y == 0 || math.IsInf(x, 0) || math.IsNaN(x) || math.IsNaN(y) {
		return "NaN"
	}
	if math.IsInf(y, 0) || x == 0 {
		return strconv.FormatUint(c.x, 10)
	}
	n := new(big.Rat).SetFloat64(x)
	d := new(big.Rat).SetFloat64(y)
	quotient := new(big.Rat).Quo(n, d)
	whole := new(big.Int).Quo(quotient.Num(), quotient.Denom())
	residue := n.Sub(n, new(big.Rat).Mul(new(big.Rat).SetInt(whole), d))
	if residue.Sign() == 0 {
		return observe(math.Copysign(0, x))
	}
	value, exact := residue.Float64()
	if !exact {
		log.Fatalf("nonrepresentable remainder for %#x %% %#x", c.x, c.y)
	}
	return observe(value)
}

func observe(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(binary64.Load(x), 10)
}
