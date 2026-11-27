package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

type cut struct {
	rest  string
	found bool
}

func outcome(rest string, found bool) cut {
	return cut{rest, found}
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
		text, affix string
		prefix      cut
		suffix      cut
	}{
		{"", "", cut{"", true}, cut{"", true}},
		{"abc", "", cut{"abc", true}, cut{"abc", true}},
		{"abc", "abc", cut{"", true}, cut{"", true}},
		{"ababa", "ab", cut{"aba", true}, cut{"ababa", false}},
		{"ababa", "ba", cut{"ababa", false}, cut{"aba", true}},
		{"abc", "b", cut{"abc", false}, cut{"abc", false}},
		{"x", "xy", cut{"x", false}, cut{"x", false}},
		{"\xffa\x00", "\xff", cut{"a\x00", true}, cut{"\xffa\x00", false}},
	}
	var containsPrefix, containsSuffix, keepPrefix, keepSuffix int
	for _, c := range cases {
		if got := strings.HasPrefix(c.text, c.affix); got != c.prefix.found {
			log.Fatalf("HasPrefix(%q, %q) = %t, want %t", c.text, c.affix, got, c.prefix.found)
		}
		if got := strings.HasSuffix(c.text, c.affix); got != c.suffix.found {
			log.Fatalf("HasSuffix(%q, %q) = %t, want %t", c.text, c.affix, got, c.suffix.found)
		}
		if got := outcome(strings.CutPrefix(c.text, c.affix)); got != c.prefix {
			log.Fatalf("CutPrefix(%q, %q) = %+v, want %+v", c.text, c.affix, got, c.prefix)
		}
		if got := outcome(strings.CutSuffix(c.text, c.affix)); got != c.suffix {
			log.Fatalf("CutSuffix(%q, %q) = %+v, want %+v", c.text, c.affix, got, c.suffix)
		}
		if contains := strings.Contains(c.text, c.affix); contains != c.prefix.found {
			containsPrefix++
		}
		if contains := strings.Contains(c.text, c.affix); contains != c.suffix.found {
			containsSuffix++
		}
		if c.prefix.found && c.text[:len(c.affix)] != c.prefix.rest {
			keepPrefix++
		}
		if c.suffix.found && c.text[len(c.text)-len(c.affix):] != c.suffix.rest {
			keepSuffix++
		}
	}
	failures := [4]int{containsPrefix, containsSuffix, keepPrefix, keepSuffix}
	if want := [4]int{2, 3, 4, 3}; failures != want {
		log.Fatalf("affix mutant failures %v, want %v", failures, want)
	}
	fmt.Printf("%s %s/%s: %d affix fixtures passed, mutant failures %v\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, len(cases), failures)
}
