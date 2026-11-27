package main

import (
	"bytes"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	text string
	seek rune
}

func main() {
	seeks := []rune{math.MinInt32, -1, 0, 'a', 0x7f, 0x80, 0xa9, 0xe9, 0x7ff,
		0x800, 0xd7ff, 0xd800, 0xdfff, 0xe000, utf8.RuneError, 0xffff,
		0x10000, '🚀', utf8.MaxRune, utf8.MaxRune + 1, math.MaxInt32}
	texts := []string{"", "a\x00\x7f", "\u0080©é", "\u07ff\u0800\ud7ff\ue000", "\uffff\U00010000\U0010ffff",
		"é🚀", "�", "\xff", "\xe2\x82", "\xed\xa0\x80", "\xc0\xaf", "\xf4\x90\x80\x80",
		strings.Repeat("©", 128), strings.Repeat("©", 128) + "é"}
	var edges []input
	for _, text := range texts {
		for _, seek := range seeks {
			edges = append(edges, input{text, seek})
		}
	}
	for i := range 256 {
		raw := string([]byte{byte(i)})
		for _, seek := range seeks {
			edges = append(edges, input{raw, seek}, input{raw + "é\xff🚀" + raw, seek})
		}
		edges = append(edges, input{raw, rune(i)})
	}
	check.CompareSamples("godiff-bytes-contains-rune", edges, sample, encode, native, expected)
}

func encode(c input) string {
	row := strconv.AppendInt(nil, int64(c.seek), 10)
	for i := range len(c.text) {
		row = strconv.AppendUint(append(row, ' '), uint64(c.text[i]), 10)
	}
	return string(row)
}

func native(c input) string {
	return strconv.FormatBool(bytes.ContainsRune([]byte(c.text), c.seek))
}

func expected(c input) string {
	for _, r := range c.text {
		if r == c.seek {
			return "true"
		}
	}
	return "false"
}

func sample(r *rand.Rand, i int) input {
	tokens := []string{"a", "\x00", "é", "©", "🚀", "�", "\xff", "\xe2\x82", "\xed\xa0\x80"}
	var text strings.Builder
	for range r.Intn(48) {
		if r.Intn(2) == 0 {
			text.WriteByte(byte(r.Intn(256)))
		} else {
			text.WriteString(tokens[r.Intn(len(tokens))])
		}
	}
	seek := rune(r.Uint32())
	switch i % 5 {
	case 0:
		seek = rune(r.Intn(128))
	case 1:
		seek = utf8.RuneError
	case 2:
		seek = rune(r.Intn(utf8.MaxRune + 1))
		text.WriteRune(seek)
	case 3:
		seek = rune(r.Intn(utf8.MaxRune + 1))
	}
	return input{text.String(), seek}
}
