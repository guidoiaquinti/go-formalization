package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

type trim func(string, func(rune) bool) string
type outcome struct {
	text, calls string
	panics      bool
}

func observe(fn trim, s string, missing bool) (out outcome) {
	var calls strings.Builder
	defer func() {
		out.calls = calls.String()
		if p := recover(); p != nil {
			if _, ok := p.(runtime.Error); !ok || !missing {
				panic(p)
			}
			out.panics = true
		}
	}()
	var f func(rune) bool
	if !missing {
		f = func(r rune) bool { calls.WriteRune(r); return r == ' ' || r == '\uFFFD' }
	}
	out.text = fn(s, f)
	return
}
func main() {
	fixtures := []struct {
		s       string
		missing bool
		want    outcome
	}{
		{"", true, outcome{}}, {"a", true, outcome{panics: true}},
		{"   ", false, outcome{calls: "   "}},
		{" a  ", false, outcome{text: " a", calls: "  a"}},
		{"\xffé\xff", false, outcome{text: "\xffé", calls: "\uFFFDé"}},
		{"a b", false, outcome{text: "a b", calls: "b"}},
		{"a\xffb", false, outcome{text: "a\xffb", calls: "b"}},
		{"a", false, outcome{text: "a", calls: "a"}},
	}
	for _, c := range []struct {
		name string
		run  trim
		want int
	}{
		{"Go", strings.TrimRightFunc, 0},
		{"Trim from the left", strings.TrimLeftFunc, 4},
		{"Normalize malformed bytes before returning", func(s string, f func(rune) bool) string {
			return strings.TrimRightFunc(strings.ToValidUTF8(s, "\ufffd"), f)
		}, 2},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := observe(c.run, f.s, f.missing); got != f.want {
				bad++
				if c.want == 0 {
					log.Fatalf("%s(%q) = %+v, want %+v", c.name, f.s, got, f.want)
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", c.name, bad, len(fixtures))
		if bad != c.want {
			log.Fatalf("%s wanted %d disagreements", c.name, c.want)
		}
	}
}
