package main

import (
	"math/rand"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct{ text, cutset string }

func main() {
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
	check.CompareSamples("godiff-strings-trim", edges, sample, encode, func(c input) string {
		return describe(strings.Trim(c.text, c.cutset))
	}, expected)
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

func expected(c input) string {
	membership := make(map[rune]bool)
	for _, r := range c.cutset {
		membership[r] = true
	}
	start, end := -1, len(c.text)
	// Track byte positions from one forward range pass. This preserves
	// malformed trailing bytes in the result whenever they are retained.
	previousMember := true
	for position, rune := range c.text {
		if !previousMember {
			end = position
		}
		previousMember = membership[rune]
		if !previousMember && start < 0 {
			start = position
		}
	}
	if start < 0 {
		return describe("")
	}
	if !previousMember {
		end = len(c.text)
	}
	return describe(c.text[start:end])
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
