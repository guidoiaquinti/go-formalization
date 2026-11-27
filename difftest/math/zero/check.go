package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []uint64{0, ^uint64(0)}
	for bit := range 64 {
		power := uint64(1) << bit
		cases = append(cases, power-1, power, power+1, ^power)
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			word := rng.Uint64()
			switch i % 4 {
			case 0:
				word = 0
			case 1:
				word &= uint64(1) << rng.Intn(64)
			case 2:
				word = word >> rng.Intn(64)
			}
			cases = append(cases, word)
		}
	}
	var leanInput, fixtureInput bytes.Buffer
	for _, word := range cases {
		fmt.Fprintln(&leanInput, word)
		fmt.Fprintln(&fixtureInput, word, expected(word))
	}
	for i, row := range check.LeanRows("godiff-math-fma-zero", len(cases), &leanInput) {
		if want := strconv.FormatUint(expected(cases[i]), 10); row != want {
			return fmt.Errorf("word %#016x: Lean=%q want=%q", cases[i], row, want)
		}
	}
	dir, err := os.MkdirTemp("", "golean-math-fma-zero-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	caseFile := filepath.Join(dir, "cases")
	if err := os.WriteFile(caseFile, fixtureInput.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_FMA_ZERO_CASES", caseFile); err != nil {
		return err
	}
	const fixtures = "difftest/math/zero/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanFmaZero", fixtures+"bridge_test.go", fixtures+"zero_test.go"); err != nil {
		return err
	}
	fmt.Printf("math.zero: %d boundaries and %d samples agree, seeds %v.\n",
		boundaries, len(cases)-boundaries, seeds)
	return nil
}

func expected(word uint64) uint64 {
	// The unsigned subtraction borrows into bit 63 only if the folded
	// low 63 bits were all zero. Folding retains every input bit.
	return (((word >> 1) | (word & 1)) - 1) >> 63
}
