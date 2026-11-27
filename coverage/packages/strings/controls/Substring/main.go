package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

type split struct {
	before, after string
	found         bool
}

func outcome(before, after string, found bool) split {
	return split{before, after, found}
}

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		log.Fatalf("Go version %s, want %s", runtime.Version(), want)
	}
	cases := []struct {
		text, needle string
		first, last  int
		cut          split
	}{
		{"", "", 0, 0, split{"", "", true}},
		{"abc", "", 0, 3, split{"", "abc", true}},
		{"abc", "abcd", -1, -1, split{"abc", "", false}},
		{"aba", "b", 1, 1, split{"a", "a", true}},
		{"aaaa", "aaa", 0, 1, split{"", "a", true}},
		{"éé", "\xa9", 1, 3, split{"\xc3", "é", true}},
		{"abc", "z", -1, -1, split{"abc", "", false}},
	}
	var indexFailures, lastFailures, containsFailures, cutFailures int
	for _, c := range cases {
		first, last := strings.Index(c.text, c.needle), strings.LastIndex(c.text, c.needle)
		if first != c.first || last != c.last {
			log.Fatalf("search(%q, %q) = (%d, %d), want (%d, %d)", c.text, c.needle, first, last, c.first, c.last)
		}
		if got := strings.Contains(c.text, c.needle); got != c.cut.found {
			log.Fatalf("Contains(%q, %q) = %t, want %t", c.text, c.needle, got, c.cut.found)
		}
		if got := outcome(strings.Cut(c.text, c.needle)); got != c.cut {
			log.Fatalf("Cut(%q, %q) = %+v, want %+v", c.text, c.needle, got, c.cut)
		}
		if last != c.first {
			indexFailures++
		}
		if first != c.last {
			lastFailures++
		}
		if (first > 0) != c.cut.found {
			containsFailures++
		}
		cutMutant := split{found: last >= 0}
		if last >= 0 {
			cutMutant.before = c.text[:last]
			cutMutant.after = c.text[last+len(c.needle):]
		}
		if cutMutant != c.cut {
			cutFailures++
		}
	}
	failures := [4]int{indexFailures, lastFailures, containsFailures, cutFailures}
	if want := [4]int{3, 3, 3, 5}; failures != want {
		log.Fatalf("substring mutants failed %v, want %v", failures, want)
	}
	fmt.Printf("%s %s/%s: %d substring fixtures passed, mutant failures %v\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, len(cases), failures)
}
