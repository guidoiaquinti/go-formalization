package main

import (
	"fmt"
	"math"
	"strings"
	"unicode"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func row(lo, hi uint32, delta rune) unicode.CaseRange {
	return unicode.CaseRange{Lo: lo, Hi: hi, Delta: [3]rune{9, delta, -9}}
}

func main() {
	cases := []struct {
		table       unicode.SpecialCase
		input, want string
	}{
		{nil, "I\xff\xe2\x82", "i���"},
		{unicode.SpecialCase{row('I', 'I', -1), row('I', 'I', 0), row('I', 'I', 1)}, "I", "I"},
		{unicode.SpecialCase{row('A', 'A', 0), row('Z', 'Z', 0), row('A', 'A', 7)}, "A", "A"},
		{unicode.SpecialCase{row('A', 'A', 7), row('Z', 'Z', 0), row('B', 'B', 0)}, "B", "b"},
		{unicode.SpecialCase{row(0x80000000, 0x7fffffff, unicode.UpperLower)}, "AB", "AC"},
		{unicode.SpecialCase{row('A', 'A', math.MinInt32)}, "A!", "!"},
		{unicode.SpecialCase{row('A', 'A', 0xd800-'A')}, "A", "�"},
		{unicode.SpecialCase{row(0xfffd, 0xfffd, '?'-0xfffd)}, "�\xff\xe2\x82", "????"},
	}
	var input strings.Builder
	for i, c := range cases {
		if got := strings.ToLowerSpecial(c.table, c.input); got != c.want {
			panic(fmt.Sprintf("native row %d: got %q want %q", i, got, c.want))
		}
		fmt.Fprintf(&input, "%d", len(c.table))
		for _, r := range c.table {
			fmt.Fprintf(&input, " %d %d %d %d %d", r.Lo, r.Hi, r.Delta[0], r.Delta[1], r.Delta[2])
		}
		for _, b := range []byte(c.input) {
			fmt.Fprintf(&input, " %d", b)
		}
		fmt.Fprintln(&input)
	}
	observed := check.LeanRows("godiff-strings-to-lower-special", len(cases), strings.NewReader(input.String()))
	for i, c := range cases {
		fields := make([]string, len(c.want))
		for j, b := range []byte(c.want) {
			fields[j] = fmt.Sprint(b)
		}
		want := "bytes:" + strings.Join(fields, ",")
		if observed[i] != want {
			panic(fmt.Sprintf("Lean row %d: got %q want %q", i, observed[i], want))
		}
	}
	fmt.Printf("%d independent pinned Go and Lean witnesses passed\n", len(cases))
}
