package strings_test

import (
	"math"
	"strings"
	"testing"
	"unicode"
)

func TestGoLeanIsSeparatorReview(t *testing.T) {
	fixtures := []struct {
		r    rune
		want bool
	}{
		{math.MinInt32, true}, {-1, true}, {0, true},
		{'0', false}, {'Z', false}, {'_', false}, {'-', true}, {0x7f, true},
		{0x80, false}, {'é', false}, {'٣', false}, {'。', false}, {'🚀', false},
		{0x85, true}, {0xa0, true}, {0x2003, true}, {0x200b, false},
		{0xfffd, false}, {math.MaxInt32, false},
	}
	for _, c := range []struct {
		name     string
		classify func(rune) bool
		want     int
	}{
		{"Go", strings.GoLeanIsSeparator, 0},
		{"Classify Unicode punctuation as separators", func(r rune) bool {
			return strings.GoLeanIsSeparator(r) || (r > 127 && unicode.IsPunct(r))
		}, 1},
		{"Treat underscore as a separator", func(r rune) bool {
			return r == '_' || strings.GoLeanIsSeparator(r)
		}, 1},
		{"Reject inputs outside Unicode", func(r rune) bool {
			return r >= 0 && r <= unicode.MaxRune && strings.GoLeanIsSeparator(r)
		}, 2},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := c.classify(f.r); got != f.want {
				bad++
				if c.want == 0 {
					t.Errorf("isSeparator(%d) = %t, want %t", f.r, got, f.want)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d", c.name, bad, len(fixtures))
		if bad != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
