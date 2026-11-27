package main

import (
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/utf8decode"
)

func main() {
	utf8decode.CompareValid("godiff-utf8-valid-string", false, func(data []byte) bool {
		return utf8.ValidString(string(data))
	})
}
