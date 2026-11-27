package main

import (
	"math"
	"math/rand"
	"slices"
	"strconv"
	"unicode/utf16"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	bounds := []rune{math.MinInt32, -1, 0, 127, 128, 0xd7ff, 0xd800, 0xdbff, 0xdc00,
		0xdfff, 0xe000, 0xfffd, 0xffff, 0x10000, 0x10001, 0x10ffff, 0x110000, math.MaxInt32}
	edges := [][]rune{nil, {}}
	for _, first := range bounds {
		edges = append(edges, []rune{first})
		for _, second := range bounds {
			edges = append(edges, []rune{first, second})
		}
	}
	check.CompareSamples("godiff-utf16-encode", edges, sample, func(runes []rune) string {
		text := []byte("0")
		for _, r := range runes {
			text = strconv.AppendInt(append(text, ' '), int64(r), 10)
		}
		return string(text)
	}, func(runes []rune) string {
		copyOfRunes := slices.Clone(runes)
		words := utf16.Encode(runes)
		if !slices.Equal(runes, copyOfRunes) {
			return "input mutated"
		}
		if words == nil {
			return "nil result"
		}
		return format(words, cap(words))
	}, expected)
}

func sample(rng *rand.Rand, index int) []rune {
	runes := make([]rune, rng.Intn(17))
	for i := range runes {
		switch (index + i) % 3 {
		case 0:
			runes[i] = rune(rng.Intn(0x10000))
		case 1:
			runes[i] = 0x10000 + rune(rng.Intn(0x100000))
		default:
			runes[i] = rune(rng.Uint32())
		}
	}
	return runes
}

func expected(runes []rune) string {
	capacity := len(runes)
	var words []uint16
	for _, r := range runes {
		plane, block := r>>16, r>>10
		if plane > 0 {
			capacity++
		}
		switch {
		case plane == 0 && block != 54 && block != 55:
			words = append(words, uint16(r))
		case 1 <= plane && plane <= 16:
			payload := r - 65536
			words = append(words, uint16(55296+payload/1024), uint16(56320+payload%1024))
		default:
			words = append(words, 65533)
		}
	}
	return format(words, capacity)
}

func format(words []uint16, capacity int) string {
	text := strconv.AppendInt([]byte("present "), int64(capacity), 10)
	for _, word := range words {
		text = strconv.AppendUint(append(text, ' '), uint64(word), 10)
	}
	return string(text)
}
