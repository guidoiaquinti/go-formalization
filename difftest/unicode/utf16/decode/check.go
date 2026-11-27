package main

import (
	"math/rand"
	"strconv"
	"unicode/utf16"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/utf16oracle"
)

func main() {
	bounds := []uint16{0, 0x7f, 0x80, 0xd7ff, 0xd800, 0xdbff, 0xdc00, 0xdfff, 0xe000, 0xffff}
	edges := [][]uint16{nil, {}}
	for _, first := range bounds {
		edges = append(edges, []uint16{first})
		for _, second := range bounds {
			edges = append(edges, []uint16{first, second})
			for _, third := range bounds {
				edges = append(edges, []uint16{first, second, third})
			}
		}
	}
	check.CompareSamples("godiff-utf16-decode", edges, sample, func(words []uint16) string {
		text := []byte("0")
		for _, word := range words {
			text = strconv.AppendUint(append(text, ' '), uint64(word), 10)
		}
		return string(text)
	}, func(words []uint16) string {
		text := []byte("runes")
		for _, r := range utf16.Decode(words) {
			text = strconv.AppendInt(append(text, ' '), int64(r), 10)
		}
		return string(text)
	}, expected)
}

func sample(rng *rand.Rand, i int) []uint16 {
	length := rng.Intn(17)
	if i%5 == 0 {
		length = 62 + rng.Intn(5)
	}
	words := make([]uint16, length)
	for j := range words {
		words[j] = uint16(rng.Uint32())
		if i%3 == 0 {
			words[j] = 0xd800 + uint16(1024*(j%2)+rng.Intn(1024))
		} else if i%3 == 1 && j%3 != 0 {
			words[j] = 0xd800 + uint16(rng.Intn(2048))
		}
	}
	return words
}

func expected(words []uint16) string {
	text := []byte("runes")
	for _, decoded := range utf16oracle.Decode(words) {
		text = strconv.AppendInt(append(text, ' '), int64(decoded), 10)
	}
	return string(text)
}
