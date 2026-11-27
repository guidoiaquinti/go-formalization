package main

import (
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/affix"
)

func main() {
	affix.CompareTrim("godiff-strings-trimsuffix", true, strings.TrimSuffix)
}
