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
		input, chars string
		want         bool
	}{
		{"", "a", false},
		{"abc", "", false},
		{"\xa9", "é", false},
		{"\xa9", "�", true},
		{"\xa9", "\xff", true},
		{"\xff", "\xc3\xa9", false},
		{"é\xff", "\xfe", true},
		{"é\x00", "\x00", true},
		{"é🚀e", "e🚀", true},
		{"é", "©", false},
		{"a\xe2\x82", "�", true},
		{"a�", "\xe2\x82", true},
		{"q\xed\xa0\x80", "\xff", true},
		{"🚀", "\x9f", false},
		{"aa", "a", true},
		{strings.Repeat("z", 8), "za", true},
		{strings.Repeat("z", 9), "za", true},
		{strings.Repeat("z", 64), "aA", false},
	}
	var input bytes.Buffer
	raw, zero := 0, 0
	for i, c := range cases {
		fmt.Fprint(&input, len(c.input))
		for _, part := range []string{c.input, c.chars} {
			for _, b := range []byte(part) {
				fmt.Fprint(&input, " ", b)
			}
		}
		fmt.Fprintln(&input)
		if got := bytes.ContainsAny([]byte(c.input), c.chars); got != c.want {
			log.Fatalf("Go fixture %d returned %t, want %t", i, got, c.want)
		}
		rawFound := false
		for _, b := range []byte(c.input) {
			rawFound = rawFound || strings.IndexByte(c.chars, b) >= 0
		}
		if rawFound != c.want {
			raw++
		}
		if (bytes.IndexAny([]byte(c.input), c.chars) > 0) != c.want {
			zero++
		}
	}
	for i, row := range check.LeanRows("godiff-bytes-contains-any", len(cases), &input) {
		if row != strconv.FormatBool(cases[i].want) {
			log.Fatalf("Lean fixture %d returned %s, want %t", i, row, cases[i].want)
		}
	}
	if raw != 9 || zero != 5 {
		log.Fatalf("mutation counts raw=%d zero=%d, want 9 and 5", raw, zero)
	}
	fmt.Printf("bytes.ContainsAny passed %d fixed Go and Lean fixtures; raw byte set and excluding offset zero mutants failed %d and %d.\n",
		len(cases), raw, zero)
}
