package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"unicode/utf8"
)

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		log.Fatalf("Go version %s, want %s", runtime.Version(), want)
	}
	runes := []struct {
		r     rune
		width int
	}{
		{-1, -1}, {0, 1}, {0x7f, 1}, {0x80, 2}, {0x7ff, 2}, {0x800, 3},
		{0xd7ff, 3}, {0xd800, -1}, {0xdfff, -1}, {0xe000, 3},
		{0xffff, 3}, {0x10000, 4}, {0x10ffff, 4}, {0x110000, -1},
	}
	var replacementFailures, surrogateFailures int
	for _, c := range runes {
		if got := utf8.RuneLen(c.r); got != c.width {
			log.Fatalf("RuneLen(%#x) = %d, want %d", c.r, got, c.width)
		}
		if got := utf8.ValidRune(c.r); got != (c.width > 0) {
			log.Fatalf("ValidRune(%#x) = %t, want %t", c.r, got, c.width > 0)
		}
		if len(string(c.r)) != c.width {
			replacementFailures++
		}
		if (0 <= c.r && c.r <= utf8.MaxRune) != (c.width > 0) {
			surrogateFailures++
		}
	}
	starts := []struct {
		b    byte
		want bool
	}{
		{0x7f, true}, {0x80, false}, {0xbf, false}, {0xc0, true},
		{0xc2, true}, {0xf4, true}, {0xf5, true}, {0xff, true},
	}
	var restrictedFailures int
	for _, c := range starts {
		if got := utf8.RuneStart(c.b); got != c.want {
			log.Fatalf("RuneStart(%#x) = %t, want %t", c.b, got, c.want)
		}
		if (c.b < 0x80 || 0xc2 <= c.b && c.b <= 0xf4) != c.want {
			restrictedFailures++
		}
	}
	if replacementFailures != 4 || surrogateFailures != 2 || restrictedFailures != 3 {
		log.Fatalf("mutant failures: replacement=%d surrogate=%d restricted=%d, want 4, 2, 3",
			replacementFailures, surrogateFailures, restrictedFailures)
	}
	fmt.Printf("%s %s/%s: UTF-8 controls passed on 14 scalars and eight bytes; all three mutants rejected\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
