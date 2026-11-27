package main

import (
	"bytes"
	"fmt"
	"log"
	"math/bits"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	width        int
	word         uint64
	want, mutant string
}

func main() {
	cases := []fixture{
		{8, 0, "32", "0"}, {8, 0x40, "96", "64"}, {8, 0x5b, "123", "91"},
		{8, 0x80, "160", "128"}, {8, 0x20, "32", "0"},
	}
	for _, width := range []int{32, 64} {
		high := uint64(1) << (width - 1)
		cases = append(cases,
			fixture{width, 0, "true", "false"}, fixture{width, high, "true", "false"},
			fixture{width, high + 1, "false", "true"}, fixture{width, high - 1, "false", "true"},
			fixture{width, high / 2, "true", "false"})
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %d\n", c.width, c.word)
	}
	f, err := os.CreateTemp("", "golean-strconv-integer-review-")
	must(err)
	defer os.Remove(f.Name())
	rejected, native := 0, 0
	for i, got := range check.LeanRows("godiff-strconv-integer-bits", len(cases), &input) {
		c := cases[i]
		if got != c.want {
			log.Fatalf("row %d: Lean %s want %s", i, got, c.want)
		}
		if got != c.mutant {
			rejected++
		}
		if c.width == 8 || c.width == bits.UintSize {
			_, err := fmt.Fprintf(f, "%d %d %s\n", c.width, c.word, c.want)
			must(err)
			native++
		}
	}
	if rejected != len(cases) {
		log.Fatalf("rejected %d of %d mutants", rejected, len(cases))
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_STRCONV_BITS_ROWS", f.Name()))
	const path = "difftest/internal/strconv/integer-bits/testdata/"
	must(check.TestInPackage("internal/strconv", "TestGoLeanIntegerBits", path+"export_test.go", path+"integer_bits_test.go"))
	fmt.Printf("integer bits review: %d Lean rows, %d native rows, %d altered results rejected\n", len(cases), native, rejected)
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
