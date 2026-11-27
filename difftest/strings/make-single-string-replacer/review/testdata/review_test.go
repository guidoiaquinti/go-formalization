package strings_test

import (
	"slices"
	"strings"
	"testing"
)

type constructor func(string, string) (string, string, [256]int, []int)

func swap(pattern, value string) (string, string, [256]int, []int) {
	return strings.GoLeanReviewSingleMaker(value, pattern)
}

func normalizeValue(pattern, value string) (string, string, [256]int, []int) {
	return strings.GoLeanReviewSingleMaker(pattern, strings.ToValidUTF8(value, "�"))
}

func TestGoLeanSingleMakerReview(t *testing.T) {
	fixtures := [][2]string{
		{"", ""}, {"", "\xff"}, {"\x00\xff", ""},
		{"mississi", "abcxxxabc"}, {"abcxxxabc", "Z"},
		{"🚀é", "\xff\x00�"}, {"\x00x", "\x00x"},
	}
	for _, c := range []struct {
		name string
		make constructor
		want int
	}{
		{"Go", strings.GoLeanReviewSingleMaker, 0},
		{"Swap constructor operands", swap, 5},
		{"Normalize the replacement", normalizeValue, 2},
	} {
		bad := 0
		for _, f := range fixtures {
			pattern, value, badTable, goodTable := c.make(f[0], f[1])
			wantBad, wantGood := strings.DumpTables(f[0])
			if pattern != f[0] || value != f[1] || goodTable == nil ||
				!slices.Equal(badTable[:], wantBad) || !slices.Equal(goodTable, wantGood) {
				bad++
				if c.want == 0 {
					t.Errorf("makeSingleStringReplacer(%x,%x) got pattern=%x value=%x good=%v",
						f[0], f[1], pattern, value, goodTable)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d", c.name, bad, len(fixtures))
		if bad != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
