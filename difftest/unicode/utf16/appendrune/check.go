package main

import (
	"math"
	"math/rand"
	"strconv"
	"unicode/utf16"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	r    rune
	dst  []uint16
	room int
}

func main() {
	runes := []rune{math.MinInt32, -1, 0, 127, 128, 0xd7ff, 0xd800, 0xdbff,
		0xdc00, 0xdfff, 0xe000, 0xfffd, 0xffff, 0x10000, 0x10001,
		0x10fffe, 0x10ffff, 0x110000, math.MaxInt32}
	prefixes := [][]uint16{nil, {}, {0}, {0xffff, 0xd800, 1}}
	var edges []input
	for _, r := range runes {
		for _, dst := range prefixes {
			for room := range 3 {
				edges = append(edges, input{r, dst, room})
			}
		}
	}
	check.CompareSamples("godiff-utf16-appendrune", edges, sample, func(c input) string {
		return string(appendWords(strconv.AppendInt(nil, int64(c.r), 10), c.dst))
	}, actual, expected)
}

func sample(rng *rand.Rand, i int) input {
	r := rune(rng.Uint32())
	switch i % 3 {
	case 0:
		r = rune(rng.Intn(0x10000))
	case 1:
		r = 0x10000 + rune(rng.Intn(0x100000))
	}
	dst := make([]uint16, rng.Intn(9))
	for j := range dst {
		dst[j] = uint16(rng.Uint32())
	}
	return input{r, dst, rng.Intn(4)}
}

func actual(c input) string {
	var dst []uint16
	if c.dst != nil || c.room != 0 {
		dst = make([]uint16, len(c.dst), len(c.dst)+c.room)
		copy(dst, c.dst)
	}
	return string(appendWords(nil, utf16.AppendRune(dst, c.r)))
}

func expected(c input) string {
	r := c.r
	out := appendWords(nil, c.dst)
	switch {
	case r>>16 == 0 && r>>10 != 54 && r>>10 != 55:
		out = appendWords(out, []uint16{uint16(r)})
	case 1 <= r>>16 && r>>16 <= 16:
		payload := r - 65536
		out = appendWords(out, []uint16{uint16(55296 + payload/1024), uint16(56320 + payload%1024)})
	default:
		out = appendWords(out, []uint16{65533})
	}
	return string(out)
}

func appendWords(out []byte, words []uint16) []byte {
	for _, word := range words {
		if len(out) != 0 {
			out = append(out, ' ')
		}
		out = strconv.AppendUint(out, uint64(word), 10)
	}
	return out
}
