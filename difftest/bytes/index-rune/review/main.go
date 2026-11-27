package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	noisy := strings.Repeat("©", 512)
	cases := []struct {
		text string
		key  rune
		want int
	}{
		{"", 0, -1},
		{"é\x00", 0, 2},
		{"©", 'é', -1},
		{noisy, 'é', -1},
		{noisy + "é" + noisy + "é", 'é', 1024},
		{"\xa9é", 'é', 1},
		{"🚀éé", 'é', 4},
		{"éé", 'é', 0},
		{"é\xf0\x9f\x9a", '🚀', -1},
		{"é🚀", '🚀', 2},
		{"é\xe0\xa0�", utf8.RuneError, 2},
		{"é\xed\xa0\x80", utf8.RuneError, 2},
		{"é\xf4\x90\x80\x80", utf8.RuneError, 2},
		{"é\xc0\xaf", utf8.RuneError, 2},
		{"é�\xff", utf8.RuneError, 2},
		{"é�", 0xd800, -1},
		{"é�", -1, -1},
		{"é�", 0x110000, -1},
		{noisy, utf8.RuneError, -1},
	}
	var input bytes.Buffer
	encoded, runeOrdinal, last := 0, 0, 0
	for i, c := range cases {
		fmt.Fprint(&input, c.key)
		for _, b := range []byte(c.text) {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		if got := bytes.IndexRune([]byte(c.text), c.key); got != c.want {
			log.Fatalf("Go fixture %d returned %d, want %d", i, got, c.want)
		}
		if bytes.Index([]byte(c.text), []byte(string(c.key))) != c.want {
			encoded++
		}
		if c.want >= 0 && utf8.RuneCountInString(c.text[:c.want]) != c.want {
			runeOrdinal++
		}
		lastOffset := -1
		for offset, r := range c.text {
			if r == c.key {
				lastOffset = offset
			}
		}
		if lastOffset != c.want {
			last++
		}
	}
	for i, row := range check.LeanRows("godiff-bytes-index-rune", len(cases), &input) {
		if row != strconv.Itoa(cases[i].want) {
			log.Fatalf("Lean fixture %d returned %s, want %d", i, row, cases[i].want)
		}
	}
	if encoded != 7 || runeOrdinal != 9 || last != 8 {
		log.Fatalf("mutation counts encoded=%d rune ordinal=%d last=%d, want 7, 9 and 8", encoded, runeOrdinal, last)
	}
	fmt.Printf("bytes.IndexRune passed %d fixed Go and Lean fixtures; encoded search, rune ordinal, and last occurrence mutants failed %d, %d and %d.\n",
		len(cases), encoded, runeOrdinal, last)
}
