package strings_test

import (
	"strings"
	"testing"
)

func withoutTail(pattern, value, input string) string {
	var out strings.Builder
	matched := false
	for {
		i := strings.Index(input, pattern)
		if i < 0 {
			if !matched {
				return input
			}
			return out.String()
		}
		matched = true
		out.WriteString(input[:i])
		out.WriteString(value)
		input = input[i+len(pattern):]
	}
}

func TestGoLeanSingleReplacerReview(t *testing.T) {
	fixtures := []struct{ pattern, value, input, want string }{
		{"a", "A", "", ""},
		{"z", "Z", "a\xff🚀", "a\xff🚀"},
		{"aa", "X", "aaaaa", "XXa"},
		{"aba", "Z", "ababa", "Zba"},
		{"a", "aa", "aaa", "aaaaaa"},
		{"\x00", "", "z\x00x\x00y", "zxy"},
		{"\xa9", "\xff", "ééx", "\xc3\xff\xc3\xffx"},
		{"🚀", "x", "é🚀x🚀", "éxxx"},
	}
	for _, c := range []struct {
		name    string
		replace func(string, string, string) string
		want    int
	}{
		{"Go", strings.GoLeanReviewSingleReplacer, 0},
		{"Replace only first match", func(pattern, value, input string) string {
			return strings.Replace(input, pattern, value, 1)
		}, 5},
		{"Omit unmatched tail after replacement", withoutTail, 4},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := c.replace(f.pattern, f.value, f.input); got != f.want {
				bad++
				if c.want == 0 {
					t.Errorf("Replace(%x,%x,%x) = %x, want %x", f.pattern, f.value, f.input, got, f.want)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d", c.name, bad, len(fixtures))
		if bad != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
