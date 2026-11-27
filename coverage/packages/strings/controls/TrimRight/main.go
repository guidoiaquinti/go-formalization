package main

import (
	"fmt"
	"log"
	"strings"
)

func byteTrim(s, cutset string) string {
	for len(s) > 0 && strings.IndexByte(cutset, s[len(s)-1]) >= 0 {
		s = s[:len(s)-1]
	}
	return s
}

func main() {
	fixtures := []struct{ s, cutset, want string }{
		{"", "a", ""},
		{"abc", "", "abc"},
		{"xxxx", "x", ""},
		{"abbaXbab", "ab", "abbaX"},
		{"éX🚀é", "é🚀", "éX"},
		{"\xffé\xff", "\uFFFD", "\xffé"},
		{"é", "\xa9", "é"},
		{"a\xffb", "ab", "a\xff"},
		{"\uFFFD", "\xff", ""},
	}
	for _, c := range []struct {
		name string
		fn   func(string, string) string
		want int
	}{
		{"Go", strings.TrimRight, 0},
		{"Trim by byte", byteTrim, 3},
		{"Normalize malformed input", func(s, cutset string) string { return strings.TrimRight(strings.ToValidUTF8(s, "\ufffd"), cutset) }, 2},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := c.fn(f.s, f.cutset); got != f.want {
				bad++
				if c.want == 0 {
					log.Fatalf("TrimRight(%q,%q) = %q, want %q", f.s, f.cutset, got, f.want)
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", c.name, bad, len(fixtures))
		if bad != c.want {
			log.Fatalf("%s wanted %d failures", c.name, c.want)
		}
	}
}
