package main

import (
	"fmt"
	"log"
	"slices"
	"strings"
)

type split func(string, func(rune) bool) []string

func main() {
	fixtures := []struct {
		input, calls, separators string
		want                     []string
	}{
		{"", "", "|", nil},
		{"aaa", "aaa", "a", nil},
		{"ab", "ab", "", []string{"ab"}},
		{"|a||b|", "|a||b|", "|", []string{"a", "b"}},
		{"🐈αβ🐈x", "🐈αβ🐈x", "🐈", []string{"αβ", "x"}},
		{"a\xffé\xef\xbf\xbdz", "a\ufffdé\ufffdz", "\ufffd", []string{"a", "é", "z"}},
		{"a\xffz", "a\ufffdz", "", []string{"a\xffz"}},
	}
	for _, candidate := range []struct {
		name string
		run  split
		want int
	}{
		{"Go", strings.FieldsFunc, 0},
		{"Replace malformed bytes", func(s string, p func(rune) bool) []string {
			return strings.FieldsFunc(strings.ToValidUTF8(s, "\ufffd"), p)
		}, 1},
		{"Evaluate predicate per byte", func(s string, p func(rune) bool) []string {
			var fields []string
			start := 0
			for i := range len(s) {
				if p(rune(s[i])) {
					if start != i {
						fields = append(fields, s[start:i])
					}
					start = i + 1
				}
			}
			if start != len(s) {
				fields = append(fields, s[start:])
			}
			return fields
		}, 3},
	} {
		bad := 0
		for _, f := range fixtures {
			var calls strings.Builder
			got := candidate.run(f.input, func(r rune) bool {
				calls.WriteRune(r)
				return strings.ContainsRune(f.separators, r)
			})
			if calls.String() != f.calls || !slices.Equal(got, f.want) {
				bad++
				if candidate.want == 0 {
					log.Fatalf("FieldsFunc(%q) = %q with calls %q, want %q with calls %q", f.input, got, calls.String(), f.want, f.calls)
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", candidate.name, bad, len(fixtures))
		if bad != candidate.want {
			log.Fatalf("%s: want %d disagreements", candidate.name, candidate.want)
		}
	}
}
