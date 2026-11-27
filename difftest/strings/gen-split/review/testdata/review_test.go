package strings_test

import (
	"math"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type split func(string, string, int, int) []string

func observe(fn split, s, sep string, save, limit int) (parts []string, bounds bool) {
	defer func() {
		if p := recover(); p != nil {
			if err, ok := p.(runtime.Error); ok && strings.Contains(err.Error(), "slice bounds") {
				bounds = true
				return
			}
			panic(p)
		}
	}()
	return fn(s, sep, save, limit), false
}

func TestGoLeanGenSplitReview(t *testing.T) {
	tests := []struct {
		s, sep      string
		save, limit int
		parts       []string
		bounds      bool
	}{
		{"ab:c", ":", -1, 2, []string{"a", "c"}, false},
		{"ab:c", ":", 2, 2, []string{"ab:c", "c"}, false},
		{"ab:c", ":", -3, 2, nil, true},
		{"ab:c", ":", math.MaxInt, 2, nil, true},
		{"ab:c", ":", math.MaxInt, 1, []string{"ab:c"}, false},
		{"no match", ":", -5, 3, []string{"no match"}, false},
		{"é\xffx", "", math.MinInt, 2, []string{"é", "\xffx"}, false},
		{"ab:c", ":", math.MaxInt, 0, nil, false},
	}
	for _, impl := range []struct {
		name string
		run  split
		want int
	}{
		{"Go", strings.GoLeanGenSplit, 0},
		{"Clamp retained width", func(s, sep string, save, limit int) []string {
			return strings.GoLeanGenSplit(s, sep, min(max(save, 0), len(sep)), limit)
		}, 4},
		{"Process a match at limit one", func(s, sep string, save, limit int) []string {
			if limit == 1 {
				limit = 2
			}
			return strings.GoLeanGenSplit(s, sep, save, limit)
		}, 1},
	} {
		bad := 0
		for _, f := range tests {
			parts, bounds := observe(impl.run, f.s, f.sep, f.save, f.limit)
			if bounds != f.bounds || !reflect.DeepEqual(parts, f.parts) {
				bad++
				if impl.want == 0 {
					t.Errorf("%s(%q,%q,%d,%d) = %q, bounds=%t", impl.name, f.s, f.sep, f.save, f.limit, parts, bounds)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d fixtures", impl.name, bad, len(tests))
		if bad != impl.want {
			t.Errorf("%s wanted %d mismatches", impl.name, impl.want)
		}
	}
}
