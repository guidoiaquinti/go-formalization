package main

import (
	"bytes"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/bytesplit"
)

func main() {
	bytesplit.Compare("godiff-bytes-split-n", bytes.SplitN, false, false)
}
