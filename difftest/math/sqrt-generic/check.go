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
		for exponent := uint64(1); exponent < 2047; exponent++ {
			power := exponent << 52
			cases = append(cases, sign|(power-1), sign|power, sign|(power+1))
		}
		for bit := range 52 {
			power := uint64(1) << bit
			cases = append(cases, sign|(power-1), sign|power, sign|(power+1))
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
				word &= (1 << 63) - 1
			case 1:
				word &= (1 << 52) - 1
			case 2:
				value := float64(uint32(word)) + 1
				word = binary64.Load(value*value) + uint64(i%3) - 1
			}
			cases = append(cases, word)
		}
	}

	var leanInput, fixtureInput bytes.Buffer
	want := make([]string, len(cases))
	for i, word := range cases {
		want[i] = expected(word)
		fmt.Fprintln(&leanInput, word)
		fmt.Fprintln(&fixtureInput, word, want[i])
	}
	for i, row := range check.LeanRows("godiff-math-sqrt-generic", len(cases), &leanInput) {
		if row != want[i] {
			return fmt.Errorf("word %#016x: Lean=%q want=%q", cases[i], row, want[i])
		}
	}
	dir, err := os.MkdirTemp("", "golean-math-sqrt-generic-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	caseFile := filepath.Join(dir, "cases")
	if err := os.WriteFile(caseFile, fixtureInput.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_SQRT_GENERIC_CASES", caseFile); err != nil {
		return err
	}
	const fixtures = "difftest/math/sqrt-generic/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanSqrtGeneric", fixtures+"bridge_test.go", fixtures+"sqrt_generic_test.go"); err != nil {
		return err
	}
	fmt.Printf("math.sqrt: %d boundaries and %d samples agree, seeds %v.\n",
		boundaries, len(cases)-boundaries, seeds)
	return nil
}

func expected(word uint64) string {
	const mantissa, infinity = uint64(1 << 52), uint64(2047 << 52)
	if word == infinity || word<<1 == 0 {
		return strconv.FormatUint(word, 10)
	}
	if word >= infinity {
		return "NaN"
	}
	exponent := int(word/mantissa) - 1023
	coefficient := word % mantissa
	if exponent == -1023 {
		exponent++
		for coefficient < mantissa {
			coefficient *= 2
			exponent--
		}
	} else {
		coefficient += mantissa
	}
	if exponent%2 != 0 {
		coefficient *= 2
	}
	var scaled, root, square big.Int
	scaled.SetUint64(coefficient)
	scaled.Lsh(&scaled, 54)
	root.Sqrt(&scaled)
	square.Mul(&root, &root)
	word = root.Uint64() / 2
	if root.Bit(0) == 1 && square.Cmp(&scaled) < 0 {
		word++
	}
	word += uint64((exponent>>1)+1022) * mantissa
	return strconv.FormatUint(word, 10)
}
