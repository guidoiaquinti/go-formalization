package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
)

type mapper func(func(rune) rune, []byte) []byte

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
		{"Go", bytes.Map, 0},
		{"Call mapper again on first rune", func(fn func(rune) rune, s []byte) []byte {
			for _, r := range string(s) {
				fn(r)
				break
			}
			return bytes.Map(fn, s)
		}, 4},
		{"Drop RuneError inputs", func(fn func(rune) rune, s []byte) []byte {
			return bytes.Map(func(r rune) rune {
				v := fn(r)
				if r == '\uFFFD' {
					return -1
				}
				return v
			}, s)
		}, 1},
		{"Return nil for empty output", func(fn func(rune) rune, s []byte) []byte {
			result := bytes.Map(fn, s)
			if len(result) == 0 {
				return nil
			}
			return result
		}, 2},
	} {
		bad := 0
		for _, f := range cases {
			var calls strings.Builder
			got := impl.run(func(r rune) rune { calls.WriteRune(r); return f.fn(r) }, []byte(f.s))
			if got == nil || string(got) != f.want || calls.String() != f.calls {
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
