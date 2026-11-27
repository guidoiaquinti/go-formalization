package main

import (
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/fullrune"
)

func main() {
	fullrune.Compare("godiff-utf8-full-rune-in-string", func(data []byte) bool {
		return utf8.FullRuneInString(string(data))
	})
}
