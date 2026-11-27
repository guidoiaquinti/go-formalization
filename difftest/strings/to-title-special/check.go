package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"unicode"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	special unicode.SpecialCase
	text    string
}

func entry(from, to uint32, title rune) unicode.CaseRange {
	return unicode.CaseRange{Lo: from, Hi: to, Delta: [3]rune{3, 2, title}}
}

func main() {
	tables := []unicode.SpecialCase{
		nil, {}, unicode.TurkishCase, unicode.AzeriCase,
		{entry('i', 'i', 0)},
		{entry('A', 'D', unicode.UpperLower)},
		{entry('A', 'D', math.MaxInt32)},
		{entry('a', 'z', -100), entry('a', 'z', 0), entry('a', 'z', 10)},
		{entry('A', 'A', 1), entry('Z', 'Z', 0), entry('A', 'C', 4)},
		{entry('A', 'A', 0xd800-'A')},
		{entry('A', 'A', math.MinInt32)},
		{entry(0xfffd, 0xfffd, '?'-0xfffd)},
		{entry(0xfffd, 0xfffd, -0x10000)},
		{entry(0x10ffff, 0x10ffff, 1)},
		{entry(0x80000000, 0x10ffff, unicode.UpperLower)},
		{entry(0x1c4, 0x1c6, 0)},
		{entry(0xffffffff, 0x7fffffff, unicode.UpperLower)},
	}
	words := []string{"", "i", "Iİiı", "ABCD", "abcd", "ZAB", "ΣKÉȺ", "Ǆǅǆ", "A\x00Z", "🚀",
		"�\xff\x80\xe2\x82", "\xed\xa0\x80", "\U0010ffff", strings.Repeat("AIİ\xff", 64)}
	var edges []input
	for _, special := range tables {
		for _, text := range words {
			edges = append(edges, input{special, text})
		}
	}
	for b := range 256 {
		raw := string([]byte{byte(b)})
		special := unicode.SpecialCase{entry(uint32(b), uint32(b), -1)}
		edges = append(edges, input{nil, raw}, input{unicode.TurkishCase, raw + "Iİ"},
			input{special, raw}, input{special, string(rune(b))})
	}
	check.CompareSamples("godiff-strings-to-title-special", edges, sample, encode, func(c input) string {
		return describe(strings.ToTitleSpecial(c.special, c.text))
	}, func(c input) string {
		var out strings.Builder
		for _, r := range c.text {
			if mapped := expected(c.special, r); mapped >= 0 {
				out.WriteRune(mapped)
			}
		}
		return describe(out.String())
	})
}

// Slicing the candidate list preserves midpoint tie breaking, including
// caller-supplied overlapping ranges. A zero-delta hit is still terminal.
func expected(table unicode.SpecialCase, r rune) rune {
	for len(table) != 0 {
		mid := len(table) / 2
		row := table[mid]
		lo, hi := rune(row.Lo), rune(row.Hi)
		switch {
		case lo <= r && r <= hi:
			delta := row.Delta[unicode.TitleCase]
			if delta > unicode.MaxRune {
				base := int64(lo)
				return rune(base + 2*((int64(r)-base)/2))
			}
			return rune(int64(r) + int64(delta))
		case r < lo:
			table = table[:mid]
		default:
			table = table[mid+1:]
		}
	}
	return unicode.ToTitle(r)
}

func sample(rng *rand.Rand, index int) input {
	tokens := []string{"I", "i", "İ", "ı", "ABC", "Ǆ", "Ⱥ", "Σ", "�", "\xff", "\xe2\x82",
		"\x00", "🚀", "\U0010ffff"}
	var text strings.Builder
	for range rng.Intn(32) {
		switch index % 3 {
		case 0:
			text.WriteByte(byte(rng.Intn(128)))
		case 1:
			text.WriteByte(byte(rng.Intn(256)))
		default:
			text.WriteString(tokens[rng.Intn(len(tokens))])
		}
	}
	if index%7 == 0 {
		return input{unicode.TurkishCase, text.String()}
	}
	special := make(unicode.SpecialCase, rng.Intn(13))
	for i := range special {
		var lo, hi uint32
		switch index % 4 {
		case 0:
			lo = uint32(i * 10)
			hi = lo + uint32(rng.Intn(10))
		case 1:
			lo = uint32(rng.Intn(120))
			hi = lo + uint32(rng.Intn(20))
		case 2:
			lo = uint32(rng.Intn(unicode.MaxRune + 1))
			hi = lo + uint32(rng.Intn(1000))
		default:
			lo, hi = rng.Uint32(), rng.Uint32()
		}
		delta := []rune{0, 1, -1, -64, 32, unicode.UpperLower, rune(rng.Uint32())}[rng.Intn(7)]
		special[i] = entry(lo, hi, delta)
	}
	return input{special, text.String()}
}

func encode(c input) string {
	row := strconv.AppendInt(nil, int64(len(c.special)), 10)
	for _, r := range c.special {
		row = fmt.Appendf(row, " %d %d %d %d %d", r.Lo, r.Hi, r.Delta[0], r.Delta[1], r.Delta[2])
	}
	for i := range len(c.text) {
		row = strconv.AppendUint(append(row, ' '), uint64(c.text[i]), 10)
	}
	return string(row)
}

func describe(raw string) string {
	row := []byte("bytes:")
	for i := range len(raw) {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(raw[i]), 10)
	}
	return string(row)
}
