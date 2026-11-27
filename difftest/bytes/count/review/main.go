package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		input, sep string
		want       int
	}{
		{"", "", 1},
		{"", "x", 0},
		{"a", "", 2},
		{"é", "", 2},
		{"🐈", "", 2},
		{"\xe2\x82", "", 3},
		{"\xffé\xfe", "", 4},
		{"�", "", 2},
		{"\xed\xa0\x80", "", 4},
		{"aé🐈\xe2\x82", "", 6},
		{"aaaa", "aa", 2},
		{"ababa", "aba", 1},
		{"\x00\x00\x00", "\x00\x00", 1},
		{"aaa", "aa", 1},
		{"xxabcabc", "abc", 2},
		{"ab", "abc", 0},
		{"xéyé", "\xa9", 2},
		{"\xff\xff�", "\xff", 2},
		{"x🐈x🐈", "\x90\x88", 2},
	}
	var rows bytes.Buffer
	overlapFailures, byteLengthFailures := 0, 0
	for i, c := range cases {
		if got := bytes.Count([]byte(c.input), []byte(c.sep)); got != c.want {
			log.Fatalf("Go fixture %d returned %d, want %d", i, got, c.want)
		}
		fmt.Fprint(&rows, len(c.sep))
		for _, part := range []string{c.sep, c.input} {
			for j := range len(part) {
				fmt.Fprint(&rows, " ", part[j])
			}
		}
		fmt.Fprintln(&rows)
		if c.sep == "" {
			if len(c.input)+1 != c.want {
				byteLengthFailures++
			}
			continue
		}
		overlapping := 0
		for j := 0; j+len(c.sep) <= len(c.input); j++ {
			if c.input[j:j+len(c.sep)] == c.sep {
				overlapping++
			}
		}
		if overlapping != c.want {
			overlapFailures++
		}
	}
	for i, result := range check.LeanRows("godiff-bytes-count", len(cases), &rows) {
		if result != strconv.Itoa(cases[i].want) {
			log.Fatalf("Lean fixture %d returned %s, want %d", i, result, cases[i].want)
		}
	}
	if overlapFailures != 4 || byteLengthFailures != 5 {
		log.Fatalf("mutation failures overlap=%d byte length=%d, want 4 and 5", overlapFailures, byteLengthFailures)
	}
	fmt.Printf("bytes.Count passed %d fixed Go and Lean fixtures; overlap and byte-length mutants failed %d and %d.\n", len(cases), overlapFailures, byteLengthFailures)
}
