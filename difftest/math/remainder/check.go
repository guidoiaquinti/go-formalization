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
	for _, value := range []float64{
		0.5, 1.5, 2.5, 3.5, 2, 3, 5, 7,
		0x1p-1021, math.Nextafter(0x1p-1021, 0), math.Nextafter(0x1p-1021, 1),
		math.MaxFloat64 / 2, math.Nextafter(math.MaxFloat64/2, math.MaxFloat64),
	} {
		words = append(words, math.Float64bits(value), math.Float64bits(-value))
	}
	words = append(words, 2, 3, 5)
	var edges []input
	for _, x := range words {
		for _, y := range words {
			edges = append(edges, input{x, y})
		}
	}
	var fixture bytes.Buffer
	check.CompareSamples("godiff-math-remainder", edges, sample, func(c input) string {
		return fmt.Sprintf("%d %d", c.x, c.y)
	}, func(c input) string {
		return observe(math.Remainder(binary64.Store(c.x), binary64.Store(c.y)))
	}, func(c input) string {
		want := expected(c)
		fmt.Fprintln(&fixture, c.x, c.y, want)
		return want
	})

	dir, err := os.MkdirTemp("", "golean-math-remainder-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_REMAINDER_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/remainder/testdata/"
	return check.TestInPackage("math", "TestGoLeanRemainder", fixtures+"bridge_test.go", fixtures+"remainder_test.go")
}

func sample(r *rand.Rand, i int) input {
	x, y := r.Uint64(), r.Uint64()
	switch i % 5 {
	case 0:
		y &= 1<<52 - 1
	case 1:
		y = y&^(uint64(2047)<<52) | x&(uint64(2047)<<52)
	case 2:
		y = x
	case 3:
		// Exact half-integer quotients, with both signs and varied exponents.
		exponent := r.Intn(1800) - 900
		x = math.Float64bits(math.Ldexp(float64(r.Intn(201)-100)+0.5, exponent))
		y = math.Float64bits(math.Ldexp(1, exponent))
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
	var whole, fraction big.Int
	whole.QuoRem(quotient.Num(), quotient.Denom(), &fraction)
	twice := new(big.Int).Lsh(new(big.Int).Abs(&fraction), 1)
	if order := twice.Cmp(quotient.Denom()); order > 0 || order == 0 && whole.Bit(0) == 1 {
		whole.Add(&whole, big.NewInt(int64(quotient.Sign())))
	}
	residue := n.Sub(n, new(big.Rat).Mul(new(big.Rat).SetInt(&whole), d))
	if residue.Sign() == 0 {
		return observe(math.Copysign(0, x))
	}
	value, exact := residue.Float64()
	if !exact {
		log.Fatalf("nonrepresentable remainder for %#x and %#x", c.x, c.y)
	}
	return observe(value)
}

func observe(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(binary64.Load(x), 10)
}
