package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		input, chars string
		want         int
	}{
		{"", "a", -1},
		{"abc", "", -1},
		{"\xa9", "é", -1},
		{"\xa9", "�", 0},
		{"\xa9", "\xff", 0},
		{"\xff", "\xc3\xa9", -1},
		{"\xfeé\xff", "\xfe", 3},
		{"é\x00", "\x00", 2},
		{"é🚀e", "e🚀", 6},
		{"éa🚀a", "a🚀", 7},
		{"é🚀a", "a", 6},
		{"é", "©", -1},
		{"a\xe2\x82�", "�", 3},
		{"a�\xe2\x82", "�", 5},
		{"q\xed\xa0\x80", "\xff", 3},
		{"🚀", "\x9f", -1},
	}
	for _, length := range []int{7, 8, 9, 64} {
		cases = append(cases, struct {
			input, chars string
			want         int
		}{strings.Repeat("z", length-2) + "aZ", "Za", length - 1})
	}
	var input bytes.Buffer
	raw, ordinal, first := 0, 0, 0
	for i, c := range cases {
		fmt.Fprint(&input, len(c.input))
		for _, part := range []string{c.input, c.chars} {
			for _, b := range []byte(part) {
				fmt.Fprint(&input, " ", b)
			}
		}
		fmt.Fprintln(&input)
		got := bytes.LastIndexAny([]byte(c.input), c.chars)
		if got != c.want {
			log.Fatalf("Go fixture %d returned %d, want %d", i, got, c.want)
		}
		rawOffset := -1
		for offset := len(c.input) - 1; offset >= 0; offset-- {
			if strings.IndexByte(c.chars, c.input[offset]) >= 0 {
				rawOffset = offset
				break
			}
		}
		if rawOffset != c.want {
			raw++
		}
		if got >= 0 && utf8.RuneCountInString(c.input[:got]) != got {
			ordinal++
		}
		if bytes.IndexAny([]byte(c.input), c.chars) != c.want {
			first++
		}
	}
	for i, row := range check.LeanRows("godiff-bytes-last-index-any", len(cases), &input) {
		if row != strconv.Itoa(cases[i].want) {
			log.Fatalf("Lean fixture %d returned %s, want %d", i, row, cases[i].want)
		}
	}
	if raw != 9 || ordinal != 6 || first != 10 {
		log.Fatalf("mutation counts raw=%d ordinal=%d first=%d, want 9, 6 and 10", raw, ordinal, first)
	}
	fmt.Printf("bytes.LastIndexAny passed %d fixed Go and Lean fixtures; raw byte set, rune ordinal and first occurrence mutants failed %d, %d and %d.\n",
		len(cases), raw, ordinal, first)
}
