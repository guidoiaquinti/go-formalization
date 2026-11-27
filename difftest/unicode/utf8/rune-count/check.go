package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	edges := [][]byte{nil, {}}
	for a := range 256 {
		edges = append(edges, []byte{byte(a)})
		for b := range 256 {
			edges = append(edges, []byte{byte(a), byte(b)})
		}
	}
	for _, text := range []string{
		"\xe0\x9f\xbf", "\xe0\xa0\x80", "\xed\x9f\xbf", "\xed\xa0\x80", "\xef\xbf\xbd",
		"\xf0\x8f\xbf\xbf", "\xf0\x90\x80\x80", "\xf4\x8f\xbf\xbf", "\xf4\x90\x80\x80",
		"\xf0\x90\x80", "\xf0\x90A\x80", "abc\x80é", "🌍", "\xffé",
	} {
		edges = append(edges, []byte(text), append(bytes.Repeat([]byte{'a'}, 64), []byte(text)...))
	}
	check.CompareSamples("godiff-utf8-rune-count", edges, sample, encode, func(data []byte) string {
		return fmt.Sprintf("%d %d", utf8.RuneCount(data), utf8.RuneCountInString(string(data)))
	}, expected)
}

func sample(rng *rand.Rand, i int) []byte {
	var data []byte
	switch i % 3 {
	case 0:
		data = make([]byte, rng.Intn(129))
		rng.Read(data)
	case 1:
		for range rng.Intn(33) {
			data = utf8.AppendRune(data, rune(rng.Intn(0x110000)))
		}
	default:
		data = bytes.Repeat([]byte{'a'}, rng.Intn(65))
		for range rng.Intn(33) {
			if rng.Intn(2) == 0 {
				data = utf8.AppendRune(data, rune(rng.Intn(0x110000)))
			} else {
				data = append(data, byte(rng.Uint32()))
			}
		}
	}
	return data
}

func encode(data []byte) string {
	row := []byte("0")
	for _, b := range data {
		row = strconv.AppendUint(append(row, ' '), uint64(b), 10)
	}
	return string(row)
}

func expected(data []byte) string {
	n := 0
	for len(data) != 0 {
		data = data[scalarWidth(data):]
		n++
	}
	return fmt.Sprintf("%d %d", n, n)
}

// Decode a candidate arithmetically, then enforce the scalar and minimum
// value constraints. Invalid candidates always advance by one byte.
func scalarWidth(data []byte) int {
	lead := data[0]
	var width, scalar, minimum int
	switch {
	case lead&0xe0 == 0xc0:
		width, scalar, minimum = 2, int(lead&0x1f), 0x80
	case lead&0xf0 == 0xe0:
		width, scalar, minimum = 3, int(lead&0x0f), 0x800
	case lead&0xf8 == 0xf0:
		width, scalar, minimum = 4, int(lead&0x07), 0x10000
	default:
		return 1
	}
	if len(data) < width {
		return 1
	}
	for _, b := range data[1:width] {
		if b&0xc0 != 0x80 {
			return 1
		}
		scalar = scalar*64 + int(b&0x3f)
	}
	if scalar < minimum || scalar > 0x10ffff || 0xd800 <= scalar && scalar < 0xe000 {
		return 1
	}
	return width
}
