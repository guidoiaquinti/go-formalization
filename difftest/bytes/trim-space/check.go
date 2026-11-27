package main

import (
	"bytes"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"unicode"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

var spaces = []rune{
	9, 10, 11, 12, 13, 32, 0x85, 0xa0, 0x1680, 0x2000, 0x2001, 0x2002,
	0x2003, 0x2004, 0x2005, 0x2006, 0x2007, 0x2008, 0x2009, 0x200a,
	0x2028, 0x2029, 0x202f, 0x205f, 0x3000,
}

func main() {
	if unicode.Version != "17.0.0" {
		log.Fatalf("Unicode version %s, want 17.0.0", unicode.Version)
	}
	membership := make(map[rune]bool, len(spaces))
	for _, r := range spaces {
		membership[r] = true
	}
	for r := rune(0); r <= unicode.MaxRune; r++ {
		if membership[r] != unicode.IsSpace(r) {
			log.Fatalf("White_Space disagreement at U+%04X", r)
		}
	}

	edges := []string{"", " \t\n\r\v\f", " \x00 ", " \xff ", " \xe2\x82 ", " \xed\xa0\x80 ",
		" \xc2", "\xc2 \xa0", "\xc2\xa0 X \xe2\x82", " X\u2003 ", "\u2003X ",
		"\u00a0\u2028", "\u200bX\ufeff", "\u0085�\u2029"}
	for _, r := range spaces {
		for _, adjacent := range []rune{r - 1, r, r + 1} {
			raw := string(adjacent)
			edges = append(edges, raw, raw+"X"+raw, "X"+raw, raw+"X", "\t"+raw+"X"+raw+" ",
				"\u2003"+raw+"\xff"+raw+"\u00a0")
		}
	}
	for value := range 256 {
		raw := string([]byte{byte(value)})
		edges = append(edges, raw, " \t"+raw+"\r ", raw+"\u00a0", "\u2003"+raw,
			"\u2028"+raw+"X"+raw+"\u2029", "\tX"+raw+"\n")
	}
	cases := make([]input, 0, len(edges)+2)
	for i, text := range edges {
		cases = append(cases, input{text, i % 7})
	}
	cases = append(cases, input{"", -1}, input{"", 4})
	check.CompareSamples("godiff-bytes-trim-space", cases, func(r *rand.Rand, i int) input {
		return input{sample(r, i), r.Intn(12)}
	}, encode, func(c input) string {
		var src []byte
		if c.spare >= 0 {
			src = make([]byte, len(c.text), len(c.text)+c.spare)
			copy(src, c.text)
		}
		got := bytes.TrimSpace(src)
		if got == nil {
			return "nil"
		}
		return describe(string(got), cap(got))
	}, func(c input) string {
		start, end := -1, len(c.text)
		lastSpace := true
		for offset, r := range c.text {
			if !lastSpace {
				end = offset
			}
			lastSpace = membership[r]
			if !lastSpace && start < 0 {
				start = offset
			}
		}
		if start < 0 {
			return "nil"
		}
		if !lastSpace {
			end = len(c.text)
		}
		return describe(c.text[start:end], len(c.text)+c.spare-start)
	})
}

type input struct {
	text  string
	spare int
}

func sample(r *rand.Rand, i int) string {
	tokens := []string{"a", "\x00", "�", "\xff", "\xe2\x82", "\xed\xa0\x80", "\u200b", "\ufeff", "🚀"}
	var text strings.Builder
	for range r.Intn(40) {
		switch r.Intn(3) {
		case 0:
			text.WriteByte(byte(r.Intn(256)))
		case 1:
			text.WriteRune(spaces[r.Intn(len(spaces))])
		default:
			text.WriteString(tokens[r.Intn(len(tokens))])
		}
	}
	value := text.String()
	switch i % 4 {
	case 0:
		return strings.Repeat(" ", r.Intn(10)) + "X" + value + "\t\r"
	case 1:
		return "\t\u00a0" + value + "\u2028 "
	case 2:
		return "\tX" + value + "\u3000 "
	default:
		return value
	}
}

func encode(c input) string {
	text := c.text
	row := strconv.AppendInt(nil, int64(max(0, c.spare)), 10)
	for i := range len(text) {
		row = strconv.AppendUint(append(row, ' '), uint64(text[i]), 10)
	}
	return string(row)
}

func describe(text string, capacity int) string {
	row := strconv.AppendInt(nil, int64(capacity), 10)
	row = append(row, ':')
	for i := range len(text) {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(text[i]), 10)
	}
	return string(row)
}
