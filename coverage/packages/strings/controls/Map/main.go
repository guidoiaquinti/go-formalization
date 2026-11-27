package main

import (
	"fmt"
	"log"
	"strings"
)

type mapper func(func(rune) rune, string) string

func main() {
	cases := []struct {
		s, calls, want string
		fn             func(rune) rune
	}{
		{"", "", "", func(r rune) rune { return r }},
		{"a\xff\uFFFD", "a\uFFFD\uFFFD", "a\uFFFD\uFFFD", func(r rune) rune { return r }},
		{"Aéx", "Aéx", "A🐈x", func(r rune) rune {
			if r == 'é' {
				return '🐈'
			}
			return r
		}},
		{"a\x00🐈", "a\x00🐈", "", func(rune) rune { return -1 }},
		{"xy", "xy", "\uFFFD\uFFFD", func(r rune) rune {
			if r == 'x' {
				return 0xd800
			}
			return 0x110000
		}},
	}
	for _, impl := range []struct {
		name string
		run  mapper
		want int
	}{
		{"Go", strings.Map, 0},
		{"Call mapper again on first rune", func(fn func(rune) rune, s string) string {
			for _, r := range s {
				fn(r)
				break
			}
			return strings.Map(fn, s)
		}, 4},
		{"Drop RuneError inputs", func(fn func(rune) rune, s string) string {
			return strings.Map(func(r rune) rune {
				v := fn(r)
				if r == '\uFFFD' {
					return -1
				}
				return v
			}, s)
		}, 1},
	} {
		bad := 0
		for _, f := range cases {
			var calls strings.Builder
			got := impl.run(func(r rune) rune { calls.WriteRune(r); return f.fn(r) }, f.s)
			if got != f.want || calls.String() != f.calls {
				bad++
				if impl.want == 0 {
					log.Fatalf("Map(%q) yielded %q with %q", f.s, got, calls.String())
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", impl.name, bad, len(cases))
		if bad != impl.want {
			log.Fatalf("%s wanted %d failures", impl.name, impl.want)
		}
	}
}
