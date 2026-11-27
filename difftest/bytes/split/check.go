package main

import (
	"bytes"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/bytesplit"
)

func main() {
	bytesplit.Compare("godiff-bytes-split", func(src, sep []byte, _ int) [][]byte {
		if len(src) == 0 && cap(src) == 0 {
			src = nil
		}
		return bytes.Split(src, sep)
	}, true, false)
}
