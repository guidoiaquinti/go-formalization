package main

import (
	"fmt"
	"log"
	"math"
	"strings"
)

type repeat func(string, int) string
type result struct{ bytes, panicText string }

func observe(fn repeat, s string, n int) (out result) {
	defer func() {
		if value := recover(); value != nil {
			out.panicText = value.(string)
		}
	}()
	out.bytes = fn(s, n)
	return
}

func main() {
	const negative = "strings: negative Repeat count"
	const overflow = "strings: Repeat output length overflow"
	fixtures := []struct {
		s    string
		n    int
		want result
	}{
		{"", math.MaxInt, result{}}, {"x", 0, result{}}, {"\xffx", 1, result{bytes: "\xffx"}},
		{"\xff\x00", 3, result{bytes: "\xff\x00\xff\x00\xff\x00"}},
		{"", math.MinInt, result{panicText: negative}}, {"abc", math.MaxInt, result{panicText: overflow}},
		{"ab", math.MaxInt, result{panicText: overflow}}, {"=", 65, result{bytes: strings.Join(make([]string, 66), "=")}},
		{"abc", 2731, result{bytes: strings.Join(make([]string, 2732), "abc")}},
	}
	for _, c := range []struct {
		name string
		fn   repeat
		want int
	}{
		{"Go", strings.Repeat, 0},
		{"Allow negative repeat of empty", func(s string, n int) string {
			if len(s) == 0 {
				return ""
			}
			return strings.Repeat(s, n)
		}, 1},
		{"Normalize bytes", func(s string, n int) string { return strings.Repeat(strings.ToValidUTF8(s, "\ufffd"), n) }, 2},
		{"Stop at the chunk size", func(s string, n int) string { v := strings.Repeat(s, n); return v[:min(len(v), 8192)] }, 1},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := observe(c.fn, f.s, f.n); got != f.want {
				bad++
				if c.want == 0 {
					log.Fatalf("Repeat(%q,%d) gave %d bytes, panic %q", f.s, f.n, len(got.bytes), got.panicText)
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", c.name, bad, len(fixtures))
		if bad != c.want {
			log.Fatalf("%s wanted %d disagreements", c.name, c.want)
		}
	}
}
