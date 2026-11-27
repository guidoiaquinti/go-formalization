package main

import (
	"fmt"
	"log"
	"unicode/utf8"
)

func main() {
	tests := []struct {
		input []byte
		want  bool
	}{
		{nil, false},
		{[]byte{}, false},
		{[]byte{0x41}, true},
		{[]byte{0xc0}, true},
		{[]byte{0xc2}, false},
		{[]byte{0xc2, 0x80}, true},
		{[]byte{0xe0, 0x9f}, true},
		{[]byte{0xe0, 0xa0}, false},
		{[]byte{0xed, 0x9f}, false},
		{[]byte{0xed, 0xa0}, true},
		{[]byte{0xf0, 0x8f}, true},
		{[]byte{0xf0, 0x90}, false},
		{[]byte{0xf4, 0x8f, 0x80}, false},
		{[]byte{0xf4, 0x90, 0x80}, true},
		{[]byte{0xf0, 0x90, 0x7f}, true},
		{[]byte{0xf0, 0x90, 0xbf}, false},
		{[]byte{0xf0, 0x90, 0xc0}, true},
		{[]byte{0xf0, 0x90, 0xbf, 0x00}, true},
	}
	candidates := []struct {
		name string
		full func([]byte) bool
		want int
	}{
		{"Go", utf8.FullRune, 0},
		{"length only", lengthOnly, 6},
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

func lengthOnly(input []byte) bool {
	return len(input) != 0 && len(input) >= leadLength(input[0])
}

func genericContinuation(input []byte) bool {
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

func ignoreThird(input []byte) bool {
	if len(input) == 3 && leadLength(input[0]) == 4 {
		return utf8.FullRune(input[:2])
	}
	return utf8.FullRune(input)
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
