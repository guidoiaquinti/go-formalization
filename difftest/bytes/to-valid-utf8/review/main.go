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

type fixture struct {
	text, replacement, want string
	growthChunks            []int
}

func main() {
	cases := []fixture{
		{"", "", "", nil},
		{"", "abc", "", nil},
		{"é�🚀", "\xff", "é�🚀", nil},
		{"\xff\xfc", "?", "?", nil},
		{"\xffé\xfc", "?", "?é?", nil},
		{"\xff�\xfe", "--", "--�--", nil},
		{"\xe2\x82X\xed\xa0\x80", "\xff!", "\xff!X\xff!", nil},
		{"\xffx\xff", "12345678", "12345678x12345678", []int{8, 1, 8}},
		{"\xffx\xffx\xff", "12345678", "12345678x12345678x12345678", []int{8, 1, 8, 1, 8}},
		{"\xffX\xfe", "", "X", nil},
		{"\xffa\xfe", "\xff\x80", "\xff\x80a\xff\x80", nil},
		{"\xff\xfe", "", "", nil},
	}
	var input bytes.Buffer
	var perByteFailures, sanitizedFailures, nilFailures, collapsedCaps int
	want := make([]string, len(cases))
	for i, c := range cases {
		row := strconv.AppendInt(nil, int64(len(c.replacement)), 10)
		for _, part := range []string{c.replacement, c.text} {
			for j := range len(part) {
				row = strconv.AppendUint(append(row, ' '), uint64(part[j]), 10)
			}
		}
		fmt.Fprintln(&input, string(row))
		initial := len(c.text) + len(c.replacement)
		got := bytes.ToValidUTF8([]byte(c.text), []byte(c.replacement))
		if got == nil || string(got) != c.want {
			log.Fatalf("ToValidUTF8(%q,%q)=%q nil=%t want %q", c.text, c.replacement, got, got == nil, c.want)
		}
		capHeader := strconv.Itoa(initial)
		wantCap := initial
		if c.growthChunks != nil {
			capHeader = "growth"
			staged := make([]byte, 0, initial)
			for _, n := range c.growthChunks {
				staged = append(staged, make([]byte, n)...)
			}
			if len(staged) != len(c.want) {
				log.Fatalf("fixture %d: wrong staged length %d", i, len(staged))
			}
			wantCap = cap(staged)
			if wantCap > len(c.want) {
				collapsedCaps++
			}
			fmt.Printf("growth fixture %d: initial=%d len=%d cap=%d\n", i, initial, len(got), cap(got))
		}
		if cap(got) != wantCap {
			log.Fatalf("fixture %d: cap=%d want %d", i, cap(got), wantCap)
		}
		want[i] = observe(c.want, capHeader)
		if perInvalidByte(c.text, c.replacement) != c.want {
			perByteFailures++
		}
		sanitized := bytes.ToValidUTF8([]byte(c.text), []byte(strings.ToValidUTF8(c.replacement, "?")))
		if string(sanitized) != c.want {
			sanitizedFailures++
		}
		if len(got) == 0 {
			nilFailures++
		}
	}
	for i, row := range check.LeanRows("godiff-bytes-to-valid-utf8", len(cases), &input) {
		if row != want[i] {
			log.Fatalf("fixture %d: Lean=%s want=%s", i, row, want[i])
		}
	}
	if perByteFailures != 2 || sanitizedFailures != 2 || nilFailures != 3 || collapsedCaps != 2 {
		log.Fatalf("per-invalid-byte failed %d, sanitized replacement failed %d, nil-empty failed %d, collapsed capacity cases %d",
			perByteFailures, sanitizedFailures, nilFailures, collapsedCaps)
	}
	fmt.Printf("ToValidUTF8 passed %d fixtures; per-invalid-byte failed %d, sanitized replacement failed %d, nil-empty failed %d; growth lower bound hides tight-cap mutation on %d fixtures\n",
		len(cases), perByteFailures, sanitizedFailures, nilFailures, collapsedCaps)
}

func perInvalidByte(text, replacement string) string {
	var result strings.Builder
	for i, r := range text {
		if r == utf8.RuneError && !strings.HasPrefix(text[i:], "\xef\xbf\xbd") {
			result.WriteString(replacement)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func observe(text, capacity string) string {
	row := append([]byte(capacity), ':')
	for i := range len(text) {
		if i > 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(text[i]), 10)
	}
	return string(row)
}
