package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	if unicode.Version != "17.0.0" {
		log.Fatalf("Unicode version %s, want 17.0.0", unicode.Version)
	}
	cases := []struct {
		left, right string
		want        bool
	}{
		{"", "", true},
		{"a\x00Z", "A\x00z", true},
		{"[", "{", false},
		{"Kelvin K", "kelvin K", true},
		{"aſΣ", "Asς", true},
		{"Θϴ", "ϑθ", true},
		{"ϑ", "Θ", true},
		{"ß", "ẞ", true},
		{"ß", "SS", false},
		{"ﬃ", "ffi", false},
		{"İ", "i", false},
		{"I", "ı", false},
		{"\U00010400", "\U00010428", true},
		{"\xffA", "�a", true},
		{"\xff\x80", "�", false},
		{"\xff\x80", "��", true},
		{"\xe2A", "�a", true},
		{"\xed\xa0\x80", "���", true},
		{"\xf4\x90\x80\x80", "����", true},
		{"éa", "É", false},
		{"é", "Éa", false},
	}
	successors := []struct{ rune, next rune }{
		{'K', 'k'}, {'k', 'K'}, {'K', 'K'},
		{'S', 's'}, {'s', 'ſ'}, {'ſ', 'S'},
		{'Θ', 'θ'}, {'θ', 'ϑ'}, {'ϑ', 'ϴ'}, {'ϴ', 'Θ'},
		{'ß', 'ẞ'}, {'ẞ', 'ß'}, {'İ', 'İ'}, {'�', '�'},
		{0x10400, 0x10428}, {0x10428, 0x10400},
		{0xd800, 0xd800}, {unicode.MaxRune, unicode.MaxRune},
	}
	var input bytes.Buffer
	lower, collapse := 0, 0
	for i, c := range cases {
		fmt.Fprint(&input, "1 ", len(c.left))
		for _, part := range []string{c.left, c.right} {
			for _, b := range []byte(part) {
				fmt.Fprint(&input, " ", b)
			}
		}
		fmt.Fprintln(&input)
		if got := strings.EqualFold(c.left, c.right); got != c.want {
			log.Fatalf("Go equality fixture %d returned %t, want %t", i, got, c.want)
		}
		if (strings.ToLower(c.left) == strings.ToLower(c.right)) != c.want {
			lower++
		}
		if strings.EqualFold(strings.ToValidUTF8(c.left, "�"), strings.ToValidUTF8(c.right, "�")) != c.want {
			collapse++
		}
	}
	for _, c := range successors {
		fmt.Fprintln(&input, 0, c.rune)
		if got := unicode.SimpleFold(c.rune); got != c.next {
			log.Fatalf("Go successor for %U returned %U, want %U", c.rune, got, c.next)
		}
	}
	for i, row := range check.LeanRows("godiff-strings-equal-fold", len(cases)+len(successors), &input) {
		want := ""
		if i < len(cases) {
			want = strconv.FormatBool(cases[i].want)
		} else {
			want = strconv.Itoa(int(successors[i-len(cases)].next))
		}
		if row != want {
			log.Fatalf("Lean fixture %d returned %s, want %s", i, row, want)
		}
	}
	if lower != 4 || collapse != 4 {
		log.Fatalf("mutation counts lowercase=%d collapse=%d, want 4 and 4", lower, collapse)
	}
	fmt.Printf("strings.EqualFold passed %d fixed Go and Lean string pairs and %d successor lookups; lowercase and collapsing invalid runs failed %d and %d.\n",
		len(cases), len(successors), lower, collapse)
}
