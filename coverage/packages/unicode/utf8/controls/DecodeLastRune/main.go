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
		{[]byte("¢A"), 'A', 1},
		{[]byte("A¢"), '¢', 2},
		{[]byte("A🐈"), '🐈', 4},
		{[]byte("A\xef\xbf\xbd"), utf8.RuneError, 3},
		{[]byte("A\xe2\x82"), utf8.RuneError, 1},
		{[]byte("€\x80"), utf8.RuneError, 1},
		{[]byte("\xff€"), '€', 3},
		{[]byte("\x80\x80\x80\x80"), utf8.RuneError, 1},
	}
	var firstRuneFailures, allInvalidFailures int
	for _, f := range fixtures {
		original := bytes.Clone(f.input)
		r, n := utf8.DecodeLastRune(f.input)
		if r != f.value || n != f.width || !bytes.Equal(f.input, original) {
			log.Fatalf("DecodeLastRune(% x) = (%U, %d), want (%U, %d) with unchanged input", original, r, n, f.value, f.width)
		}
		if r, n := utf8.DecodeLastRuneInString(string(f.input)); r != f.value || n != f.width {
			log.Fatalf("string decoder disagreed for % x", f.input)
		}
		if r, n := utf8.DecodeRune(f.input); r != f.value || n != f.width {
			firstRuneFailures++
		}
		if !utf8.Valid(f.input) && (r != utf8.RuneError || n != 1) {
			allInvalidFailures++
		}
	}
	if firstRuneFailures != 7 || allInvalidFailures != 1 {
		log.Fatalf("first-rune mutant failed %d; reject-any-invalid-byte mutant failed %d", firstRuneFailures, allInvalidFailures)
	}
	fmt.Printf("Both Go decoders passed %d fixtures; mutants failed %d and %d\n", len(fixtures), firstRuneFailures, allInvalidFailures)
}
