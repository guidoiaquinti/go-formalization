package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	word, code, length         uint32
	packed, gotLength, gotCode uint32
	zero                       bool
	set                        uint32
}

func main() {
	cases := []fixture{
		{0, 0x1234, 0xfe, 0x1234fe, 0, 0, true, 0x1234fe},
		{0x80000000, 0, 0, 0, 0, 0x800000, false, 0},
		{0xff1234e7, 0x8000, 0x80, 0x800080, 0xe7, 0xff1234, false, 0x800080},
		{0x00000100, 0xffff, 0xff, 0xffffff, 0, 1, false, 0xffffff},
	}
	var input bytes.Buffer
	truth, err := os.CreateTemp("", "review-flate-hcode-")
	must(err)
	defer os.Remove(truth.Name())
	wants := make([]string, 0, len(cases))
	for _, c := range cases {
		args := fmt.Sprintf("%d %d %d", c.word, c.code, c.length)
		want := fmt.Sprintf("%d %d %d %t %d", c.packed, c.gotLength, c.gotCode, c.zero, c.set)
		wants = append(wants, want)
		fmt.Fprintln(&input, args)
		_, err := fmt.Fprintln(truth, args, want)
		must(err)
	}
	must(truth.Close())
	must(os.Setenv("GOLEAN_FLATE_HCODE_ROWS", truth.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanHCode", "difftest/compress/flate/hcode/testdata/hcode_test.go"))

	for i, line := range check.LeanRows("godiff-flate-hcode", len(cases), &input) {
		if line != wants[i] {
			log.Fatalf("fixture %d got %q, want %q", i, line, wants[i])
		}
	}
	var failures [5]int
	for i, c := range cases {
		mutants := []any{
			c.code | c.length<<16,           // swap code and length placement
			c.word & 0xf,                    // truncate length to a nibble
			(c.word >> 8) & 0xffff,          // drop the high byte from the extracted code
			c.word&0xff == 0,                // use zero length as the unset predicate
			c.word | (c.code<<8 | c.length), // preserve stale receiver bits
		}
		fields := strings.Fields(wants[i])
		for column, mutated := range mutants {
			if fmt.Sprint(mutated) != fields[column] {
				failures[column]++
			}
		}
	}
	if failures != [5]int{2, 1, 2, 2, 2} {
		log.Fatalf("unexpected mutation mismatch counts %v", failures)
	}
	log.Print("four pinned Go and Lean fixtures passed with mutation mismatch counts 2/1/2/2/2")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
