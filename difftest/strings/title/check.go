package main

import (
	"math/rand"
	"strconv"
	"strings"
	"unicode"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	edges := []string{
		"", "hello WORLD", "one_two", "a-b.c/d", "1one 2two", "a\x00b\x7fc",
		"ǆǆ ǆ", "ǄǅǆǇǈǉǊǋǌǱǲǳ", "élan déjà", "Straße", "οσ.οσ",
		"a—b", "a’b", "a。b", "a🚀b", "a\u0301b", "a\u0661b",
		"a�b", "a\xffb", "\xffa\x80b", "a\xe2\x82b", "a\xc0\x80b",
		"a\xed\xa0\x80b", "a\xf4\x90\x80\x80b", "a\xf0\x9f\x9a\x80b",
		strings.Repeat("one—two three_4ǆ\xffa\x00", 128),
	}
	for b := range 256 {
		raw := string([]byte{byte(b)})
		edges = append(edges, raw, "a"+raw+"b", raw+"ǆ", "\xff"+raw+"z")
	}
	for _, r := range []rune{
		0x84, 0x85, 0x86, 0x9f, 0xa0, 0xa1, 0x167f, 0x1680, 0x1681,
		0x1fff, 0x2000, 0x2001, 0x2009, 0x200a, 0x200b, 0x2027, 0x2028,
		0x2029, 0x202a, 0x202e, 0x202f, 0x2030, 0x205e, 0x205f, 0x2060,
		0x2fff, 0x3000, 0x3001, 0xfeff, unicode.MaxRune,
	} {
		edges = append(edges, "a"+string(r)+"ǆ")
	}
	check.CompareSamples("godiff-strings-title", edges, sample, encode,
		func(input string) string { return describe(strings.Title(input)) },
		func(input string) string { return describe(expected(input)) })
}

// Above ASCII, the only separators are Unicode spaces. None are
// letters or decimal digits, so their tests do not remove any spaces.
func expected(input string) string {
	var output strings.Builder
	atBoundary := true
	for _, original := range input {
		mapped := original
		if atBoundary {
			mapped = unicode.ToTitle(original)
		}
		output.WriteRune(mapped)
		if original > unicode.MaxASCII {
			atBoundary = unicode.IsSpace(original)
		} else {
			atBoundary = !strings.ContainsRune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz", original)
		}
	}
	return output.String()
}

func sample(rng *rand.Rand, index int) string {
	tokens := []string{"abc", "a b", "ǆ", "ǳ", "Σ", "é", "—", "🚀", "�", "_", "-",
		"\u200a", "\u200b", "\u2028", "\u1680", "\x00", "\xff", "\xe2\x82", "\xed\xa0\x80", "\xc0\x80"}
	var output strings.Builder
	for range rng.Intn(40) {
		switch index % 4 {
		case 0:
			output.WriteByte(byte(rng.Intn(128)))
		case 1:
			output.WriteRune(rune(rng.Intn(unicode.MaxRune + 1)))
		case 2:
			output.WriteByte(byte(rng.Intn(256)))
		default:
			output.WriteString(tokens[rng.Intn(len(tokens))])
		}
	}
	return output.String()
}

func encode(input string) string {
	row := []byte("0")
	for i := range len(input) {
		row = strconv.AppendUint(append(row, ' '), uint64(input[i]), 10)
	}
	return string(row)
}

func describe(output string) string {
	row := []byte("bytes:")
	for i := range len(output) {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(output[i]), 10)
	}
	return string(row)
}
