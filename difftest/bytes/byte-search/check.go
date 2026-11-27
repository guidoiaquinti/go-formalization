package main

import (
	"bytes"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/bytefind"
)

func main() {
	bytefind.Compare("godiff-bytes-byte-search", func(data []byte, key byte) (int, int) {
		return bytes.IndexByte(data, key), bytes.LastIndexByte(data, key)
	})
}
