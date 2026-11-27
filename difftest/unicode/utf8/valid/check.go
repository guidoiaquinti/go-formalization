package main

import (
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/utf8decode"
)

func main() {
	utf8decode.CompareValid("godiff-utf8-valid", true, utf8.Valid)
}
