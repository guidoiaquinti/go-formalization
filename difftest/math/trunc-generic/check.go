package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := binary64.Boundaries()
	for exponent := uint64(1022); exponent <= 1075; exponent++ {
		quantum := uint64(1 << 52)
		if exponent >= 1023 {
			quantum >>= exponent - 1023
		}
		for _, sign := range []uint64{0, 1 << 63} {
			for _, payload := range []uint64{0, 1, quantum / 2, quantum - 1, 1<<51 + 1, 1<<52 - 1} {
				cases = append(cases, sign|exponent<<52|payload)
			}
		}
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			word := rng.Uint64()
			if i%2 == 0 {
				word = word&^uint64(2047<<52) | uint64(1022+rng.Intn(54))<<52
			}
			cases = append(cases, word)
		}
	}
	var leanInput, fixtureInput bytes.Buffer
	want := make([]uint64, len(cases))
	for i, word := range cases {
		want[i] = expectedWord(word)
		if numeric := numericWord(word); numeric != want[i] {
			return fmt.Errorf("word %#016x: numeric oracle=%#016x, word oracle=%#016x", word, numeric, want[i])
		}
		fmt.Fprintln(&leanInput, word)
		fmt.Fprintln(&fixtureInput, word, want[i])
	}
	for i, row := range check.LeanRows("godiff-math-trunc-generic", len(cases), &leanInput) {
		if expected := strconv.FormatUint(want[i], 10); row != expected {
			return fmt.Errorf("word %#016x: Lean=%q want=%q", cases[i], row, expected)
		}
	}
	dir, err := os.MkdirTemp("", "golean-math-trunc-generic-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	caseFile := filepath.Join(dir, "cases")
	if err := os.WriteFile(caseFile, fixtureInput.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_TRUNC_GENERIC_CASES", caseFile); err != nil {
		return err
	}
	const fixtures = "difftest/math/trunc-generic/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanTruncGeneric", fixtures+"bridge_test.go", fixtures+"trunc_generic_test.go"); err != nil {
		return err
	}
	fmt.Printf("math.trunc: %d boundaries and %d samples agree, seeds %v.\n",
		boundaries, len(cases)-boundaries, seeds)
	return nil
}

// Evaluate the kernel theorem's quotient formula using unsigned words.
func expectedWord(word uint64) uint64 {
	switch exponent := word >> 52 & 2047; {
	case exponent < 1023:
		return word - word%(1<<63)
	case exponent < 1075:
		quantum := uint64(1) << (1075 - exponent)
		return word / quantum * quantum
	default:
		return word
	}
}

// Converting finite values below 2^52 to int64 truncates their fractions
// without losing integer bits. Larger finite values are already integral.
func numericWord(word uint64) uint64 {
	x := binary64.Store(word)
	if math.IsNaN(x) || x <= -(1<<52) || x >= 1<<52 {
		return word
	}
	integer := int64(x)
	if integer == 0 {
		return word & (1 << 63)
	}
	return binary64.Load(float64(integer))
}
