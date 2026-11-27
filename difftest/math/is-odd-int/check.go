package main

import (
	"bytes"
	"fmt"
	"log"
	"math/big"
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
	for _, sign := range []uint64{0, 1 << 63} {
		for _, exponent := range []uint64{1022, 1023, 1024, 1073, 1074, 1075, 1076, 1077} {
			for _, payload := range []uint64{0, 1, 2, 3, 1<<51 - 1, 1 << 51, 1<<51 + 1, 1<<52 - 1} {
				cases = append(cases, sign|exponent<<52|payload)
			}
		}
	}
	for power := 1; power <= 53; power++ {
		value := int64(1) << power
		for _, integer := range []int64{value - 1, value, value + 1} {
			cases = append(cases, binary64.Load(float64(integer)), binary64.Load(float64(-integer)))
		}
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			word := rng.Uint64()
			switch i % 4 {
			case 0:
				word = word&^uint64(2047<<52) | uint64(1022+rng.Intn(56))<<52
			case 1, 2:
				integer := int64((word & ((1 << 53) - 1)) >> rng.Intn(53))
				if word>>63 != 0 {
					integer = -integer
				}
				word = binary64.Load(float64(integer))
				if i%4 == 2 {
					word++
				}
			}
			cases = append(cases, word)
		}
	}
	var leanInput, fixtureInput bytes.Buffer
	for _, word := range cases {
		fmt.Fprintln(&leanInput, word)
		fmt.Fprintln(&fixtureInput, word, expected(word))
	}
	for i, row := range check.LeanRows("godiff-math-is-odd-int", len(cases), &leanInput) {
		if want := strconv.FormatUint(expected(cases[i]), 10); row != want {
			return fmt.Errorf("word %#016x: Lean=%q want=%q", cases[i], row, want)
		}
	}
	dir, err := os.MkdirTemp("", "golean-math-is-odd-int-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	caseFile := filepath.Join(dir, "cases")
	if err := os.WriteFile(caseFile, fixtureInput.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_IS_ODD_INT_CASES", caseFile); err != nil {
		return err
	}
	const fixtures = "difftest/math/is-odd-int/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanIsOddInt", fixtures+"bridge_test.go", fixtures+"is_odd_int_test.go"); err != nil {
		return err
	}
	fmt.Printf("math.isOddInt: %d boundaries and %d samples agree, seeds %v.\n",
		boundaries, len(cases)-boundaries, seeds)
	return nil
}

func expected(word uint64) uint64 {
	exact := new(big.Rat).SetFloat64(binary64.Store(word))
	if exact != nil && exact.IsInt() && exact.Num().Bit(0) == 1 {
		return 1
	}
	return 0
}
