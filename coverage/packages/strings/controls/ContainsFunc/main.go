package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

type search func(string, func(rune) bool) bool
type result struct {
	found  bool
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
	out.found = fn(s, predicate)
	return
}
func main() {
	fixtures := []struct {
		s       string
		key     rune
		nilPred bool
		want    result
	}{
		{"", 'x', true, result{found: false}}, {"x", 'x', true, result{panics: true}},
		{"éx🐈", 'x', false, result{found: true, calls: "éx"}},
		{"é\xffz", '\uFFFD', false, result{found: true, calls: "é\uFFFD"}},
		{"éx", 'z', false, result{found: false, calls: "éx"}},
		{"xx", 'x', false, result{found: true, calls: "x"}},
	}
	for _, c := range []struct {
		name string
		fn   search
		want int
	}{
		{"Go", strings.ContainsFunc, 0},
		{"Skip RuneError", func(s string, f func(rune) bool) bool {
			return strings.ContainsFunc(s, func(r rune) bool { return r != '\uFFFD' && f(r) })
		}, 1},
		{"Continue after matching", func(s string, f func(rune) bool) bool {
			matched := false
			for _, r := range s {
				if f(r) {
					matched = true
				}
			}
			return matched
		}, 3},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := observe(c.fn, f.s, f.key, f.nilPred); got != f.want {
				bad++
				if c.want == 0 {
					log.Fatalf("ContainsFunc(%q) = %+v, want %+v", f.s, got, f.want)
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", c.name, bad, len(fixtures))
		if bad != c.want {
			log.Fatalf("%s wanted %d failures", c.name, c.want)
		}
	}
}
