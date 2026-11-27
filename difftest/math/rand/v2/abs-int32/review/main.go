package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		word, magnitude uint32
	}{
		{0, 0},
		{0x40000000, 0x40000000},
		{0x7fffffff, 0x7fffffff},
		{0x80000000, 0x80000000},
		{0xc0000000, 0x40000000},
		{0xffffffff, 1},
	}
	var input bytes.Buffer
	truth, err := os.CreateTemp("", "review-abs-int32-")
	must(err)
	defer os.Remove(truth.Name())
	for _, c := range cases {
		fmt.Fprintln(&input, c.word)
		_, err = fmt.Fprintln(truth, c.word, c.magnitude)
		must(err)
	}
	must(truth.Close())
	must(os.Setenv("GOLEAN_ABS_INT32_ROWS", truth.Name()))
	must(check.TestInPackage("math/rand/v2", "TestGoLeanAbsInt32",
		"difftest/math/rand/v2/abs-int32/testdata/abs_test.go"))

	var missingNegation, saturation, clearSign int
	for i, row := range check.LeanRows("godiff-rand-v2-abs-int32", len(cases), &input) {
		c := cases[i]
		got, err := strconv.ParseUint(row, 10, 32)
		must(err)
		if uint32(got) != c.magnitude {
			log.Fatalf("word %#08x got %d, want %d", c.word, got, c.magnitude)
		}
		recovered := uint32(got)
		if int32(c.word) < 0 {
			recovered = -recovered
		}
		if recovered != c.word {
			log.Fatalf("word %#08x was recovered as %#08x", c.word, recovered)
		}
		if c.word != c.magnitude {
			missingNegation++
		}
		if min(c.magnitude, math.MaxInt32) != c.magnitude {
			saturation++
		}
		if c.word&math.MaxInt32 != c.magnitude {
			clearSign++
		}
	}
	if missingNegation != 2 || saturation != 1 || clearSign != 2 {
		log.Fatalf("unexpected mutant mismatch counts %d/%d/%d", missingNegation, saturation, clearSign)
	}
	fmt.Println("six fixed cases passed, with mismatch counts 2/1/2 for identity, saturation and clearing the sign bit")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
