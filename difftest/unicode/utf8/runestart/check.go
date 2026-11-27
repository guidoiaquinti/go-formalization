package main

import (
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	check.Exhaustive("godiff-utf8-runestart", 1<<8, func(b int) (actual, expected int) {
		if utf8.RuneStart(byte(b)) {
			actual = 1
		}
		if b < 128 || 192 <= b {
			expected = 1
		}
		return actual, expected
	})
}
