package main

import (
	"bytes"
	"fmt"
	"log"
	"math/bits"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const (
	fractionWeight = 1 << 52
	signWeight     = 1 << 63
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := binary64.Boundaries()
	for k := range 52 {
		power := uint64(1) << k
		for _, sign := range []uint64{0, signWeight} {
			cases = append(cases, sign+power-1, sign+power, sign+power+1)
		}
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			word := rng.Uint64()
			if i%2 == 0 {
				word &= signWeight | (fractionWeight - 1)
			}
			cases = append(cases, word)
		}
	}
	var modelInput, fixtureInput bytes.Buffer
	for _, word := range cases {
		fmt.Fprintln(&modelInput, word)
		fmt.Fprintln(&fixtureInput, word, expected(word))
	}
	for i, row := range check.LeanRows("godiff-math-frexp-generic", len(cases), &modelInput) {
		if want := expected(cases[i]); row != want {
			return fmt.Errorf("word %#016x: Lean=%q want=%q", cases[i], row, want)
		}
	}
	dir, err := os.MkdirTemp("", "golean-math-frexp-generic-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	caseFile := filepath.Join(dir, "cases")
	if err := os.WriteFile(caseFile, fixtureInput.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_FREXP_GENERIC_CASES", caseFile); err != nil {
		return err
	}
	const fixtures = "difftest/math/frexp-generic/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanFrexpGeneric", fixtures+"bridge_test.go", fixtures+"frexp_test.go"); err != nil {
		return err
	}
	fmt.Printf("math.frexp: %d boundaries and %d samples agree exactly, seeds %v.\n", boundaries, len(cases)-boundaries, seeds)
	return nil
}

func expected(word uint64) string {
	magnitude := word % signWeight
	sign := word - magnitude
	e, f := magnitude/fractionWeight, magnitude%fractionWeight
	switch {
	case magnitude == 0 || e == 2047:
		return result(word, 0)
	case e == 0:
		top := bits.Len64(f) - 1
		payload := (f - (uint64(1) << top)) << (52 - top)
		return result(sign+1022*fractionWeight+payload, top-1073)
	default:
		return result(sign+1022*fractionWeight+f, int(e)-1022)
	}
}

func result(fraction uint64, exponent int) string {
	return fmt.Sprintf("%d %d", fraction, exponent)
}
