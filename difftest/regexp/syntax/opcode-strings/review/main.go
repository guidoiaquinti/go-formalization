package main

import (
	"bytes"
	"fmt"
	"log"
	"regexp/syntax"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		raw      uint8
		inst, op string
	}{
		{0, "InstAlt", "Op(0)"}, {1, "InstAltMatch", "NoMatch"}, {10, "InstRuneAnyNotNL", "EndText"},
		{11, "", "WordBoundary"}, {19, "", "Alternate"}, {20, "", "Op(20)"},
		{127, "", "Op(127)"}, {128, "", "opPseudo"}, {129, "", "Op(129)"}, {255, "", "Op(255)"},
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&input, c.raw)
	}
	var rejected [4]int
	for i, got := range check.LeanRows("godiff-regexp-syntax-opcode-strings", len(cases), &input) {
		c := cases[i]
		if want := c.inst + "|" + c.op; got != want || syntax.InstOp(c.raw).String() != c.inst || syntax.Op(c.raw).String() != c.op {
			log.Fatalf("opcode %d: pinned Go or Lean differs from literal %q", c.raw, want)
		}
		fallback, truncated := c.inst, c.inst
		missingPseudo, signed := c.op, c.op
		if c.raw >= 11 {
			fallback = fmt.Sprintf("InstOp(%d)", c.raw)
		}
		if c.raw >= 10 {
			truncated = ""
		}
		if c.raw == 128 {
			missingPseudo = "Op(128)"
		}
		if c.raw >= 20 && c.raw != 128 {
			signed = fmt.Sprintf("Op(%d)", int8(c.raw))
		}
		for j, wrong := range []string{fallback + "|" + c.op, truncated + "|" + c.op,
			c.inst + "|" + missingPseudo, c.inst + "|" + signed} {
			if wrong != got {
				rejected[j]++
			}
		}
	}
	if want := [4]int{7, 1, 1, 2}; rejected != want {
		log.Fatalf("mutation counts %v want %v", rejected, want)
	}
	fmt.Printf("opcode review: %d pinned native and Lean rows, mutation counts %v\n", len(cases), rejected)
}
