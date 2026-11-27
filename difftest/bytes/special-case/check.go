package main

import (
	"bytes"
	"fmt"
	"log"
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

func entry(from, to uint32, delta rune) unicode.CaseRange {
	return unicode.CaseRange{Lo: from, Hi: to, Delta: [3]rune{delta, delta, delta}}
}

func main() {
	tables := []unicode.SpecialCase{
		nil, {}, unicode.TurkishCase, unicode.AzeriCase,
		{entry('i', 'i', 0)},
		{{Lo: 65, Hi: 65, Delta: [3]rune{1, 2, 3}}},
		{entry(65, 65, -66)},
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
	check.CompareSamples("godiff-bytes-special-case", edges, sample, encode, actual, reference)

}

// Slicing the candidate list preserves midpoint tie breaking, including
// caller-supplied overlapping ranges. A zero-delta hit is still terminal.
func expected(table unicode.SpecialCase, r rune, variant int) rune {
	for len(table) != 0 {
		mid := len(table) / 2
		row := table[mid]
		lo, hi := rune(row.Lo), rune(row.Hi)
		switch {
		case lo <= r && r <= hi:
			delta := row.Delta[variant]
			if delta > unicode.MaxRune {
				base := int64(lo)
				parity := int64(0)
				if variant == unicode.LowerCase {
					parity = 1
				}
				return rune(base + 2*((int64(r)-base)/2) + parity)
			}
			return rune(int64(r) + int64(delta))
		case r < lo:
			table = table[:mid]
		default:
			table = table[mid+1:]
		}
	}
	return []func(rune) rune{unicode.ToUpper, unicode.ToLower, unicode.ToTitle}[variant](r)
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
		special[i] = unicode.CaseRange{Lo: lo, Hi: hi, Delta: [3]rune{delta, -delta, rune(rng.Uint32())}}
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

func capacity(input, output int) string {
	if output > input {
		return "growth"
	}
	return strconv.Itoa(input)
}

func describe(raw string, nilSlice bool, capacity string) string {
	tag := "slice"
	if nilSlice {
		tag = "nil"
	}
	row := []byte(tag + ":" + capacity + ":")
	for i := range len(raw) {
		if i > 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(raw[i]), 10)
	}
	return string(row)
}

func actual(c input) string {
	backing := make([]byte, len(c.text)+2)
	copy(backing, c.text)
	backing[len(c.text)], backing[len(c.text)+1] = 0xa5, 0x59
	source := backing[:len(c.text)]
	if c.text == "" {
		source = nil
	}
	var results []string
	for _, convert := range []func(unicode.SpecialCase, []byte) []byte{bytes.ToUpperSpecial, bytes.ToLowerSpecial, bytes.ToTitleSpecial} {
		out := convert(c.special, source)
		size := capacity(len(source), len(out))
		if cap(out) < len(source) || (size != "growth" && strconv.Itoa(cap(out)) != size) {
			log.Fatalf("special-case cap(%q→%q)=%d want %s", c.text, out, cap(out), size)
		}
		results = append(results, describe(string(out), out == nil, size))
		for i := range out {
			out[i] ^= 0xff
		}
		if string(backing[:len(source)]) != c.text || backing[len(source)] != 0xa5 || backing[len(source)+1] != 0x59 {
			log.Fatal("special casing shares input storage")
		}
		if c.text == "" {
			empty := convert(c.special, backing[:0])
			if empty == nil || cap(empty) != 0 || len(empty) != 0 {
				log.Fatal("nonnil empty conversion mismatch")
			}
		}
	}
	return strings.Join(results, "|")
}

func reference(c input) string {
	var results []string
	for variant := range 3 {
		var out strings.Builder
		for _, r := range c.text {
			if mapped := expected(c.special, r, variant); mapped >= 0 {
				out.WriteRune(mapped)
			}
		}
		result := out.String()
		results = append(results, describe(result, false, capacity(len(c.text), len(result))))
	}
	return strings.Join(results, "|")
}
