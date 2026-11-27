package main

import (
	"fmt"
	"log"
	"strings"
	"unicode"
	"unicode/utf8"
)

func main() {
	fixtures := []struct{ input, want string }{
		{"", ""},
		{"\t\n\v\f\r ", ""},
		{" \tX Y\r ", "X Y"},
		{" \u0085\u00a0\u3000\t", ""},
		{" \u2003X\u200bY\u2028 ", "X\u200bY"},
		{" \tX\u202f ", "X"},
		{" \u202fX \t", "X"},
		{" \ufeffX\u200b ", "\ufeffX\u200b"},
		{" \xffX\xe2\x82 ", "\xffX\xe2\x82"},
		{" \u0085�\xff\u00a0 ", "�\xff"},
		{" \x00X\u180e ", "\x00X\u180e"},
		{" \u2029X\u205fY\u1680 ", "X\u205fY"},
	}
	for _, c := range []struct {
		name string
		trim func(string) string
		want int
	}{
		{"Go", strings.TrimSpace, 0},
		{"ASCII whitespace only", func(s string) string { return strings.Trim(s, " \t\n\v\f\r") }, 6},
		{"RuneError is whitespace", func(s string) string {
			return strings.TrimFunc(s, func(r rune) bool { return unicode.IsSpace(r) || r == utf8.RuneError })
		}, 2},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := c.trim(f.input); got != f.want {
				bad++
				if c.want == 0 {
					log.Fatalf("TrimSpace(%q) = %q, want %q", f.input, got, f.want)
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", c.name, bad, len(fixtures))
		if bad != c.want {
			log.Fatalf("%s wanted %d failures", c.name, c.want)
		}
	}
}
