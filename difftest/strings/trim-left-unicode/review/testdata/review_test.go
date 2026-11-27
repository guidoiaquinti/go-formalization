package strings_test

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func stepOneByte(input, cutset string) string {
	for len(input) != 0 {
		r, _ := utf8.DecodeRuneInString(input)
		if !strings.ContainsRune(cutset, r) {
			break
		}
		input = input[1:]
	}
	return input
}

func TestGoLeanTrimLeftUnicodeReview(t *testing.T) {
	fixtures := []struct{ input, cutset, want string }{
		{"", "�", ""},
		{"\xffé", "", "\xffé"},
		{"ééXé", "é", "Xé"},
		{"\x80🚀X\xff", "�🚀", "X\xff"},
		{"�\xffX", "\x80", "X"},
		{"é", "\xa9", "é"},
		{"\xff\xffX", "�", "X"},
		{"\xe2\x82X", "�", "X"},
		{"🚀", "🚀", ""},
		{"aaaa", "a", ""},
	}
	for _, c := range []struct {
		name string
		trim func(string, string) string
		want int
	}{
		{"Go", strings.GoLeanTrimLeftUnicode, 0},
		{"Consume one byte per rune", stepOneByte, 2},
		{"Normalize malformed input", func(input, cutset string) string {
			return strings.GoLeanTrimLeftUnicode(strings.ToValidUTF8(input, "�"), cutset)
		}, 2},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := c.trim(f.input, f.cutset); got != f.want {
				bad++
				if c.want == 0 {
					t.Errorf("trimLeftUnicode(%x,%x) = %x, want %x", f.input, f.cutset, got, f.want)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d", c.name, bad, len(fixtures))
		if bad != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
