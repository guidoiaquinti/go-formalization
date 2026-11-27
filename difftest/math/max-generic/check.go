package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type pair struct{ a, b uint64 }

const (
	sign = uint64(1 << 63)
	inf  = uint64(0x7ff0000000000000)
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	words := binary64.Boundaries()
	var cases []pair
	for _, a := range words {
		for _, b := range words {
			cases = append(cases, pair{a, b})
		}
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			cases = append(cases, sample(rng, i))
		}
	}

	var modelInput, fixtureInput bytes.Buffer
	for _, p := range cases {
		fmt.Fprintln(&modelInput, p.a, p.b)
		fmt.Fprintln(&fixtureInput, p.a, p.b, expected(p))
	}
	for i, row := range check.LeanRows("godiff-math-max-generic", len(cases), &modelInput) {
		if want := expected(cases[i]); row != want {
			return fmt.Errorf("input %#v: Lean=%s want=%s", cases[i], row, want)
		}
	}
	dir, err := os.MkdirTemp("", "golean-math-max-generic-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	caseFile := filepath.Join(dir, "cases")
	if err := os.WriteFile(caseFile, fixtureInput.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_MAX_GENERIC_CASES", caseFile); err != nil {
		return err
	}
	if err := check.TestInPackage("math", "TestGoLeanMaxGeneric", "difftest/math/max-generic/testdata/bridge_test.go", "difftest/math/max-generic/testdata/max_test.go"); err != nil {
		return err
	}
	fmt.Printf("math.max: %d boundaries and %d samples agree exactly, seeds %v.\n", boundaries, len(cases)-boundaries, seeds)
	return nil
}

func sample(rng *rand.Rand, i int) pair {
	a, b := rng.Uint64(), rng.Uint64()
	switch i % 5 {
	case 1:
		b = a
	case 2:
		b = a ^ sign
	case 3:
		a = inf
	case 4:
		b |= inf | 1 // Force a NaN, in either sign.
	}
	return pair{a, b}
}

// Key orders finite values and infinities, with negative zero first among zeros.
func key(word uint64) uint64 {
	if word&sign != 0 {
		return ^word
	}
	return word | sign
}

func expected(p pair) string {
	switch {
	case p.a == inf || p.b == inf:
		return strconv.FormatUint(inf, 10)
	case p.a%sign > inf || p.b%sign > inf:
		return strconv.FormatUint(0x7ff8000000000001, 10)
	case key(p.a) >= key(p.b):
		return strconv.FormatUint(p.a, 10)
	default:
		return strconv.FormatUint(p.b, 10)
	}
}
