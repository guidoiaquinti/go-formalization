package main

import (
	"fmt"
	"log"
	"unicode/utf8"
)

type decode func(string) (rune, int)

func rejectShortAsEmpty(s string) (rune, int) {
	r, n := utf8.DecodeRuneInString(s)
	if len(s) != 0 && !utf8.FullRuneInString(s) {
		return utf8.RuneError, 0
	}
	return r, n
}

func replacementIsAlwaysMalformed(s string) (rune, int) {
	r, n := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && n != 0 {
		return r, 1
	}
	return r, n
}

func finalPayloadIsZero(s string) (rune, int) {
	r, n := utf8.DecodeRuneInString(s)
	if n > 1 {
		r &^= 63
	}
	return r, n
}

func main() {
	fixtures := []struct {
		s string
		r rune
		n int
	}{
		{"", utf8.RuneError, 0},
		{"\x00\xff", 0, 1},
		{"\xc2\xa2!", 0xa2, 2},
		{"\xe2\x82", utf8.RuneError, 1},
		{"\xe0\x9f\xbf", utf8.RuneError, 1},
		{"\xed\xa0\x80", utf8.RuneError, 1},
		{"\xe2\x82\xac!", 0x20ac, 3},
		{"\xef\xbf\xbd", utf8.RuneError, 3},
		{"\xf0\x9f\x90\x88!", 0x1f408, 4},
		{"\xf4\x8f\xbf\xbf", 0x10ffff, 4},
		{"\xf4\x90\x80\x80", utf8.RuneError, 1},
		{"\xf1\x80\x80\xc0", utf8.RuneError, 1},
	}
	for _, candidate := range []struct {
		name string
		fn   decode
		want int
	}{
		{"Go", utf8.DecodeRuneInString, 0},
		{"Treat truncation as empty", rejectShortAsEmpty, 1},
		{"Treat valid replacement rune as malformed", replacementIsAlwaysMalformed, 1},
		{"Ignore final payload", finalPayloadIsZero, 5},
	} {
		bad := 0
		for _, f := range fixtures {
			r, n := candidate.fn(f.s)
			if r != f.r || n != f.n {
				bad++
				if candidate.want == 0 {
					log.Fatalf("DecodeRuneInString(% x) = (%U, %d), want (%U, %d)", f.s, r, n, f.r, f.n)
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", candidate.name, bad, len(fixtures))
		if bad != candidate.want {
			log.Fatalf("%s expected %d discrepancies", candidate.name, candidate.want)
		}
	}
}
