package strings_test

import (
	"runtime"
	"strings"
	"testing"
)

type search func(string, func(rune) bool, bool) int
type outcome struct {
	index    int
	calls    string
	panicked bool
}

func observe(fn search, s string, key rune, truth, missing bool) (out outcome) {
	var calls strings.Builder
	defer func() {
		out.calls = calls.String()
		if p := recover(); p != nil {
			if _, ok := p.(runtime.Error); !ok || !missing {
				panic(p)
			}
			out.panicked = true
		}
	}()
	var fnArg func(rune) bool
	if !missing {
		fnArg = func(r rune) bool { calls.WriteRune(r); return r == key }
	}
	out.index = fn(s, fnArg, truth)
	return
}
func TestGoLeanLastIndexFuncReview(t *testing.T) {
	fixtures := []struct {
		s              string
		key            rune
		truth, missing bool
		want           outcome
	}{
		{"éxé", 'é', false, false, outcome{index: 2, calls: "éx"}},
		{"é\xffx", '\uFFFD', true, false, outcome{index: 2, calls: "x\uFFFD"}},
		{"éé", 'é', false, false, outcome{index: -1, calls: "éé"}},
		{"x", 'x', true, false, outcome{index: 0, calls: "x"}},
		{"", 'x', false, true, outcome{index: -1}},
		{"x", 'x', false, true, outcome{panicked: true}},
	}
	for _, impl := range []struct {
		name string
		run  search
		want int
	}{
		{"Go", strings.GoLeanLastIndexFunc, 0},
		{"Ignore truth", func(s string, f func(rune) bool, _ bool) int { return strings.GoLeanLastIndexFunc(s, f, true) }, 2},
		{"Count rune ordinals", func(s string, f func(rune) bool, truth bool) int {
			index := strings.GoLeanLastIndexFunc(s, f, truth)
			if index < 0 {
				return index
			}
			return len([]rune(s[:index]))
		}, 2},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := observe(impl.run, f.s, f.key, f.truth, f.missing); got != f.want {
				bad++
				if impl.want == 0 {
					t.Errorf("fixture %+v got %+v", f, got)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d fixtures", impl.name, bad, len(fixtures))
		if bad != impl.want {
			t.Errorf("%s wanted %d failures", impl.name, impl.want)
		}
	}
}
