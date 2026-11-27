// Command gen emits the pinned strconv.IsPrint intervals used by the Lean model.
// Run from the repository root using GOTOOLCHAIN=go1.27.0 GOWORK=off go run
// ./difftest/net/url/error-methods/gen.
package main

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

func main() {
	fmt.Println("-- Generated from Go 1.27.0 strconv.IsPrint by difftest/net/url/error-methods/gen.")
	fmt.Println("namespace GoLean.Net.Url")
	fmt.Println("def goQuotePrintRanges : Array (Nat × Nat) := #[")
	for r := rune(0); r <= utf8.MaxRune; r++ {
		if !strconv.IsPrint(r) {
			continue
		}
		start := r
		for r < utf8.MaxRune && strconv.IsPrint(r+1) {
			r++
		}
		fmt.Printf("  (%d, %d),\n", start, r)
	}
	fmt.Println("]")
	fmt.Println("end GoLean.Net.Url")
}
