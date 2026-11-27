package main

import (
	"bytes"
	"fmt"
	"log"
	"math/bits"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	const sign = uint64(1 << 63)
	var words []uint64
	for _, high := range []uint64{0, sign} {
		for power := range 52 {
			leading := uint64(1) << power
			words = append(words, high|leading, high|(leading+1))
		}
		words = append(words, high|((1<<52)-1), high|(1<<52), high|0x3ff0000000000000, high|0x7fefffffffffffff)
	}
	boundaries := len(words)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			word := rng.Uint64()
			if i%2 == 0 {
				word = (word & sign) | (1 + (word & ((1 << 52) - 1)))
			}
			if word>>52&2047 == 2047 || word&(sign-1) == 0 {
				word ^= 1 << 52
			}
			words = append(words, word)
		}
	}
	var leanInput, fixtureInput bytes.Buffer
	for _, word := range words {
		fmt.Fprintln(&leanInput, word)
		fmt.Fprintf(&fixtureInput, "%d %s\n", word, expected(word))
	}
	for i, row := range check.LeanRows("godiff-math-normalize", len(words), &leanInput) {
		if want := expected(words[i]); row != want {
			return fmt.Errorf("word %#016x: Lean=%q want=%q", words[i], row, want)
		}
	}
	dir, err := os.MkdirTemp("", "golean-normalize-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixtureInput.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_NORMALIZE_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/math/normalize/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanNormalize", fixtures+"bridge_test.go", fixtures+"normalize_test.go"); err != nil {
		return err
	}
	fmt.Printf("math.normalize and math.ilogb: %d boundaries and %d samples agree, seeds %v.\n",
		boundaries, len(words)-boundaries, seeds)
	return nil
}

func expected(word uint64) string {
	const sign = uint64(1 << 63)
	const mantissa = uint64(1 << 52)
	magnitude := word & (sign - 1)
	if e := magnitude / mantissa; e != 0 {
		return fmt.Sprintf("%d 0 %d", word, int(e)-1023)
	}
	k := bits.Len64(magnitude) - 1
	low := (magnitude - (uint64(1) << k)) * (mantissa >> k)
	normal := (word & sign) + uint64(k+1)*mantissa + low
	return fmt.Sprintf("%d -52 %d", normal, k-1074)
}
