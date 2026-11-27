package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

func main() {
	tests := []struct {
		s, sep string
		n      int
		want   []string
	}{
		{"a:b", ":", 0, nil}, {"", "", 1, []string{}}, {"", ":", -1, []string{""}},
		{"aaa", "aa", 1, []string{"aaa"}}, {"aaa", "aa", -3, []string{"aa", "a"}},
		{"a:b:c", ":", 2, []string{"a:", "b:c"}}, {"é\xffx", "", 2, []string{"é", "\xffx"}},
		{"::", "::", 1, []string{"::"}},
	}
	for _, impl := range []struct {
		name string
		run  func(string, string, int) []string
		want int
	}{
		{"Go", strings.SplitAfterN, 0},
		{"Drop separators", strings.SplitN, 2},
		{"Ignore limit", func(s, sep string, n int) []string { return strings.SplitAfter(s, sep) }, 5},
		{"Collapse empty to nil", func(s, sep string, n int) []string {
			got := strings.SplitAfterN(s, sep, n)
			if len(got) == 0 {
				return nil
			}
			return got
		}, 1},
	} {
		bad := 0
		for _, test := range tests {
			got := impl.run(test.s, test.sep, test.n)
			if !reflect.DeepEqual(got, test.want) {
				bad++
				if impl.want == 0 {
					log.Fatalf("SplitAfterN(%q,%q,%d) = %#v, want %#v", test.s, test.sep, test.n, got, test.want)
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", impl.name, bad, len(tests))
		if bad != impl.want {
			log.Fatalf("%s wanted %d disagreements", impl.name, impl.want)
		}
	}
}
