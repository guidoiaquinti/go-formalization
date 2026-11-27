package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

func row(lo, hi uint32, delta rune) unicode.CaseRange {
	return unicode.CaseRange{Lo: lo, Hi: hi, Delta: [3]rune{delta, 9, -9}}
}

func main() {
	cases := []struct {
		table       unicode.SpecialCase
		input, want string
	}{
		{nil, "i\xff\xe2\x82", "I���"},
		{unicode.SpecialCase{row('i', 'i', -1), row('i', 'i', 0), row('i', 'i', 1)}, "i", "i"},
		{unicode.SpecialCase{row('a', 'a', 0), row('z', 'z', 0), row('a', 'a', 7)}, "a", "a"},
		{unicode.SpecialCase{row('a', 'a', 7), row('z', 'z', 0), row('b', 'b', 0)}, "b", "B"},
		{unicode.SpecialCase{row(0x80000000, 0x7fffffff, unicode.UpperLower)}, "AB", "@B"},
		{unicode.SpecialCase{row('a', 'a', math.MinInt32)}, "a!", "!"},
		{unicode.SpecialCase{row('a', 'a', 0xd800-'a')}, "a", "�"},
		{unicode.SpecialCase{row(0xfffd, 0xfffd, '?'-0xfffd)}, "�\xff\xe2\x82", "????"},
		{unicode.TurkishCase, "iıİI", "İIİI"},
	}
	var input strings.Builder
	for i, c := range cases {
		if got := strings.ToUpperSpecial(c.table, c.input); got != c.want {
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
	if wrong := strings.Map(cases[1].table.ToLower, cases[1].input); wrong == cases[1].want {
		// The overlap row picks zero for upper and a distinct lower delta.
		panic("wrong case variant unexpectedly matched")
	}
	if wrong := strings.ToUpper(cases[1].input); wrong == cases[1].want {
		panic("fallback on zero hit unexpectedly matched")
	}
	midpoint := cases[1].table[0]
	if wrong := string(rune('i') + midpoint.Delta[unicode.UpperCase]); wrong == cases[1].want {
		panic("first matching range unexpectedly matched")
	}
	binary := "lean/.lake/build/bin/godiff-strings-to-upper-special"
	if len(os.Args) == 2 {
		binary = os.Args[1]
	}
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input.String())
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	observed := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(observed) != len(cases) {
		panic(fmt.Sprintf("got %d Lean rows", len(observed)))
	}
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
