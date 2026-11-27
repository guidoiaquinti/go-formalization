package utf8decode

import (
	"fmt"
	"math/rand"
	"strconv"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func Compare(executable string, decode func([]byte) (rune, int)) {
	edges := [][]byte{nil, {}}
	for a := range 256 {
		edges = append(edges, []byte{byte(a)})
		for b := range 256 {
			edges = append(edges, []byte{byte(a), byte(b)})
		}
	}
	for _, word := range []string{
		"\xe0\x9f\xbf", "\xe0\xa0\x80", "\xed\x9f\xbf", "\xed\xa0\x80",
		"\xf0\x8f\xbf\xbf", "\xf0\x90\x80\x80", "\xf4\x8f\xbf\xbf", "\xf4\x90\x80\x80",
		"\xef\xbf\xbd", "\xe2\x82", "\xf0\x90\x80", "\xf1\x80\xff\x80", "\xf1\x80\x80\xff",
		"€x", "🐈x",
	} {
		edges = append(edges, []byte(word))
	}
	check.CompareSamples(executable, edges, func(r *rand.Rand, i int) []byte {
		if i%2 == 0 {
			raw := make([]byte, r.Intn(9))
			r.Read(raw)
			return raw
		}
		raw := utf8.AppendRune(nil, rune(r.Intn(0x110000)))
		switch i % 6 {
		case 1:
			return raw[:r.Intn(len(raw)+1)]
		case 3:
			raw[r.Intn(len(raw))] = byte(r.Uint32())
		}
		return append(raw, byte(r.Uint32()))
	}, func(data []byte) string {
		row := []byte("0")
		for _, b := range data {
			row = strconv.AppendUint(append(row, ' '), uint64(b), 10)
		}
		return string(row)
	}, func(data []byte) string {
		value, width := decode(data)
		return fmt.Sprintf("%d %d", value, width)
	}, func(data []byte) string {
		value, width := Expected(data)
		return fmt.Sprintf("%d %d", value, width)
	})
}

// Expected accumulates each payload first, then checks scalar legality instead of using
// the source's lead-byte and second-byte table.
func Expected(data []byte) (rune, int) {
	if len(data) == 0 {
		return 0xfffd, 0
	}
	b := data[0]
	var value, minimum rune
	var width int
	switch {
	case b < 0x80:
		return rune(b), 1
	case b&0xe0 == 0xc0:
		width, value, minimum = 2, rune(b&31), 0x80
	case b&0xf0 == 0xe0:
		width, value, minimum = 3, rune(b&15), 0x800
	case b&0xf8 == 0xf0:
		width, value, minimum = 4, rune(b&7), 0x10000
	default:
		return 0xfffd, 1
	}
	if len(data) < width {
		return 0xfffd, 1
	}
	for _, next := range data[1:width] {
		if next&0xc0 != 0x80 {
			return 0xfffd, 1
		}
		value = value*64 + rune(next&63)
	}
	if value < minimum || value > 0x10ffff || (0xd800 <= value && value <= 0xdfff) {
		return 0xfffd, 1
	}
	return value, width
}
