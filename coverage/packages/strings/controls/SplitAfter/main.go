package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

func main() {
	cases := []struct {
		s, sep string
		want   []string
	}{
		{"", "", []string{}}, {"", ":", []string{""}},
		{"::a::", "::", []string{"::", "a::", ""}},
		{"aaaaa", "aa", []string{"aa", "aa", "a"}},
		{"é\xffx", "", []string{"é", "\xff", "x"}},
		{"a:b:c", ":", []string{"a:", "b:", "c"}},
		{"abc", "|", []string{"abc"}},
	}
	for _, c := range []struct {
		name string
		fn   func(string, string) []string
		want int
	}{
		{"Go", strings.SplitAfter, 0},
		{"Omit separators", strings.Split, 3},
		{"Normalize malformed input", func(s, sep string) []string { return strings.SplitAfter(strings.ToValidUTF8(s, "\ufffd"), sep) }, 1},
		{"Conflate nil and empty", func(s, sep string) []string {
			r := strings.SplitAfter(s, sep)
			if len(r) == 0 {
				return nil
			}
			return r
		}, 1},
	} {
		bad := 0
		for _, f := range cases {
			if got := c.fn(f.s, f.sep); !reflect.DeepEqual(got, f.want) {
				bad++
				if c.want == 0 {
					log.Fatalf("SplitAfter(%q, %q) = %#v, want %#v", f.s, f.sep, got, f.want)
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", c.name, bad, len(cases))
		if bad != c.want {
			log.Fatalf("%s wanted %d disagreements", c.name, c.want)
		}
	}
}
