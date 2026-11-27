package main

import (
	"fmt"
	"log"
	"slices"
	"strings"
)

func main() {
	fixtures := []struct {
		input string
		want  []string
	}{
		{"", nil},
		{"\t\r\v\f", nil},
		{"  a\t\tb  ", []string{"a", "b"}},
		{"\u0085\u1680\u3000", nil},
		{"a\u00a0\xff\u202fb", []string{"a", "\xff", "b"}},
		{"a\u200bb\ufeff", []string{"a\u200bb\ufeff"}},
		{"a\x00b", []string{"a\x00b"}},
	}
	for _, candidate := range []struct {
		name string
		run  func(string) []string
		want int
	}{
		{"Go", strings.Fields, 0},
		{"Split on literal spaces", func(s string) []string { return strings.Split(s, " ") }, 5},
		{"Only ASCII whitespace", func(s string) []string {
			return strings.FieldsFunc(s, func(r rune) bool { return r == ' ' || '\t' <= r && r <= '\r' })
		}, 2},
		{"Replace malformed bytes", func(s string) []string { return strings.Fields(strings.ToValidUTF8(s, "\uFFFD")) }, 1},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := candidate.run(f.input); !slices.Equal(got, f.want) {
				bad++
				if candidate.want == 0 {
					log.Fatalf("Fields(%q) = %q, want %q", f.input, got, f.want)
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", candidate.name, bad, len(fixtures))
		if bad != candidate.want {
			log.Fatalf("%s: wanted %d mismatches", candidate.name, candidate.want)
		}
	}
}
