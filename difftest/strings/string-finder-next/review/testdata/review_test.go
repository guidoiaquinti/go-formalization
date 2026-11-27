package strings_test

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func runeIndex(pattern, text string) int {
	i := strings.GoLeanStringFinderNext(pattern, text)
	if i < 0 {
		return i
	}
	return utf8.RuneCountInString(text[:i])
}

func skipMatchedSuffix(pattern, text string) int {
	bad, good := strings.DumpTables(pattern)
	for end := len(pattern) - 1; end < len(text); {
		i, j := end, len(pattern)-1
		for j >= 0 && text[i] == pattern[j] {
			i--
			j--
		}
		if j < 0 {
			return i + 1
		}
		end += max(bad[text[i]], good[j])
	}
	return -1
}

func TestGoLeanStringFinderNextReview(t *testing.T) {
	fixtures := []struct {
		pattern, text string
		want          int
	}{
		{"", "", 0}, {"", "aa", 0}, {"abc", "ab", -1},
		{"aaa", "baaa", 1}, {"aaa", "aaaaaa", 0}, {"ab", "zabab", 1},
		{"🚀", "é🚀é🚀", 2}, {"\xa9", "é", 1},
		{"\xff\x00", "é\xff\x00", 2}, {"\xff", "\xfe", -1},
		{"abcxxxabc", "abcxxxabqabcxxxabc", 9},
	}
	for _, c := range []struct {
		name   string
		search func(string, string) int
		want   int
	}{
		{"Go", strings.GoLeanStringFinderNext, 0},
		{"Return last match", func(pattern, text string) int { return strings.LastIndex(text, pattern) }, 4},
		{"Return a rune offset", runeIndex, 2},
		{"Shift from window end after matching a suffix", skipMatchedSuffix, 1},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := c.search(f.pattern, f.text); got != f.want {
				bad++
				if c.want == 0 {
					t.Errorf("next(%x,%x) = %d, want %d", f.pattern, f.text, got, f.want)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d", c.name, bad, len(fixtures))
		if bad != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
