package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		log.Fatalf("Go version %s, want %s", runtime.Version(), want)
	}
	trims := []struct{ text, affix, prefix, suffix string }{
		{"", "", "", ""},
		{"abc", "", "abc", "abc"},
		{"abc", "abcd", "abc", "abc"},
		{"ababXabab", "ab", "abXabab", "ababXab"},
		{"xaba", "ab", "xaba", "xaba"},
		{"\xffx\x00", "\xff", "x\x00", "\xffx\x00"},
	}
	var prefixFailures, suffixFailures int
	for _, c := range trims {
		prefix, suffix := strings.TrimPrefix(c.text, c.affix), strings.TrimSuffix(c.text, c.affix)
		if prefix != c.prefix || suffix != c.suffix {
			log.Fatalf("trim(%q, %q) = (%q, %q), want (%q, %q)", c.text, c.affix, prefix, suffix, c.prefix, c.suffix)
		}
		if strings.TrimPrefix(prefix, c.affix) != c.prefix {
			prefixFailures++
		}
		if strings.TrimSuffix(suffix, c.affix) != c.suffix {
			suffixFailures++
		}
	}
	counts := []struct {
		text, separator string
		want            int
	}{
		{"", "", 1},
		{"é\xff猫", "", 4},
		{"\xe0\xa0", "", 3},
		{"aaaaa", "aa", 2},
		{"ababa", "aba", 1},
		{"\x00\xff\x00", "\x00", 2},
		{"x", "xy", 0},
	}
	var countFailures int
	for _, c := range counts {
		if got := strings.Count(c.text, c.separator); got != c.want {
			log.Fatalf("Count(%q, %q) = %d, want %d", c.text, c.separator, got, c.want)
		}
		mutant := len(c.text) + 1
		if c.separator != "" {
			mutant = 0
			for i := 0; i+len(c.separator) <= len(c.text); i++ {
				if c.text[i:i+len(c.separator)] == c.separator {
					mutant++
				}
			}
		}
		if mutant != c.want {
			countFailures++
		}
	}
	if prefixFailures != 1 || suffixFailures != 1 || countFailures != 3 {
		log.Fatalf("mutant failures: prefix=%d, suffix=%d, count=%d, want 1, 1, 3", prefixFailures, suffixFailures, countFailures)
	}
	fmt.Printf("%s %s/%s: trim and count controls passed and rejected all three mutants\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
