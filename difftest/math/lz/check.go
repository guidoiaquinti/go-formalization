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

type words struct{ high, low uint64 }

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []words{{}}
	for bit := range 64 {
		top := uint64(1) << bit
		cases = append(cases, words{top, 0}, words{0, top},
			words{top | (top - 1), ^uint64(0)}, words{0, top | (top - 1)})
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			high, low := rng.Uint64(), rng.Uint64()
			if i%2 == 0 {
				cases = append(cases, words{high >> rng.Intn(64), low})
			} else {
				cases = append(cases, words{0, low >> rng.Intn(64)})
			}
		}
	}

	var leanInput, fixtureInput bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&leanInput, c.high, c.low)
		fmt.Fprintln(&fixtureInput, c.high, c.low, expected(c))
	}
	for i, row := range check.LeanRows("godiff-math-lz", len(cases), &leanInput) {
		if want := strconv.Itoa(expected(cases[i])); row != want {
			return fmt.Errorf("words (%#016x, %#016x): Lean=%q want=%q", cases[i].high, cases[i].low, row, want)
		}
	}
	dir, err := os.MkdirTemp("", "golean-math-lz-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	caseFile := filepath.Join(dir, "cases")
	if err := os.WriteFile(caseFile, fixtureInput.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_LZ_CASES", caseFile); err != nil {
		return err
	}
	const fixtures = "difftest/math/lz/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanLz", fixtures+"bridge_test.go", fixtures+"lz_test.go"); err != nil {
		return err
	}
	fmt.Printf("math.lz: %d boundaries and %d samples agree, seeds %v.\n",
		boundaries, len(cases)-boundaries, seeds)
	return nil
}

// Walk the joined 128-bit value from its highest bit until a one is found.
func expected(c words) int {
	for index := range 128 {
		value, bit := c.high, 63-index
		if index >= 64 {
			value, bit = c.low, 127-index
		}
		if value>>bit&1 == 1 {
			return index
		}
	}
	return 128
}
