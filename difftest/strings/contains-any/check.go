package main

import (
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
	tokens := []string{"", "a", "\x00", "\x7f", "\x80", "\xff", "\xc3", "é", "e\u0301", "\ufffd", "😀", "\xf0\x9f", "\xed\xa0\x80", "\xf4\x90\x80\x80"}
	for _, text := range tokens {
		for _, chars := range tokens {
			edges = append(edges, input{text, chars}, input{"zzzzzzzz" + text, chars},
				input{"zzzzzzzzz" + text, chars + "y"})
		}
	}
	check.CompareSamples("godiff-strings-contains-any", edges, sample, func(c input) string {
		row := strconv.AppendInt(nil, int64(len(c.text)), 10)
		for _, part := range []string{c.text, c.chars} {
			for i := range len(part) {
				row = strconv.AppendUint(append(row, ' '), uint64(part[i]), 10)
			}
		}
		return string(row)
	}, func(c input) string {
		return strconv.FormatBool(strings.ContainsAny(c.text, c.chars))
	}, func(c input) string {
		text := []rune(c.text)
		for _, key := range c.chars {
			for _, r := range text {
				if r == key {
					return "true"
				}
			}
		}
		return "false"
	})
}

func sample(rng *rand.Rand, i int) input {
	text := make([]byte, rng.Intn(97))
	chars := make([]byte, rng.Intn(21))
	rng.Read(text)
	rng.Read(chars)
	switch i % 4 {
	case 0:
		for j := range chars {
			chars[j] &= 0x7f
		}
	case 1:
		chars = []byte{byte(rng.Intn(256))}
	case 2:
		chars = append(chars, "😀"...)
		text = append(text, "😀"...)
	}
	return input{string(text), string(chars)}
}
