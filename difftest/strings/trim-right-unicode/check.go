package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct{ text, cutset string }

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dir, err := os.MkdirTemp("", "golean-strings-trim-right-unicode-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.Setenv("GOLEAN_STRINGS_TRIM_RIGHT_UNICODE_CASES", path); err != nil {
		return err
	}

	texts := []string{"", "a", "aaaa", "ababa", "éaé", "🚀%🚀", "x🚀éy",
		"\x00x\x00", "\xffa\xfe", "�x\x80", "\xe2\x82x\xe2\x82", "\xc0\x80x\xed\xa0\x80",
		"\xed\xa0x\x80", "\xf4\x90\x80\x80a�"}
	cutsets := []string{"", "a", "ab", "\x00", "a\x00", "é", "🚀", "é🚀", "éa🚀",
		"�", "\xff", "\x80a", "\xe2\x82", "aé\x80"}
	var edges []input
	for _, text := range texts {
		for _, cutset := range cutsets {
			edges = append(edges, input{text, cutset})
		}
	}
	for value := range 256 {
		raw := string([]byte{byte(value)})
		for _, cutset := range []string{raw, raw + "a", "a�", "é", "\x00abc"} {
			for _, text := range []string{raw, raw + "a" + raw, raw + "é" + raw, "é" + raw + "🚀"} {
				edges = append(edges, input{text, cutset})
			}
		}
	}
	compare := func(label string, cases []input) error {
		var modelInput, fixture bytes.Buffer
		expected := make([]string, len(cases))
		for i, c := range cases {
			fmt.Fprintln(&modelInput, encode(c))
			want := expectedPrefix(c)
			fmt.Fprintf(&fixture, "%x\t%x\t%x\n", c.cutset, c.text, want)
			expected[i] = describe(want)
		}
		for i, got := range check.LeanRows("godiff-strings-trim-right-unicode", len(cases), &modelInput) {
			if got != expected[i] {
				return fmt.Errorf("%s trimRightUnicode(%x,%x): Lean=%s want=%s", label, cases[i].text,
					cases[i].cutset, got, expected[i])
			}
		}
		if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
			return err
		}
		const fixtures = "difftest/strings/trim-right-unicode/testdata/"
		if err := check.TestInPackage("strings", "TestGoLeanTrimRightUnicode", fixtures+"bridge_test.go",
			fixtures+"trim_right_unicode_test.go"); err != nil {
			return err
		}
		fmt.Printf("strings.trimRightUnicode %s: %d inputs agree.\n", label, len(cases))
		return nil
	}
	if err := compare("boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]input, 20_000)
		for i := range cases {
			cases[i] = sample(r, i)
		}
		if err := compare(fmt.Sprintf("seed=%d", seed), cases); err != nil {
			return err
		}
	}
	return nil
}

func sample(r *rand.Rand, i int) input {
	tokens := []string{"a", "b", "\x00", "é", "🚀", "�", "\xff", "\xe2\x82", "\xed\xa0\x80"}
	var text, cutset strings.Builder
	for range r.Intn(40) {
		if r.Intn(2) == 0 {
			text.WriteByte(byte(r.Intn(256)))
		} else {
			text.WriteString(tokens[r.Intn(len(tokens))])
		}
	}
	switch i % 4 {
	case 0:
		cutset.WriteByte(byte(r.Intn(128)))
	case 1:
		for range 2 + r.Intn(8) {
			cutset.WriteByte(byte(r.Intn(128)))
		}
	case 2:
		for range 1 + r.Intn(6) {
			cutset.WriteString(tokens[r.Intn(len(tokens))])
		}
	default:
		for range r.Intn(12) {
			cutset.WriteByte(byte(r.Intn(256)))
		}
	}
	value, chars := text.String(), cutset.String()
	if r.Intn(2) == 0 && chars != "" {
		_, width := utf8.DecodeRuneInString(chars)
		edge := chars[:width]
		value = strings.Repeat(edge, r.Intn(5)) + value + strings.Repeat(edge, r.Intn(5))
	}
	return input{value, chars}
}

func encode(c input) string {
	row := strconv.AppendInt(nil, int64(len(c.cutset)), 10)
	for _, text := range []string{c.cutset, c.text} {
		for i := range len(text) {
			row = strconv.AppendUint(append(row, ' '), uint64(text[i]), 10)
		}
	}
	return string(row)
}

func expectedPrefix(c input) string {
	membership := make(map[rune]bool)
	for _, r := range c.cutset {
		membership[r] = true
	}
	end := 0
	for position, r := range c.text {
		if !membership[r] {
			_, width := utf8.DecodeRuneInString(c.text[position:])
			end = position + width
		}
	}
	return c.text[:end]
}

func describe(text string) string {
	row := []byte("bytes:")
	for i := range len(text) {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(text[i]), 10)
	}
	return string(row)
}
