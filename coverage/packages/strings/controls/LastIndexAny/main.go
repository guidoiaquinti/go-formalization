package main

import (
	"fmt"
	"log"
	"strings"
)

func main() {
	cases := []struct {
		s, chars string
		want     int
	}{
		{"abc", "", -1},
		{"éx", "x", 2},
		{"é\xffa", "\xc2", 2},
		{"\xffé\uFFFD", "\xff", 3},
		{"0123456789x", "zx", 10},
		{"é", "\xa9", -1},
		{"aéβ", "βé", 3},
	}
	for _, impl := range []struct {
		name string
		f    func(string, string) int
		want int
	}{
		{"Go", strings.LastIndexAny, 0},
		{"Compare individual bytes", func(s, chars string) int {
			for i := len(s) - 1; i >= 0; i-- {
				if strings.IndexByte(chars, s[i]) >= 0 {
					return i
				}
			}
			return -1
		}, 4},
		{"Search from the start", strings.IndexAny, 2},
	} {
		bad := 0
		for _, c := range cases {
			if got := impl.f(c.s, c.chars); got != c.want {
				bad++
				if impl.want == 0 {
					log.Fatalf("LastIndexAny(%q,%q) = %d, want %d", c.s, c.chars, got, c.want)
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", impl.name, bad, len(cases))
		if bad != impl.want {
			log.Fatalf("%s wanted %d disagreements", impl.name, impl.want)
		}
	}
}
