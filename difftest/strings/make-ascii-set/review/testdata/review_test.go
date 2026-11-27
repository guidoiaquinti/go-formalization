package strings_test

import (
	"strings"
	"testing"
)

type makeSet func(string) ([256]bool, bool)

func TestGoLeanMakeASCIISetReview(t *testing.T) {
	cases := []struct {
		input, marked string
		ok            bool
	}{
		{"", "", true}, {"\x00aa\x7f", "\x00a\x7f", true},
		{"ab\x80z", "ab", false}, {"\xffz", "", false}, {"abéz", "ab", false},
	}
	for _, c := range []struct {
		name string
		fn   makeSet
		want int
	}{
		{"Go", strings.GoLeanMakeASCIISet, 0},
		{"Clear table on failure", func(s string) ([256]bool, bool) {
			as, ok := strings.GoLeanMakeASCIISet(s)
			if !ok {
				as = [256]bool{}
			}
			return as, ok
		}, 2},
		{"Keep inserting after failure", func(s string) ([256]bool, bool) {
			var as [256]bool
			ok := true
			for i := range len(s) {
				if s[i] >= 128 {
					ok = false
				} else {
					as[s[i]] = true
				}
			}
			return as, ok
		}, 3},
	} {
		bad := 0
		for _, f := range cases {
			var want [256]bool
			for i := range len(f.marked) {
				want[f.marked[i]] = true
			}
			if got, ok := c.fn(f.input); got != want || ok != f.ok {
				bad++
				if c.want == 0 {
					t.Errorf("makeASCIISet(%q) mismatched", f.input)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d fixtures", c.name, bad, len(cases))
		if bad != c.want {
			t.Errorf("%s wanted %d failures", c.name, c.want)
		}
	}
}
