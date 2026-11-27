package main

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"
)

func allowOverlongNUL(s string) bool {
	return utf8.ValidString(strings.ReplaceAll(s, "\xc0\x80", "\x00"))
}

func validFirstRune(s string) bool {
	r, width := utf8.DecodeRuneInString(s)
	return r != utf8.RuneError || width != 1
}

func main() {
	ascii := strings.Repeat("A", 65)
	fixtures := []struct {
		input string
		valid bool
	}{
		{"", true},
		{ascii, true},
		{ascii + "\x00\xef\xbf\xbd", true},
		{ascii + "\xf4\x8f\xbf\xbf", true},
		{"\xc0\x80", false},
		{ascii + "\xc0\x80", false},
		{"€\xed\xa0\x80", false},
		{ascii + "\xf0\x90\x80", false},
		{"\xff\x00", false},
	}
	var overlongFailures, firstRuneFailures int
	for _, f := range fixtures {
		if got := utf8.ValidString(f.input); got != f.valid {
			log.Fatalf("ValidString(% x) = %t, want %t", f.input, got, f.valid)
		}
		if allowOverlongNUL(f.input) != f.valid {
			overlongFailures++
		}
		if validFirstRune(f.input) != f.valid {
			firstRuneFailures++
		}
	}
	if overlongFailures != 2 || firstRuneFailures != 3 {
		log.Fatalf("allow-overlong mutant failed %d, first-rune mutant failed %d", overlongFailures, firstRuneFailures)
	}
	fmt.Printf("Go passed %d fixtures; mutants failed %d and %d\n", len(fixtures), overlongFailures, firstRuneFailures)
}
