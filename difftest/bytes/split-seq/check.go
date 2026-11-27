package main

import (
	"bytes"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/splitseq"
)

func main() {
	splitseq.Compare("godiff-bytes-split-seq", bytes.SplitSeq, strings.Split)
}
