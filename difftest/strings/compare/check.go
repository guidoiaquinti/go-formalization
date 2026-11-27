package main

import (
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/byteorder"
)

func main() {
	byteorder.Compare("godiff-strings-compare", func(a, b []byte) string {
		return strconv.Itoa(strings.Compare(string(a), string(b)))
	}, strconv.Itoa)
}
