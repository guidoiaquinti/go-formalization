package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	noisy := strings.Repeat("©", 512)
	cases := []struct {
		text string
		key  rune
		want bool
	}{
		{"", 0, false},
		{"\x00a", 0, true},
		{"©", 0xa9, true},
		{"©", 'é', false},
		{noisy, 'é', false},
		{noisy + "é", 'é', true},
		{"\xa9é", 'é', true},
		{"\xc3", 'é', false},
		{"\xf0\x9f\x9a", '🚀', false},
		{"\xf0\x9f\x9a\x80", '🚀', true},
		{"\xe0\xa0", utf8.RuneError, true},
		{"\xed\xa0\x80", utf8.RuneError, true},
		{"\xf4\x90\x80\x80", utf8.RuneError, true},
		{"\xc0\xaf", utf8.RuneError, true},
		{"�", utf8.RuneError, true},
		{"�", 0xd800, false},
		{"�", -1, false},
		{"�", 0x110000, false},
		{noisy, utf8.RuneError, false},
	}
	var input bytes.Buffer
	encoded, invalid := 0, 0
	for _, c := range cases {
		fmt.Fprint(&input, c.key)
		for _, b := range []byte(c.text) {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		got := bytes.ContainsRune([]byte(c.text), c.key)
		if got != c.want {
			log.Fatalf("Go key %#x returned %t, want %t", c.key, got, c.want)
		}
		if bytes.Contains([]byte(c.text), []byte(string(c.key))) != c.want {
			encoded++
		}
		key := c.key
		if !utf8.ValidRune(key) {
			key = utf8.RuneError
		}
		if bytes.ContainsRune([]byte(c.text), key) != c.want {
			invalid++
		}
	}
	for i, row := range check.LeanRows("godiff-bytes-contains-rune", len(cases), &input) {
		if row != fmt.Sprint(cases[i].want) {
			log.Fatalf("Lean fixture %d returned %s, want %t", i, row, cases[i].want)
		}
	}
	if encoded != 7 || invalid != 3 {
		log.Fatalf("mutation counts encoded=%d invalid=%d, want 7 and 3", encoded, invalid)
	}
	fmt.Printf("bytes.ContainsRune passed %d fixed Go and Lean fixtures; encoded search and invalid-replacement mutants failed %d and %d.\n",
		len(cases), encoded, invalid)
}
