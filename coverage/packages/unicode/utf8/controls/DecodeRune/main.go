package main

import (
	"bytes"
	"fmt"
	"log"
	"unicode/utf8"
)

func main() {
	fixtures := []struct {
		input []byte
		value rune
		width int
	}{
		{nil, utf8.RuneError, 0},
		{[]byte{}, utf8.RuneError, 0},
		{[]byte{0x7f, 0xff}, 0x7f, 1},
		{[]byte{0xe0, 0xa0, 0x80, 0xff}, 0x800, 3},
		{[]byte{0xe0, 0x9f, 0xbf}, utf8.RuneError, 1},
		{[]byte{0xef, 0xbf, 0xbd}, utf8.RuneError, 3},
		{[]byte{0xf0, 0x90, 0x80}, utf8.RuneError, 1},
		{[]byte{0xf4, 0x8f, 0xbf, 0xbf}, 0x10ffff, 4},
	}
	var lengthFailures, errorValueFailures int
	for _, f := range fixtures {
		before := bytes.Clone(f.input)
		r, n := utf8.DecodeRune(f.input)
		if r != f.value || n != f.width || !bytes.Equal(before, f.input) {
			log.Fatalf("DecodeRune(% x) = (%U, %d), want (%U, %d), input must stay intact",
				before, r, n, f.value, f.width)
		}
		if r, n := utf8.DecodeRuneInString(string(f.input)); r != f.value || n != f.width {
			log.Fatalf("DecodeRuneInString(% x) = (%U, %d), want (%U, %d)",
				before, r, n, f.value, f.width)
		}
		if len(f.input) != f.width {
			lengthFailures++
		}
		if r == utf8.RuneError && n > 1 {
			errorValueFailures++
		}
	}
	if lengthFailures != 4 || errorValueFailures != 1 {
		log.Fatalf("mutant failures: entire slice width %d, collapse replacement width %d", lengthFailures, errorValueFailures)
	}
	fmt.Printf("Go decoders passed %d fixtures; entire slice width and collapse replacement width mutants failed %d and %d\n",
		len(fixtures), lengthFailures, errorValueFailures)
}
