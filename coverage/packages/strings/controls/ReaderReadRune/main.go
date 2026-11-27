package main

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

func main() {
	cases := []struct {
		bytes string
		rune  rune
		width int
	}{
		{"A", 'A', 1},
		{"é", 'é', 2},
		{"€", '€', 3},
		{"\xf4\x8f\xbf\xbf", 0x10ffff, 4},
		{"\xed\xa0\x80", utf8.RuneError, 1},
		{"\xe2\x82", utf8.RuneError, 1},
		{"\xef\xbf\xbd", utf8.RuneError, 3},
	}
	var advanceOne, wholeInvalid int
	for _, c := range cases {
		r := strings.NewReader(c.bytes)
		got, width, err := r.ReadRune()
		if got != c.rune || width != c.width || err != nil || r.Len() != len(c.bytes)-c.width {
			panic("incorrect rune, width, or cursor")
		}
		if err := r.UnreadRune(); err != nil || r.Len() != len(c.bytes) {
			panic("previous rune did not retain the starting byte offset")
		}
		if c.width != 1 {
			advanceOne++
		}
		if c.rune == utf8.RuneError && c.width == 1 && len(c.bytes) > 1 {
			wholeInvalid++
		}
	}
	end := strings.NewReader("é")
	if _, _, err := end.ReadRune(); err != nil {
		panic(err)
	}
	if got, width, err := end.ReadRune(); got != 0 || width != 0 || err != io.EOF || end.UnreadRune() == nil {
		panic("EOF did not clear the marker")
	}
	if advanceOne != 4 || wholeInvalid != 2 {
		panic("unexpected mutant failures")
	}
	fmt.Printf("Reader.ReadRune: Go passed eight fixtures, byte-only advance=%d/7, whole-invalid-sequence=%d/7\n", advanceOne, wholeInvalid)
}
