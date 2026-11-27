package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func asciiOnly(input string) string {
	raw := []byte(input)
	for i, b := range raw {
		if b >= 'a' && b <= 'z' {
			raw[i] = b - ('a' - 'A')
		}
	}
	return string(raw)
}

func wire(input string) string {
	row := []byte("bytes:")
	for i := range len(input) {
		if i > 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(input[i]), 10)
	}
	return string(row)
}

func main() {
	fixtures := []struct{ input, want string }{
		{"", ""}, {"\x00aAZ\x7f", "\x00AAZ\x7f"}, {"ǳǆ", "ǲǅ"},
		{"σς", "ΣΣ"}, {"ı", "I"}, {"𐐨", "𐐀"}, {"\xff\x80a", "��A"},
		{"\xe2\x82", "��"}, {"é", "É"}, {"�É", "�É"},
	}
	var rows bytes.Buffer
	for _, f := range fixtures {
		fmt.Fprint(&rows, "1")
		for i := range len(f.input) {
			fmt.Fprint(&rows, " ", f.input[i])
		}
		fmt.Fprintln(&rows)
		if got := strings.ToTitle(f.input); got != f.want {
			log.Fatalf("ToTitle(%x)=%x, want %x", f.input, got, f.want)
		}
	}
	for i, got := range check.LeanRows("godiff-strings-to-title", len(fixtures), &rows) {
		if want := wire(fixtures[i].want); got != want {
			log.Fatalf("Lean fixture %d gave %q, want %q", i, got, want)
		}
	}
	for _, c := range []struct {
		name  string
		upper func(string) string
		want  int
	}{
		{"ASCII only", asciiOnly, 7},
		{"Use uppercase instead", strings.ToUpper, 1},
		{"Collapse malformed runs", func(s string) string { return strings.ToTitle(strings.ToValidUTF8(s, "�")) }, 2},
	} {
		bad := 0
		for _, f := range fixtures {
			if c.upper(f.input) != f.want {
				bad++
			}
		}
		if bad != c.want {
			log.Fatalf("%s disagreed on %d, want %d", c.name, bad, c.want)
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", c.name, bad, len(fixtures))
	}
	fmt.Printf("%d Go and Lean titlecase fixtures passed\n", len(fixtures))
}
