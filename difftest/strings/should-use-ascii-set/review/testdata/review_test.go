package strings_test

import (
	"math"
	"strings"
	"testing"
)

func TestGoLeanShouldUseASCIISetReview(t *testing.T) {
	fixtures := []struct {
		length int
		want   bool
	}{
		{math.MinInt, false}, {math.MinInt + 1, false}, {-9, false}, {-1, false},
		{0, false}, {1, false}, {7, false}, {8, false},
		{9, true}, {10, true}, {math.MaxInt - 1, true}, {math.MaxInt, true},
	}
	for _, c := range []struct {
		name   string
		useSet func(int) bool
		want   int
	}{
		{"Go", strings.GoLeanShouldUseASCIISet, 0},
		{"Inclusive threshold at eight", func(n int) bool { return n >= 8 }, 1},
		{"Unsigned conversion", func(n int) bool { return uint(n) > 8 }, 4},
		{"Add before comparing", func(n int) bool { return n+1 > 9 }, 1},
	} {
		mismatches := 0
		for _, f := range fixtures {
			if got := c.useSet(f.length); got != f.want {
				mismatches++
				if c.want == 0 {
					t.Errorf("shouldUseASCIISet(%d) = %t, want %t", f.length, got, f.want)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d", c.name, mismatches, len(fixtures))
		if mismatches != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
