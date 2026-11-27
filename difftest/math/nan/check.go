package main

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	// Positive sign, all exponent bits set, quiet bit set, and payload one.
	const want uint64 = 2047<<52 + 1<<51 + 1
	actual := math.Float64bits(math.NaN())
	model := check.LeanRows("godiff-nan", 1, nil)[0]
	if actual != want || model != strconv.FormatUint(want, 10) {
		log.Fatalf("NaN: Go bits=%#016x Lean decimal=%q want bits=%#016x", actual, model, want)
	}
	fmt.Printf("math.NaN: Go and Lean agree on the only input, bits=%#016x\n", want)
}
