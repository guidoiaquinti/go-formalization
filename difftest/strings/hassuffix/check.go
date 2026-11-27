package main

import (
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/affix"
)

func main() {
	affix.Compare("godiff-strings-hassuffix", true, strings.HasSuffix)
}
