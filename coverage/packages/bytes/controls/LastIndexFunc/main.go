package main

import (
	"bytes"

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
		{"éx🐈", 'x', false, result{index: 2, calls: "🐈x"}},
		{"é\xffz", '\uFFFD', false, result{index: 2, calls: "z\uFFFD"}},
		{"éx", 'z', false, result{index: -1, calls: "xé"}},
		{"xx", 'x', false, result{index: 1, calls: "x"}},
	}
	for _, c := range []struct {
		name string
		fn   search
		want int
	}{
		{"Go", func(s string, p func(rune) bool) int { return bytes.LastIndexFunc([]byte(s), p) }, 0},
		{"Scan forward", strings.IndexFunc, 4},
		{"Return rune ordinal", func(s string, f func(rune) bool) int {
			index := strings.LastIndexFunc(s, f)
			if index < 0 {
				return index
			}
			return len([]rune(s[:index]))
		}, 2},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := observe(c.fn, f.s, f.key, f.nilPred); got != f.want {
				bad++
				if c.want == 0 {
					log.Fatalf("LastIndexFunc(%q) = %+v, want %+v", f.s, got, f.want)
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", c.name, bad, len(fixtures))
		if bad != c.want {
			log.Fatalf("%s wanted %d failures", c.name, c.want)
		}
	}
}
