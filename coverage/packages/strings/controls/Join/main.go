package main

import (
	"fmt"
	"log"
	"strings"
)

func main() {
	tests := []struct {
		parts     []string
		sep, want string
	}{
		{nil, "<>", ""}, {[]string{}, "<>", ""}, {[]string{"\xffx"}, "<>", "\xffx"},
		{[]string{"", "x", ""}, "<>", "<>x<>"}, {[]string{"", "", ""}, "x", "xx"},
		{[]string{"a", "\xff", "b"}, "\x00", "a\x00\xff\x00b"}, {[]string{"a", "b"}, "", "ab"},
	}
	for _, impl := range []struct {
		name string
		join func([]string, string) string
		want int
	}{
		{"Go", strings.Join, 0},
		{"Discard empty elements", func(parts []string, sep string) string {
			var nonempty []string
			for _, part := range parts {
				if part != "" {
					nonempty = append(nonempty, part)
				}
			}
			return strings.Join(nonempty, sep)
		}, 2},
		{"Add separator to the front", func(parts []string, sep string) string {
			if len(parts) == 0 {
				return ""
			}
			return sep + strings.Join(parts, sep)
		}, 4},
	} {
		bad := 0
		for _, f := range tests {
			if got := impl.join(f.parts, f.sep); got != f.want {
				bad++
				if impl.want == 0 {
					log.Fatalf("Join(%q,%q) = %q, want %q", f.parts, f.sep, got, f.want)
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", impl.name, bad, len(tests))
		if bad != impl.want {
			log.Fatalf("%s wanted %d mismatches", impl.name, impl.want)
		}
	}
}
