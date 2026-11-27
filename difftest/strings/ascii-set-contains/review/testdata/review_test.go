package strings_test

import (
	"runtime"
	"strings"
	"testing"
)

type lookup func(*[256]bool, byte) bool
type result struct{ value, panicked bool }

func observe(fn lookup, table *[256]bool, b byte) (out result) {
	defer func() {
		if p := recover(); p != nil {
			if _, ok := p.(runtime.Error); !ok || table != nil {
				panic(p)
			}
			out.panicked = true
		}
	}()
	out.value = fn(table, b)
	return
}
func TestGoLeanASCIISetContainsReview(t *testing.T) {
	table := &[256]bool{127: true, 128: true, 255: true}
	cases := []struct {
		table *[256]bool
		b     byte
		want  result
	}{
		{table, 0, result{}}, {table, 127, result{value: true}}, {table, 128, result{value: true}},
		{table, 255, result{value: true}}, {table, 254, result{}}, {nil, 128, result{panicked: true}},
	}
	for _, impl := range []struct {
		name string
		fn   lookup
		want int
	}{
		{"Go", strings.GoLeanASCIISetContains, 0},
		{"Mask the input to seven bits", func(as *[256]bool, b byte) bool { return strings.GoLeanASCIISetContains(as, b&127) }, 1},
		{"Treat nil as an empty table", func(as *[256]bool, b byte) bool { return as != nil && strings.GoLeanASCIISetContains(as, b) }, 1},
	} {
		bad := 0
		for _, f := range cases {
			if got := observe(impl.fn, f.table, f.b); got != f.want {
				bad++
				if impl.want == 0 {
					t.Errorf("lookup %d got %+v, want %+v", f.b, got, f.want)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d fixtures", impl.name, bad, len(cases))
		if bad != impl.want {
			t.Errorf("%s wanted %d failures", impl.name, impl.want)
		}
	}
}
