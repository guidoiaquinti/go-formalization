package main

import (
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/utf8decode"
)

func main() {
	utf8decode.Compare("godiff-utf8-decode-rune-in-string", func(data []byte) (rune, int) {
		return utf8.DecodeRuneInString(string(data))
	})
}
