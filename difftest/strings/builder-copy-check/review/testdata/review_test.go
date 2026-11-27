package strings_test

import (
	"runtime"
	"strings"
	"testing"
)

func status(fn func(*strings.Builder), b *strings.Builder) (out string) {
	defer func() {
		switch p := recover().(type) {
		case nil:
		case runtime.Error:
			if b != nil {
				panic(p)
			}
			out = "nil"
		case string:
			if p != "strings: illegal use of non-zero Builder copied by value" {
				panic(p)
			}
			out = "copy"
		default:
			panic(p)
		}
	}()
	fn(b)
	return "ok"
}

func TestGoLeanBuilderCopyCheckReview(t *testing.T) {
	cases := []struct {
		mode, capacity   int
		contents, status string
	}{
		{0, -1, "", "ok"}, {0, 0, "", "ok"}, {0, 17, "\xff\x00", "ok"},
		{1, 11, "xyz", "ok"}, {2, -1, "", "copy"}, {2, 9, "x", "copy"},
		{3, -1, "", "nil"},
	}
	for _, implementation := range []struct {
		name string
		fn   func(*strings.Builder)
		want int
	}{
		{"Go", strings.GoLeanBuilderCopyCheck, 0},
		{"Skip checking empty buffers", func(b *strings.Builder) {
			if b != nil && b.Len() == 0 {
				return
			}
			strings.GoLeanBuilderCopyCheck(b)
		}, 3},
		{"Suppress foreign identity panic", func(b *strings.Builder) {
			if b != nil {
				identity, _, _, _ := strings.GoLeanBuilderCopyCheckObserve(b)
				if identity == "foreign" {
					return
				}
			}
			strings.GoLeanBuilderCopyCheck(b)
		}, 2},
	} {
		bad := 0
		for _, f := range cases {
			b := strings.GoLeanBuilderCopyCheckSetup(f.mode, f.capacity, f.contents)
			var old *strings.Builder
			if b != nil {
				_, old, _, _ = strings.GoLeanBuilderCopyCheckObserve(b)
			}
			pass := status(implementation.fn, b) == f.status
			if b != nil {
				identity, addr, capacity, contents := strings.GoLeanBuilderCopyCheckObserve(b)
				if f.mode == 2 {
					pass = pass && identity == "foreign" && addr == old
				} else {
					pass = pass && identity == "self" && addr == b
				}
				pass = pass && capacity == f.capacity && contents == f.contents
			}
			if !pass {
				bad++
				if implementation.want == 0 {
					t.Errorf("case %+v failed", f)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d fixtures", implementation.name, bad, len(cases))
		if bad != implementation.want {
			t.Errorf("%s wanted %d disagreements", implementation.name, implementation.want)
		}
	}
}
