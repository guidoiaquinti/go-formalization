package main

import (
	"bytes"
	"fmt"
	"log"
	"math/bits"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		word  uint16
		count byte
		want  uint16
	}{
		{0x8001, 0, 0x8001},
		{0x8001, 16, 0x8001},
		{0x0001, 17, 0x0001},
		{0x8001, 255, 0x4000},
		{0xc005, 3, 0x0005},
		{0xc005, 19, 0x0005},
	}
	var input bytes.Buffer
	file, err := os.CreateTemp("", "review-flate-reverse-")
	must(err)
	defer os.Remove(file.Name())
	for _, c := range cases {
		fmt.Fprintln(&input, c.word, c.count)
		_, err := fmt.Fprintln(file, c.word, c.count, c.want)
		must(err)
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_FLATE_REVERSE_ROWS", file.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanReverseBits",
		"difftest/compress/flate/reverse-bits/testdata/reverse_test.go"))

	var failures [3]int
	for i, row := range check.LeanRows("godiff-flate-reverse-bits", len(cases), &input) {
		c := cases[i]
		got, err := strconv.ParseUint(row, 10, 16)
		must(err)
		if uint16(got) != c.want {
			log.Fatalf("word %#04x count %d got %d, want %d", c.word, c.count, got, c.want)
		}
		width := c.count % 16
		mutants := [3]uint16{
			bits.Reverse16(c.word << (16 - c.count)), // omit the shift mask
			bits.Reverse16(c.word),                   // ignore the requested width
			bits.Reverse16(c.word) >> (16 - width),   // treat zero residue as zero selected bits
		}
		for index, mutated := range mutants {
			if mutated != uint16(got) {
				failures[index]++
			}
		}
	}
	if failures != [3]int{4, 4, 2} {
		log.Fatalf("unexpected mutation mismatch counts %v", failures)
	}
	log.Print("six pinned Go and Lean fixtures passed with mutation mismatch counts 4/4/2")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
