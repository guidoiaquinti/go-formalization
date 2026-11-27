package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	text string
	key  rune
}

func main() {
	edges := make([]input, 0, 66000)
	for first := range 256 {
		for second := range 256 {
			edges = append(edges, input{string([]byte{byte(first), byte(second)}), utf8.RuneError})
		}
	}
	for _, text := range []string{"", "a", "a\x00a", "é\xff", "\uFFFDé", "\xed\xa0\x80", "\xf4\x90\x80\x80", "\xf0\x9f", "🚀🚀", "\u07ff\u0800\ud7ff\ue000\uffff\U00010000\U0010ffff"} {
		for _, r := range []rune{-2147483648, -1, 0, 'a', 127, 128, 'é', 0x7ff, 0x800, 0xd7ff, 0xd800, 0xdfff, 0xe000, utf8.RuneError, 0xffff, 0x10000, '🚀', 0x10ffff, 0x110000, 2147483647} {
			edges = append(edges, input{text, r})
		}
	}
	check.CompareSamples("godiff-strings-index-rune", edges, sample, func(c input) string {
		out := strconv.AppendInt(nil, int64(c.key), 10)
		for i := range len(c.text) {
			out = strconv.AppendUint(append(out, ' '), uint64(c.text[i]), 10)
		}
		return string(out)
	}, func(c input) string {
		return fmt.Sprintf("%d %t", strings.IndexRune(c.text, c.key), strings.ContainsRune(c.text, c.key))
	}, func(c input) string {
		offset := -1
		for i, got := range c.text {
			if got == c.key {
				offset = i
				break
			}
		}
		return fmt.Sprintf("%d %t", offset, offset >= 0)
	})
}

func sample(rng *rand.Rand, i int) input {
	raw := make([]byte, rng.Intn(97))
	rng.Read(raw)
	r := rune(rng.Uint32())
	switch i % 4 {
	case 0:
		r = utf8.RuneError
	case 1:
		r = rune(rng.Intn(128))
	case 2:
		r = rune(rng.Intn(utf8.MaxRune + 1))
		if utf8.ValidRune(r) {
			pos := rng.Intn(len(raw) + 1)
			raw = append(append(append([]byte(nil), raw[:pos]...), string(r)...), raw[pos:]...)
		}
	}
	return input{string(raw), r}
}
