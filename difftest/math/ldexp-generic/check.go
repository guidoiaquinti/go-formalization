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
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const (
	fractionWeight = 1 << 52
	signWeight     = 1 << 63
)

type input struct {
	word uint64
	exp  int
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	words := append(binary64.Boundaries(),
		0x3fe0000000000000, 0xbfe0000000000000,
		0x4000000000000000, 0xc000000000000000,
		0x3ff0000000000001, 0x3ff0000000000002,
		0x3ff0000000000003, 0x3fffffffffffffff)
	var boundaries []input
	for _, word := range words {
		_, _, highest := fields(word)
		for _, exp := range []int{math.MinInt, math.MinInt + 1, math.MinInt + 1074,
			-1076 - highest, -1075 - highest, -1074 - highest,
			-1024 - highest, -1023 - highest, -1022 - highest, -1021 - highest,
			1022 - highest, 1023 - highest, 1024 - highest,
			-1, 0, 1, math.MaxInt - 1023, math.MaxInt - 1, math.MaxInt} {
			boundaries = append(boundaries, input{word, exp})
		}
	}
	cases := boundaries
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			cases = append(cases, sample(rng, i))
		}
	}
	var modelInput, fixtureInput bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&modelInput, strconv.IntSize, c.word, c.exp)
		fmt.Fprintln(&fixtureInput, c.word, c.exp, expected(c))
	}
	for i, row := range check.LeanRows("godiff-math-ldexp-generic", len(cases), &modelInput) {
		c := cases[i]
		if want := strconv.FormatUint(expected(c), 10); row != want {
			return fmt.Errorf("word %#016x exp %d: Lean=%s want=%s", c.word, c.exp, row, want)
		}
	}
	dir, err := os.MkdirTemp("", "golean-math-ldexp-generic-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	caseFile := filepath.Join(dir, "cases")
	if err := os.WriteFile(caseFile, fixtureInput.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_LDEXP_GENERIC_CASES", caseFile); err != nil {
		return err
	}
	const fixtures = "difftest/math/ldexp-generic/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanLdexpGeneric", fixtures+"bridge_test.go", fixtures+"ldexp_test.go"); err != nil {
		return err
	}
	fmt.Printf("math.ldexp: %d boundaries and %d samples agree exactly, seeds %v.\n", len(boundaries), len(cases)-len(boundaries), seeds)
	return nil
}

func sample(rng *rand.Rand, index int) input {
	word := rng.Uint64()
	if index%2 == 0 {
		word &= signWeight | (fractionWeight - 1)
	}
	_, _, highest := fields(word)
	var exp int
	switch index % 5 {
	case 0:
		exp = int(rng.Uint64())
	case 1:
		exp = math.MinInt + rng.Intn(1100)
	case 2:
		exp = math.MaxInt - rng.Intn(1100)
	case 3:
		exp = rng.Intn(65) - 1080 - highest
	default:
		exp = rng.Intn(2200) - 1100 - highest
	}
	return input{word, exp}
}

// fields also gives the unbiased exponent of the highest set value bit.
func fields(word uint64) (exponent, payload uint64, highest int) {
	exponent, payload = (word>>52)&2047, word%fractionWeight
	if exponent == 0 {
		return exponent, payload, bits.Len64(payload) - 1 - 1074
	}
	return exponent, payload, int(exponent) - 1023
}

// Compute the exact dyadic product at the wrapped exponent, then rely on
// big.Float's independent rounding implementation to convert to binary64.
func expected(c input) uint64 {
	e, payload, highest := fields(c.word)
	if e == 0 && payload == 0 || e == 2047 {
		return c.word
	}
	sign := c.word & signWeight
	wrapped := c.exp + highest
	if wrapped < -1075 {
		return sign
	}
	if wrapped > 1023 {
		return sign | 0x7ff0000000000000
	}
	coefficient := fractionWeight + payload
	if e == 0 {
		coefficient = payload << (53 - bits.Len64(payload))
	}
	f := new(big.Float).SetUint64(coefficient)
	f.SetMantExp(f, wrapped-52)
	rounded, _ := f.Float64()
	return binary64.Load(rounded) | sign
}
