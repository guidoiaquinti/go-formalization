package main

import (
	"github.com/guidoiaquinti/go-formalization/difftest/internal/utf8decode"
	"unicode/utf8"
)

func main() {
	utf8decode.CompareLast("godiff-utf8-decode-last-rune-in-string", false, func(bytes []byte) (rune, int) {
		return utf8.DecodeLastRuneInString(string(bytes))
	})
}
