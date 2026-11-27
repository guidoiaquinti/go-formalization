package main

import (
	"bytes"
	"math/rand"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct{ text, chars string }

func main() {
	var edges []input
	for a := range 256 {
		for b := range 256 {
			edges = append(edges, input{string([]byte{byte(a)}), string([]byte{byte(b)})})
		}
	}
	tokens := []string{"", "\x00", "\x7f", "a", "é", "🚀", "\ufffd", "\xff", "\xc3", "\xf0\x9f", "\xed\xa0\x80", "\xf4\x90\x80\x80", "é🚀é", "�\xff", "\xff�"}
	charsets := append(append([]string(nil), tokens...), "aa", "a\x00\x7f", "aé", "a\xff", "\ufffdé", "🚀a")
	for _, text := range tokens {
		for _, chars := range charsets {
			for _, length := range []int{0, 7, 8, 9, 10, 15, 16, 31, 32, 65} {
				padded := strings.Repeat("q", length)
				edges = append(edges, input{padded + text, chars}, input{text + padded, chars})
			}
		}
	}
	check.CompareSamples("godiff-bytes-last-index-any", edges, sample, func(c input) string {
		row := strconv.AppendInt(nil, int64(len(c.text)), 10)
		for _, part := range []string{c.text, c.chars} {
			for i := range len(part) {
				row = strconv.AppendUint(append(row, ' '), uint64(part[i]), 10)
			}
		}
		return string(row)
	}, func(c input) string {
		return strconv.Itoa(bytes.LastIndexAny([]byte(c.text), c.chars))
	}, func(c input) string {
		members := make(map[rune]bool)
		for _, runeValue := range c.chars {
			members[runeValue] = true
		}
		last := -1
		for i, runeValue := range c.text {
			if members[runeValue] {
				last = i
			}
		}
		return strconv.Itoa(last)
	})
}

func sample(rng *rand.Rand, i int) input {
	text := make([]byte, rng.Intn(97))
	chars := make([]byte, rng.Intn(21))
	rng.Read(text)
	rng.Read(chars)
	switch i % 4 {
	case 0:
		for i := range chars {
			chars[i] &= 0x7f
		}
	case 1:
		chars = []byte{byte(rng.Intn(256))}
	case 2:
		chars = append(chars, "🚀"...)
		text = append(text, "🚀"...)
	}
	return input{string(text), string(chars)}
}
