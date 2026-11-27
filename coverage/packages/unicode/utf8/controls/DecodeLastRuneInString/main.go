package main

import (
	"fmt"
	"log"
	"unicode/utf8"
)

func acceptPartialSuffix(s string) (rune, int) {
	if len(s) == 0 || s[len(s)-1] < utf8.RuneSelf {
		return utf8.DecodeLastRuneInString(s)
	}
	start := len(s) - 1
	limit := max(0, len(s)-utf8.UTFMax)
	for start--; start >= limit; start-- {
		if utf8.RuneStart(s[start]) {
			break
		}
	}
	return utf8.DecodeRuneInString(s[max(start, 0):])
}

func main() {
	fixtures := []struct {
		input string
		value rune
		width int
	}{
		{"", utf8.RuneError, 0},
		{"\xffA", 'A', 1},
		{"A¢", '¢', 2},
		{"A€", '€', 3},
		{"A🐈", '🐈', 4},
		{"A\xef\xbf\xbd", utf8.RuneError, 3},
		{"A\xe2\x82", utf8.RuneError, 1},
		{"🐈\x80", utf8.RuneError, 1},
		{"€\x80", utf8.RuneError, 1},
		{"A\x80", utf8.RuneError, 1},
		{"\x80\x80\x80\x80", utf8.RuneError, 1},
	}
	var partialFailures int
	for _, f := range fixtures {
		r, n := utf8.DecodeLastRuneInString(f.input)
		if r != f.value || n != f.width {
			log.Fatalf("DecodeLastRuneInString(% x) = (%U, %d), want (%U, %d)",
				f.input, r, n, f.value, f.width)
		}
		if r, n := acceptPartialSuffix(f.input); r != f.value || n != f.width {
			partialFailures++
		}
	}
	if partialFailures != 3 {
		log.Fatalf("partial suffix mutant failed %d fixtures, want 3", partialFailures)
	}
	fmt.Printf("Go passed %d fixtures; accepting a partial decoded suffix failed %d\n", len(fixtures), partialFailures)
}
