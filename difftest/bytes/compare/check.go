package main

import (
	"bytes"
	"fmt"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/byteorder"
)

func main() {
	byteorder.Compare("godiff-bytes-compare", func(a, b []byte) string {
		return fmt.Sprintf("%t %d", bytes.Equal(a, b), bytes.Compare(a, b))
	}, func(order int) string {
		return fmt.Sprintf("%t %d", order == 0, order)
	})
}
