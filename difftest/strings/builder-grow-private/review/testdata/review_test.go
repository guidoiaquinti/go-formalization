package strings_test

import (
	"runtime"
	"slices"
	"strings"
	"testing"
)

type grow func(*strings.Builder, int)

func growStatus(fn grow, b *strings.Builder, n int) (status string) {
	defer func() {
		if p := recover(); p != nil {
			if e, ok := p.(runtime.Error); ok {
				switch {
				case b == nil:
					status = "nil"
				case strings.Contains(e.Error(), "makeslice"):
					status = "allocation"
				case strings.Contains(e.Error(), "slice bounds"):
					status = "bounds"
				default:
					panic(p)
				}
			} else if p == "strings: illegal use of non-zero Builder copied by value" {
				status = "copy"
			} else if p == "strings.Builder.Grow: negative count" {
				status = "negative"
			} else {
				panic(p)
			}
		}
	}()
	fn(b, n)
	return "ok"
}

func TestGoLeanBuilderGrowPrivateReview(t *testing.T) {
	cases := []struct {
		mode, capacity, n, wantCap int
		text, status               string
	}{
		{0, -1, 0, 0, "", "ok"},
		{0, -1, -1, -1, "", "allocation"},
		{0, 8, -16, 8, "abc", "bounds"},
		{1, 8, -17, 8, "abc", "allocation"},
		{2, 8, 1, cap(slices.Grow([]byte(nil), 17)), "\xffx", "ok"},
		{3, -1, 0, -1, "", "nil"},
	}
	for _, impl := range []struct {
		name string
		run  grow
		want int
	}{
		{"Go", strings.GoLeanBuilderGrowPrivate, 0},
		{"Run the public copy and sign checks", (*strings.Builder).Grow, 5},
		{"Skip zero request", func(b *strings.Builder, n int) {
			if n != 0 {
				strings.GoLeanBuilderGrowPrivate(b, n)
			}
		}, 2},
	} {
		bad := 0
		for _, f := range cases {
			b := strings.GoLeanBuilderCopyCheckSetup(f.mode, f.capacity, f.text)
			var prior *strings.Builder
			if b != nil {
				_, prior, _, _ = strings.GoLeanBuilderCopyCheckObserve(b)
			}
			got := growStatus(impl.run, b, f.n)
			pass := got == f.status
			if b != nil {
				_, addr, c, text := strings.GoLeanBuilderCopyCheckObserve(b)
				pass = pass && addr == prior && c == f.wantCap && text == f.text
			}
			if !pass {
				bad++
				if impl.want == 0 {
					t.Errorf("fixture %+v returned %s", f, got)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d fixtures", impl.name, bad, len(cases))
		if bad != impl.want {
			t.Errorf("%s wanted %d failures", impl.name, impl.want)
		}
	}
}
