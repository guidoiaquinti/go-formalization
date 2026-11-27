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

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	words := binary64.Boundaries()
	for _, sign := range []uint64{0, 1 << 63} {
		for _, exponent := range []uint64{1022, 1023, 1024, 1073, 1074, 1075, 1076} {
			for _, fraction := range []uint64{0, 1, (1 << 51) - 1, 1 << 51, (1 << 52) - 1} {
				words = append(words, sign|exponent<<52|fraction)
			}
		}
	}
	boundaries := len(words)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			word := rng.Uint64()
			if i%2 == 0 {
				word = word&^uint64(0x7ff<<52) | uint64(1022+rng.Intn(55))<<52
			}
			words = append(words, word)
		}
	}
	var leanInput, fixture bytes.Buffer
	want := make([]string, 0, len(words))
	for _, word := range words {
		x := binary64.Store(word)
		row := oracle(x)
		if actual := observe(math.Floor(x)) + " " + observe(math.Ceil(x)); actual != row {
			return fmt.Errorf("input %#016x: public Go=%q want=%q", word, actual, row)
		}
		fmt.Fprintln(&leanInput, word)
		fmt.Fprintln(&fixture, word, row)
		want = append(want, row)
	}
	for i, row := range check.LeanRows("godiff-math-floor-ceil", len(words), &leanInput) {
		if row != want[i] {
			return fmt.Errorf("input %#016x: Lean=%q want=%q", words[i], row, want[i])
		}
	}
	dir, err := os.MkdirTemp("", "golean-floor-ceil-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_FLOOR_CEIL_CASES", path); err != nil {
		return err
	}
	const prefix = "difftest/math/floor-ceil/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanFloorCeil", prefix+"bridge_test.go", prefix+"floor_ceil_test.go"); err != nil {
		return err
	}
	fmt.Printf("math.Floor and math.Ceil: %d boundaries and %d samples agree, seeds %v.\n",
		boundaries, len(words)-boundaries, seeds)
	return nil
}

func oracle(x float64) string {
	if x == 0 || math.IsInf(x, 0) || math.IsNaN(x) {
		return observe(x) + " " + observe(x)
	}
	r := new(big.Rat).SetFloat64(x)
	var whole, remainder big.Int
	whole.QuoRem(r.Num(), r.Denom(), &remainder)
	var low, high big.Int
	low.Set(&whole)
	high.Set(&whole)
	if remainder.Sign() < 0 {
		low.Sub(&low, big.NewInt(1))
	} else if remainder.Sign() > 0 {
		high.Add(&high, big.NewInt(1))
	}
	return integerWord(&low, x) + " " + integerWord(&high, x)
}

func integerWord(n *big.Int, input float64) string {
	if n.Sign() == 0 {
		return observe(math.Copysign(0, input))
	}
	x, _ := n.Float64()
	return observe(x)
}

func observe(x float64) string {
	if math.IsNaN(x) {
		return "NaN"
	}
	return strconv.FormatUint(binary64.Load(x), 10)
}
