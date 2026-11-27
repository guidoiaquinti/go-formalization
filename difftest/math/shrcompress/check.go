package main

import (
	"bytes"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	high, low uint64
	shift     uint
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var cases []input
	const maxWord = ^uint64(0)
	for _, shift := range []uint{0, 1, 31, 32, 62, 63, 64, 65, 66, 95, 96, 126, 127, 128, 129, ^uint(0)} {
		for _, high := range []uint64{0, 1, 1 << 63, 0xabcdef1234567890, maxWord} {
			for _, low := range []uint64{0, 1, 1 << 63, 0x123456abcdef9870, maxWord} {
				cases = append(cases, input{high, low, shift})
			}
		}
	}
	for _, shift := range []uint{1, 32, 63, 64, 65, 96, 127, 128} {
		for _, bit := range []uint{shift - 1, shift, shift + 1} {
			if bit < 64 {
				cases = append(cases, input{0, uint64(1) << bit, shift})
			} else if bit < 128 {
				cases = append(cases, input{uint64(1) << (bit - 64), 0, shift})
			}
		}
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			shift := uint(rng.Uint64())
			if i%4 != 0 {
				shift = uint(rng.Intn(132))
			}
			cases = append(cases, input{rng.Uint64(), rng.Uint64(), shift})
		}
	}

	var leanInput, fixtureInput bytes.Buffer
	want := make([]string, len(cases))
	for i, c := range cases {
		high, low := expected(c)
		want[i] = fmt.Sprintf("%d %d", high, low)
		fmt.Fprintln(&leanInput, c.high, c.low, c.shift)
		fmt.Fprintln(&fixtureInput, c.high, c.low, c.shift, high, low)
	}
	for i, row := range check.LeanRows("godiff-math-shrcompress", len(cases), &leanInput) {
		if row != want[i] {
			return fmt.Errorf("input %+v: Lean=%q want=%q", cases[i], row, want[i])
		}
	}
	dir, err := os.MkdirTemp("", "golean-math-shrcompress-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	caseFile := filepath.Join(dir, "cases")
	if err := os.WriteFile(caseFile, fixtureInput.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_SHRCOMPRESS_CASES", caseFile); err != nil {
		return err
	}
	const fixtures = "difftest/math/shrcompress/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanShrCompress", fixtures+"bridge_test.go", fixtures+"shrcompress_test.go"); err != nil {
		return err
	}
	fmt.Printf("math.shrcompress: %d boundaries and %d samples agree, seeds %v.\n",
		boundaries, len(cases)-boundaries, seeds)
	return nil
}

func expected(c input) (uint64, uint64) {
	shift := min(c.shift, 128)
	var combined, low, shifted, restored big.Int
	combined.SetUint64(c.high)
	combined.Lsh(&combined, 64)
	low.SetUint64(c.low)
	combined.Add(&combined, &low)
	shifted.Rsh(&combined, shift)
	restored.Lsh(&shifted, shift)
	if restored.Cmp(&combined) != 0 {
		shifted.SetBit(&shifted, 0, 1)
	}
	low.SetUint64(^uint64(0))
	low.And(&shifted, &low)
	shifted.Rsh(&shifted, 64)
	return shifted.Uint64(), low.Uint64()
}
