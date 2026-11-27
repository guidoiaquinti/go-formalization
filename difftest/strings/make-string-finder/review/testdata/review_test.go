package strings_test

import (
	"slices"
	"strings"
	"testing"
)

func includeLast(pattern string) (string, [256]int, []int) {
	got, bad, good := strings.GoLeanReviewFinder(pattern)
	if pattern != "" {
		bad[pattern[len(pattern)-1]] = 0
	}
	return got, bad, good
}

func keepOnlyPrefixBorders(pattern string) (string, [256]int, []int) {
	got, bad, good := strings.GoLeanReviewFinder(pattern)
	for i := range good {
		tail := pattern[i+1:]
		border := len(tail)
		for !strings.HasPrefix(pattern, tail[len(tail)-border:]) {
			border--
		}
		good[i] = len(pattern) - border + len(tail)
	}
	return got, bad, good
}

func TestGoLeanMakeStringFinderReview(t *testing.T) {
	fixtures := []struct {
		pattern string
		bad     map[byte]int
		good    []int
	}{
		{"", nil, []int{}},
		{"\xff", nil, []int{1}},
		{"aaaa", map[byte]int{'a': 1}, []int{4, 4, 4, 4}},
		{"abc", map[byte]int{'a': 2, 'b': 1}, []int{5, 4, 1}},
		{"mississi", map[byte]int{'m': 7, 'i': 3, 's': 1}, []int{15, 14, 13, 7, 11, 10, 7, 1}},
		{"abcxxxabc", map[byte]int{'a': 2, 'b': 1, 'c': 6, 'x': 3}, []int{14, 13, 12, 11, 10, 9, 11, 10, 1}},
		{"\xff\x00\xff", map[byte]int{0xff: 2, 0: 1}, []int{4, 3, 1}},
	}
	for _, c := range []struct {
		name  string
		build func(string) (string, [256]int, []int)
		want  int
	}{
		{"Go", strings.GoLeanReviewFinder, 0},
		{"Enter final byte in bad-character table", includeLast, 6},
		{"Omit good-suffix repeat overrides", keepOnlyPrefixBorders, 4},
	} {
		mismatches := 0
		for _, f := range fixtures {
			pattern, bad, good := c.build(f.pattern)
			var wantBad [256]int
			for i := range wantBad {
				wantBad[i] = len(f.pattern)
			}
			for b, shift := range f.bad {
				wantBad[b] = shift
			}
			if pattern != f.pattern || bad != wantBad || good == nil || !slices.Equal(good, f.good) {
				mismatches++
				if c.want == 0 {
					t.Errorf("makeStringFinder(%x): pattern=%x bad matches=%t good=%v nil=%t, want good=%v",
						f.pattern, pattern, bad == wantBad, good, good == nil, f.good)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d", c.name, mismatches, len(fixtures))
		if mismatches != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
