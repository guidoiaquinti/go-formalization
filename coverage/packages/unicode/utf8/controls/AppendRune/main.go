package main

import (
	"fmt"
	"log"
	"slices"
	"unicode/utf8"
)

func main() {
	fixtures := []struct {
		prefix string
		r      rune
		suffix string
	}{
		{"", 0, "\x00"},
		{"x", 0x7f, "\x7f"},
		{"\xff", 0x80, "\xc2\x80"},
		{"ab", 0x800, "\xe0\xa0\x80"},
		{"\x00", 0x10000, "\xf0\x90\x80\x80"},
		{"a", 0x10ffff, "\xf4\x8f\xbf\xbf"},
		{"", -1, "\xef\xbf\xbd"},
		{"\xff", 0xd800, "\xef\xbf\xbd"},
		{"ab", 0x110000, "\xef\xbf\xbd"},
	}
	var droppedPrefix, skippedInvalid int
	for _, f := range fixtures {
		want := f.prefix + f.suffix
		for _, capacity := range []int{0, utf8.UTFMax} {
			dst := make([]byte, len(f.prefix), len(f.prefix)+capacity)
			copy(dst, f.prefix)
			got := utf8.AppendRune(dst, f.r)
			if string(got) != want || !slices.Equal(dst, []byte(f.prefix)) {
				log.Fatalf("AppendRune(% x, %U) = % x, want % x and preserved prefix", f.prefix, f.r, got, want)
			}
		}
		if string(utf8.AppendRune(nil, f.r)) != want {
			droppedPrefix++
		}
		if !utf8.ValidRune(f.r) && f.prefix != want {
			skippedInvalid++
		}
	}
	if droppedPrefix != 7 || skippedInvalid != 3 {
		log.Fatalf("drop-prefix mutant failed %d, skip-invalid mutant failed %d", droppedPrefix, skippedInvalid)
	}
	fmt.Printf("Go passed %d fixtures at both capacities; mutants failed %d and %d\n", len(fixtures), droppedPrefix, skippedInvalid)
}
