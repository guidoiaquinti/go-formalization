package main

import (
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/bytefind"
)

func main() {
	bytefind.Compare("godiff-strings-byte-search", func(data []byte, key byte) (int, int) {
		text := string(data)
		return strings.IndexByte(text, key), strings.LastIndexByte(text, key)
	})
}
