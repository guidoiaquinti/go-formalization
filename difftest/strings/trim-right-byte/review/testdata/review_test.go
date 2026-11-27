package strings_test

import (
	"strings"
	"testing"
)

func stripBothEnds(s string, c byte) string {
	for len(s) > 0 && s[0] == c {
		s = s[1:]
	}
	for len(s) > 0 && s[len(s)-1] == c {
		s = s[:len(s)-1]
	}
	return s
}

func TestGoLeanTrimRightByteReview(t *testing.T) {
	fixtures := []struct {
		input string
		cut   byte
		want  string
	}{
		{"", 0xff, ""},
		{"\x00X\x00\x00", 0, "\x00X"},
		{"aaaa", 'a', ""},
		{"aaXaa", 'a', "aaX"},
		{"aaaX", 'a', "aaaX"},
		{"é", 0xa9, "\xc3"},
		{"é", 0xc3, "é"},
		{"\xffx\xff\xff", 0xff, "\xffx"},
		{"\x80�\x80", 0x80, "\x80�"},
		{"xÿ", 0xff, "xÿ"},
	}
	for _, c := range []struct {
		name string
		trim func(string, byte) string
		want int
	}{
		{"Go", strings.GoLeanTrimRightByte, 0},
		{"Trim both ends", stripBothEnds, 6},
		{"Compare decoded runes", func(s string, c byte) string {
			return strings.TrimRightFunc(s, func(r rune) bool { return r == rune(c) })
		}, 4},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := c.trim(f.input, f.cut); got != f.want {
				bad++
				if c.want == 0 {
					t.Errorf("trimRightByte(%x,%02x) = %x, want %x", f.input, f.cut, got, f.want)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d", c.name, bad, len(fixtures))
		if bad != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
