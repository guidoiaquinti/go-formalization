package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		word uint32
		want string
	}{
		{0, "0 0 0 0"},
		{0xffffffff, "3221225472 255 4194303 255"},
		{0x800100a5, "2147483648 165 65701 0"},
		{0x407fffff, "1073741824 255 4194303 1"},
	}
	var input bytes.Buffer
	file, err := os.CreateTemp("", "review-flate-token-fields-")
	must(err)
	defer os.Remove(file.Name())
	for _, c := range cases {
		fmt.Fprintln(&input, c.word)
		_, err := fmt.Fprintln(file, c.word, c.want)
		must(err)
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_FLATE_TOKEN_ROWS", file.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanTokenFields",
		"difftest/compress/flate/token-fields/testdata/token_test.go"))
	var failures [4]int
	for i, row := range check.LeanRows("godiff-flate-token-fields", len(cases), &input) {
		c := cases[i]
		if row != c.want {
			log.Fatalf("word %#08x got %q, want %q", c.word, row, c.want)
		}
		mutants := [4]uint32{
			c.word >> 30,          // shift the type code to the low bits
			c.word & 0xffff,       // use the low offset-shaped word as a literal
			c.word & 0xffff,       // drop the embedded offset code
			(c.word >> 22) & 0x7f, // drop the top length bit
		}
		fields := strings.Fields(row)
		for i, mutant := range mutants {
			if fmt.Sprint(mutant) != fields[i] {
				failures[i]++
			}
		}
	}
	if failures != [4]int{3, 2, 3, 1} {
		log.Fatalf("unexpected mutation mismatch counts %v", failures)
	}
	log.Print("four pinned Go and Lean cases passed with mutation mismatch counts 3/2/3/1")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
