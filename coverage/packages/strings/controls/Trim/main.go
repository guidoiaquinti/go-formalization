package main

import (
	"fmt"
	"log"
	"strings"
)

func byteTrim(s, cutset string) string {
	for len(s) > 0 && strings.IndexByte(cutset, s[0]) >= 0 {
		s = s[1:]
	}
	for len(s) > 0 && strings.IndexByte(cutset, s[len(s)-1]) >= 0 {
		s = s[:len(s)-1]
	}
	return s
}
func main() {
	fixtures := []struct{ s, cutset, want string }{
		{"", "", ""}, {"abc", "", "abc"}, {"xxxx", "x", ""}, {"abbaXbab", "ab", "X"},
		{"é🚀Xé", "é🚀", "X"}, {"\xffé\xff", "\uFFFD", "é"},
		{"é", "\xa9", "é"}, {"a\xffb", "ab", "\xff"}, {"\uFFFD", "\xff", ""},
	}
	for _, c := range []struct {
		name string
		fn   func(string, string) string
		want int
	}{
		{"Go", strings.Trim, 0}, {"Trim by byte", byteTrim, 3},
		{"Normalize malformed input", func(s, cutset string) string { return strings.Trim(strings.ToValidUTF8(s, "\ufffd"), cutset) }, 1},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := c.fn(f.s, f.cutset); got != f.want {
				bad++
				if c.want == 0 {
					log.Fatalf("Trim(%q,%q) = %q, want %q", f.s, f.cutset, got, f.want)
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", c.name, bad, len(fixtures))
		if bad != c.want {
			log.Fatalf("%s wanted %d failures", c.name, c.want)
		}
	}
}
