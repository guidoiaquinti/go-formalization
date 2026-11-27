package main

import (
	"fmt"
	"log"
	"strings"
)

func main() {
	fixtures := []struct {
		s, chars string
		want     bool
	}{
		{"", "x", false}, {"abc", "", false}, {"éx", "x", true}, {"é", "\xa9", false},
		{"é\uFFFD", "\xff", true}, {"é\xff", "\uFFFD", true}, {"0123456789x", "azx", true},
	}
	for _, c := range []struct {
		name string
		run  func(string, string) bool
		want int
	}{
		{"Go", strings.ContainsAny, 0},
		{"Compare individual bytes", func(s, chars string) bool {
			for i := range len(s) {
				if strings.IndexByte(chars, s[i]) >= 0 {
					return true
				}
			}
			return false
		}, 3},
		{"Skip RuneError", func(s, chars string) bool {
			for _, r := range s {
				if r != '\uFFFD' && strings.ContainsRune(chars, r) {
					return true
				}
			}
			return false
		}, 2},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := c.run(f.s, f.chars); got != f.want {
				bad++
				if c.want == 0 {
					log.Fatalf("ContainsAny(%q,%q) = %t, want %t", f.s, f.chars, got, f.want)
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", c.name, bad, len(fixtures))
		if bad != c.want {
			log.Fatalf("%s wanted %d failures", c.name, c.want)
		}
	}
}
