package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"

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
	for _, sign := range []uint64{0, 1 << 63} {
		for bit := range 52 {
			top := uint64(1) << bit
			cases = append(cases, sign|(top-1), sign|top, sign|(top+1))
		}
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			word := rng.Uint64()
			if i%2 == 0 {
				word &= (1 << 63) | (1<<52 - 1)
			}
			cases = append(cases, word)
		}
	}

	var leanInput, fixtureInput bytes.Buffer
	want := make([]string, len(cases))
	for i, word := range cases {
		sign, exponent, significand := expected(word)
		want[i] = fmt.Sprintf("%d %d %d", sign, exponent, significand)
		fmt.Fprintln(&leanInput, word)
		fmt.Fprintln(&fixtureInput, word, sign, exponent, significand)
	}
	for i, row := range check.LeanRows("godiff-math-fma-split", len(cases), &leanInput) {
		if row != want[i] {
			return fmt.Errorf("word %#016x: Lean=%q want=%q", cases[i], row, want[i])
		}
	}
	dir, err := os.MkdirTemp("", "golean-math-fma-split-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	caseFile := filepath.Join(dir, "cases")
	if err := os.WriteFile(caseFile, fixtureInput.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_FMA_SPLIT_CASES", caseFile); err != nil {
		return err
	}
	const fixtures = "difftest/math/split/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanFmaSplit", fixtures+"bridge_test.go", fixtures+"split_test.go"); err != nil {
		return err
	}
	fmt.Printf("math.split: %d boundaries and %d samples agree, seeds %v.\n",
		boundaries, len(cases)-boundaries, seeds)
	return nil
}

func expected(word uint64) (uint32, int32, uint64) {
	const signWeight, fractionWeight = 1 << 63, 1 << 52
	sign := uint32(word / signWeight)
	magnitude := word % signWeight
	exponent := int32(magnitude / fractionWeight)
	fraction := magnitude % fractionWeight
	if exponent != 0 {
		return sign, exponent, fractionWeight + fraction
	}
	if fraction == 0 {
		return sign, -52, 0
	}
	exponent = 1
	for fraction < fractionWeight {
		exponent--
		fraction *= 2
	}
	return sign, exponent, fraction
}
