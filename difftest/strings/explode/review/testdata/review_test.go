package strings_test

import (
	"math"
	"reflect"
	"strings"
	"testing"
)

func TestGoLeanExplodeReview(t *testing.T) {
	fixtures := []struct {
		s    string
		n    int
		want []string
	}{
		{"", math.MinInt, []string{}},
		{"a", 0, []string{}},
		{"é\xffx", 1, []string{"é\xffx"}},
		{"é\xffx", 2, []string{"é", "\xffx"}},
		{"é\xffx", math.MaxInt, []string{"é", "\xff", "x"}},
		{"\xe2\x82x", -8, []string{"\xe2", "\x82", "x"}},
	}
	for _, candidate := range []struct {
		name string
		fn   func(string, int) []string
		want int
	}{
		{"Go", strings.GoLeanExplode, 0},
		{"Discard bytes after the budget", func(s string, n int) []string {
			full := strings.GoLeanExplode(s, -1)
			if n >= 0 && n < len(full) {
				return full[:n]
			}
			return full
		}, 2},
		{"Convert fields through runes", func(s string, n int) []string {
			fields := strings.GoLeanExplode(s, n)
			for i, field := range fields {
				fields[i] = string([]rune(field))
			}
			return fields
		}, 4},
		{"Return nil at zero", func(s string, n int) []string {
			if n == 0 {
				return nil
			}
			return strings.GoLeanExplode(s, n)
		}, 1},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := candidate.fn(f.s, f.n); !reflect.DeepEqual(got, f.want) {
				bad++
				if candidate.want == 0 {
					t.Errorf("explode(%q,%d) = %#v, want %#v", f.s, f.n, got, f.want)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d fixtures", candidate.name, bad, len(fixtures))
		if bad != candidate.want {
			t.Errorf("%s wanted %d disagreements", candidate.name, candidate.want)
		}
	}
}
