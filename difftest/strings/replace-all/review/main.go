package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		input, old, replacement, want string
	}{
		{"aaaa", "aa", "x", "xx"},
		{"ababa", "aba", "X", "Xba"},
		{"aaa", "a", "aa", "aaaaaa"},
		{"", "", ":", ":"},
		{"é🔧", "", ".", ".é.🔧."},
		{"\xe2\x82", "", "-", "-\xe2-\x82-"},
		{"\xff�", "", "!", "!\xff!�!"},
		{"éé", "\xa9", "X", "\xc3X\xc3X"},
		{"ab", "abc", "!", "ab"},
		{"XX", "X", "", ""},
	}
	var rows bytes.Buffer
	single, byteInsertion := 0, 0
	for i, c := range cases {
		if got := strings.ReplaceAll(c.input, c.old, c.replacement); got != c.want {
			log.Fatalf("Go fixture %d: %q, want %q", i, got, c.want)
		}
		fmt.Fprint(&rows, len(c.old), " ", len(c.replacement))
		for _, part := range []string{c.old, c.replacement, c.input} {
			for j := range len(part) {
				fmt.Fprint(&rows, " ", part[j])
			}
		}
		fmt.Fprintln(&rows)
		if strings.Replace(c.input, c.old, c.replacement, 1) != c.want {
			single++
		}
		if c.old == "" {
			var mutant strings.Builder
			mutant.WriteString(c.replacement)
			for j := range len(c.input) {
				mutant.WriteByte(c.input[j])
				mutant.WriteString(c.replacement)
			}
			if mutant.String() != c.want {
				byteInsertion++
			}
		}
	}
	for i, got := range check.LeanRows("godiff-strings-replace-all", len(cases), &rows) {
		fields := make([]string, len(cases[i].want))
		for j, b := range []byte(cases[i].want) {
			fields[j] = strconv.Itoa(int(b))
		}
		if want := "bytes:" + strings.Join(fields, ","); got != want {
			log.Fatalf("Lean fixture %d: %q, want %q", i, got, want)
		}
	}
	if single != 7 || byteInsertion != 2 {
		log.Fatalf("mutation counts single=%d byte insertion=%d, want 7 and 2", single, byteInsertion)
	}
	fmt.Printf("strings.ReplaceAll passed %d fixed Go and Lean fixtures; single and byte-insertion mutants failed %d and %d.\n",
		len(cases), single, byteInsertion)
}
