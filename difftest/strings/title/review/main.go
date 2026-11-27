package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func wire(s string) string {
	fields := make([]byte, 0, len(s)*4)
	for i := range len(s) {
		if i > 0 {
			fields = append(fields, ',')
		}
		fields = fmt.Appendf(fields, "%d", s[i])
	}
	return "bytes:" + string(fields)
}

func main() {
	cases := []struct{ input, want string }{
		{"", ""},
		{"one-two_three", "One-Two_three"},
		{"a'b", "A'B"},
		{"a—b", "A—b"},
		{"a\u00a0b", "A\u00a0B"},
		{"a\u200bb", "A\u200bb"},
		{"1a_2b.3c d", "1a_2b.3c D"},
		{"a\x00b", "A\x00B"},
		{"ǳ ǳǳ", "ǲ ǲǳ"},
		{"\xffabc", "�abc"},
		{"\xe2\x82 a", "�� A"},
		{"A�b\xff-c", "A�b�-C"},
	}
	var input bytes.Buffer
	for i, c := range cases {
		if got := strings.Title(c.input); got != c.want {
			panic(fmt.Sprintf("native row %d: got %q want %q", i, got, c.want))
		}
		input.WriteByte('0')
		for _, b := range []byte(c.input) {
			fmt.Fprintf(&input, " %d", b)
		}
		fmt.Fprintln(&input)
	}
	binary := "lean/.lake/build/bin/godiff-strings-title"
	if len(os.Args) == 2 {
		binary = os.Args[1]
	}
	cmd := exec.Command(binary)
	cmd.Stdin = &input
	raw, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	rows := strings.Split(strings.TrimSuffix(string(raw), "\n"), "\n")
	if len(rows) != len(cases) {
		panic(fmt.Sprintf("Lean produced %d rows", len(rows)))
	}
	for i, c := range cases {
		if rows[i] != wire(c.want) {
			panic(fmt.Sprintf("Lean row %d: got %q want %q", i, rows[i], wire(c.want)))
		}
	}
	for _, mutant := range []func(string) string{strings.ToTitle, strings.ToUpper, func(s string) string { return s }} {
		mismatch := 0
		for _, c := range cases {
			if mutant(c.input) != c.want {
				mismatch++
			}
		}
		if mismatch == 0 {
			panic("whole-string case or identity mutant survived")
		}
	}
	fmt.Printf("strings.Title passed %d independent pinned Go and Lean fixtures\n", len(cases))
}
