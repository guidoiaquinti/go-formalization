package strings_test

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func commonPrefix(a, b string) int {
	for i := range min(len(a), len(b)) {
		if a[i] != b[i] {
			return i
		}
	}
	return min(len(a), len(b))
}

func decodedSuffix(a, b string) int {
	var width int
	for a != "" && b != "" {
		ar, aw := utf8.DecodeLastRuneInString(a)
		br, bw := utf8.DecodeLastRuneInString(b)
		if ar != br || aw != bw {
			break
		}
		width += aw
		a, b = a[:len(a)-aw], b[:len(b)-bw]
	}
	return width
}

func TestGoLeanLongestCommonSuffixReview(t *testing.T) {
	fixtures := []struct {
		a, b string
		want int
	}{
		{"", "aa", 0}, {"aa", "", 0}, {"xabc", "abc", 3},
		{"prefixX", "prefixY", 0}, {"xxbcb", "xyacb", 2}, {"aaaa", "aaaa", 4},
		{"é", "©", 1}, {"\xff", "\xfe", 0}, {"A\xff", "B\xff", 1},
		{"a🚀", "b🚀", 4}, {"\x00X\x00", "X\x00", 2},
	}
	for _, c := range []struct {
		name  string
		count func(string, string) int
		want  int
	}{
		{"Go", strings.GoLeanLongestCommonSuffix, 0},
		{"Compare from start", commonPrefix, 7},
		{"Compare decoded runes", decodedSuffix, 2},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := c.count(f.a, f.b); got != f.want {
				bad++
				if c.want == 0 {
					t.Errorf("longestCommonSuffix(%x,%x) = %d, want %d", f.a, f.b, got, f.want)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d", c.name, bad, len(fixtures))
		if bad != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
