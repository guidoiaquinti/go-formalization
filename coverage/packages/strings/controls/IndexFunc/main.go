package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

type search func(string, func(rune) bool) int
type result struct {
	index  int
	calls  string
	panics bool
}

func observe(fn search, s string, key rune, useNil bool) (out result) {
	var calls strings.Builder
	defer func() {
		out.calls = calls.String()
		if p := recover(); p != nil {
			if _, ok := p.(runtime.Error); !ok {
				panic(p)
			}
			out.panics = true
		}
	}()
	var predicate func(rune) bool
	if !useNil {
		predicate = func(r rune) bool { calls.WriteRune(r); return r == key }
	}
	out.index = fn(s, predicate)
	return
}
func main() {
	fixtures := []struct {
		s       string
		key     rune
		nilPred bool
		want    result
	}{
		{"", 'x', true, result{index: -1}}, {"x", 'x', true, result{panics: true}},
		{"éx🐈", 'x', false, result{index: 2, calls: "éx"}},
		{"é\xffz", '\uFFFD', false, result{index: 2, calls: "é\uFFFD"}},
		{"éx", 'z', false, result{index: -1, calls: "éx"}},
		{"xx", 'x', false, result{index: 0, calls: "x"}},
	}
	for _, c := range []struct {
		name string
		fn   search
		want int
	}{
		{"Go", strings.IndexFunc, 0},
		{"Return rune ordinal", func(s string, f func(rune) bool) int {
			i := 0
			for _, r := range s {
				if f(r) {
					return i
				}
				i++
			}
			return -1
		}, 2},
		{"Continue after matching", func(s string, f func(rune) bool) int {
			first := -1
			for i, r := range s {
				if f(r) && first < 0 {
					first = i
				}
			}
			return first
		}, 3},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := observe(c.fn, f.s, f.key, f.nilPred); got != f.want {
				bad++
				if c.want == 0 {
					log.Fatalf("IndexFunc(%q) = %+v, want %+v", f.s, got, f.want)
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", c.name, bad, len(fixtures))
		if bad != c.want {
			log.Fatalf("%s wanted %d failures", c.name, c.want)
		}
	}
}
