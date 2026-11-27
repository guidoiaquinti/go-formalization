package main

import (
	"fmt"
	"log"
	"unicode/utf8"
)

func main() {
	tests := []struct {
		input string
		want  bool
	}{
		{"", false},
		{"\x41", true},
		{"\xc0", true},
		{"\xc2", false},
		{"\xc2\x80", true},
		{"\xe0\x9f", true},
		{"\xe0\xa0", false},
		{"\xed\x9f", false},
		{"\xed\xa0", true},
		{"\xf0\x8f", true},
		{"\xf0\x90", false},
		{"\xf4\x8f\x80", false},
		{"\xf4\x90\x80", true},
		{"\xf0\x90\x7f", true},
		{"\xf0\x90\xbf", false},
		{"\xf0\x90\xc0", true},
		{"\xf0\x90\xbf\x00", true},
	}
	candidates := []struct {
		name string
		full func(string) bool
		want int
	}{
		{"Go", utf8.FullRuneInString, 0},
		{"length only", lengthOnly, 6},
		{"range-decoded rune exists", func(input string) bool {
			for range input {
				return true
			}
			return false
		}, 6},
		{"generic continuation range", genericContinuation, 4},
		{"ignore third continuation", ignoreThird, 2},
	}
	for _, candidate := range candidates {
		var mismatches int
		for _, test := range tests {
			if got := candidate.full(test.input); got != test.want {
				mismatches++
				if candidate.want == 0 {
					log.Fatalf("%s(% x) = %t, want %t", candidate.name, test.input, got, test.want)
				}
			}
		}
		if mismatches != candidate.want {
			log.Fatalf("%s disagreed on %d fixtures, want %d", candidate.name, mismatches, candidate.want)
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", candidate.name, mismatches, len(tests))
	}
}

func lengthOnly(input string) bool {
	return len(input) != 0 && len(input) >= leadLength(input[0])
}

func genericContinuation(input string) bool {
	if lengthOnly(input) {
		return true
	}
	for i := 1; i < len(input); i++ {
		b := input[i]
		if b < 0x80 || b > 0xbf {
			return true
		}
	}
	return false
}

func ignoreThird(input string) bool {
	if len(input) == 3 && leadLength(input[0]) == 4 {
		return utf8.FullRuneInString(input[:2])
	}
	return utf8.FullRuneInString(input)
}

func leadLength(lead byte) int {
	switch {
	case 0xc2 <= lead && lead <= 0xdf:
		return 2
	case 0xe0 <= lead && lead <= 0xef:
		return 3
	case 0xf0 <= lead && lead <= 0xf4:
		return 4
	default:
		return 1
	}
}
