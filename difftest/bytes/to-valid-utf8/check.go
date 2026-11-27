package main

import (
	"bytes"
	"math/rand"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	text, replacement string
}

func main() {
	texts := []string{"", "ascii", "�", "é🚀", "\xff", "\xff\xfc", "\xffa\xfc",
		"\xff�\xfc", "\xffé\xfc", "\xe2\x82", "\xe2\x82a\xa0", "\xed\xa0\x80",
		"\xc0\xaf", "\xf4\x90\x80\x80", strings.Repeat("\xffX", 80)}
	replacements := []string{"", "?", "�", "\x00", "\xff\xe2", "替",
		strings.Repeat("x", 64), strings.Repeat("\xff\x80", 160)}
	var edges []input
	for _, text := range texts {
		for _, replacement := range replacements {
			edges = append(edges, input{text, replacement})
		}
	}
	for i := range 256 {
		raw := string([]byte{byte(i)})
		for _, replacement := range replacements {
			edges = append(edges, input{raw, replacement}, input{raw + "\xffé\x80" + raw, replacement})
		}
	}
	check.CompareSamples("godiff-bytes-to-valid-utf8", edges, sample, encode, native, expected)
}

func native(c input) string {
	got := bytes.ToValidUTF8([]byte(c.text), []byte(c.replacement))
	if got == nil {
		return "nil"
	}
	initial := len(c.text) + len(c.replacement)
	capacity := strconv.Itoa(cap(got))
	if len(got) > initial && cap(got) >= len(got) {
		capacity = "growth"
	}
	return describe(string(got), capacity)
}

func expected(c input) string {
	var output strings.Builder
	invalid := false
	for offset, r := range c.text {
		if r == utf8.RuneError && !strings.HasPrefix(c.text[offset:], "\xef\xbf\xbd") {
			if !invalid {
				output.WriteString(c.replacement)
			}
			invalid = true
			continue
		}
		invalid = false
		output.WriteRune(r)
	}
	initial := len(c.text) + len(c.replacement)
	capacity := strconv.Itoa(initial)
	if output.Len() > initial {
		capacity = "growth"
	}
	return describe(output.String(), capacity)
}

func sample(r *rand.Rand, i int) input {
	var text, replacement strings.Builder
	tokens := []string{"a", "é", "🚀", "�", "\xff", "\xe2\x82", "\xed\xa0\x80", "\x00"}
	for range r.Intn(56) {
		if r.Intn(2) == 0 {
			text.WriteByte(byte(r.Intn(256)))
		} else {
			text.WriteString(tokens[r.Intn(len(tokens))])
		}
	}
	for range r.Intn(4 + i%20) {
		replacement.WriteByte(byte(r.Intn(256)))
	}
	return input{text.String(), replacement.String()}
}

func encode(c input) string {
	row := strconv.AppendInt(nil, int64(len(c.replacement)), 10)
	for _, text := range []string{c.replacement, c.text} {
		for i := range len(text) {
			row = strconv.AppendUint(append(row, ' '), uint64(text[i]), 10)
		}
	}
	return string(row)
}

func describe(output, capacity string) string {
	row := append([]byte(capacity), ':')
	for i := range len(output) {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(output[i]), 10)
	}
	return string(row)
}
