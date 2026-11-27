package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	text, old, replacement string
	limit                  int
}

func main() {
	texts := []string{"", "aaaaa", "abababa", "é🚀\x00é", "\xff\x80\xe2\x82", "\xe2\x82�",
		"\xf4\x90\x80\x80", "\xed\xa0\x80", "ééё", strings.Repeat("a", 256)}
	olds := []string{"", "a", "aa", "aba", "\xff", "\x82", "�", "é", "🚀", "\x00", "missing"}
	replacements := []string{"", "a", "aa", "é", "\xff\x00", "�"}
	limits := []int{math.MinInt, -2, -1, 0, 1, 2, 3, math.MaxInt}
	var edges []input
	for _, text := range texts {
		for _, old := range olds {
			for _, replacement := range replacements {
				for _, limit := range limits {
					edges = append(edges, input{text, old, replacement, limit})
				}
			}
		}
	}
	for text := range 63 {
		for old := range 63 {
			for _, limit := range []int{-1, 0, 1, 2, 5} {
				edges = append(edges, input{binaryBytes(text), binaryBytes(old), "\xfe\x00", limit})
			}
		}
	}
	check.CompareSamples("godiff-strings-replace", edges, sample, encode, func(c input) string {
		return describe(strings.Replace(c.text, c.old, c.replacement, c.limit))
	}, func(c input) string {
		return describe(expected(c))
	})
}

func binaryBytes(index int) string {
	var raw []byte
	for bits := index + 1; bits > 1; bits >>= 1 {
		raw = append(raw, byte(bits&1)*255)
	}
	return string(raw)
}

func expected(c input) string {
	var starts []int
	if c.old == "" {
		for start := range c.text {
			starts = append(starts, start)
		}
		starts = append(starts, len(c.text))
	} else {
		width := len(c.old)
		for i := 0; i+width <= len(c.text); {
			if c.text[i:i+width] == c.old {
				starts = append(starts, i)
				i += width
			} else {
				i++
			}
		}
	}
	if c.limit >= 0 {
		starts = starts[:min(c.limit, len(starts))]
	}
	var output strings.Builder
	consumed := 0
	for _, offset := range starts {
		output.WriteString(c.text[consumed:offset])
		output.WriteString(c.replacement)
		consumed = offset + len(c.old)
	}
	output.WriteString(c.text[consumed:])
	return output.String()
}

func raw(rng *rand.Rand, count, alphabet int) string {
	data := make([]byte, count)
	for i := range data {
		data[i] = byte(rng.Intn(alphabet))
	}
	return string(data)
}

func sample(rng *rand.Rand, index int) input {
	tokens := []string{"a", "é", "🚀", "\xff", "\x00", "�", "\xed\xa0\x80"}
	alphabet := [...]int{3, 23, 256}
	size := alphabet[index%len(alphabet)]
	text := raw(rng, rng.Intn(56), size)
	old := raw(rng, rng.Intn(11), size)
	replacement := raw(rng, rng.Intn(15), 256)
	switch index % 6 {
	case 0:
		old = ""
		for range rng.Intn(12) {
			text += tokens[rng.Intn(len(tokens))]
		}
	case 1:
		text += old + old
	case 2:
		replacement = old
	case 3:
		text = old + text + old + old
	}
	limit := rng.Intn(15) - 4
	if index%7 == 0 {
		limit = int(rng.Uint64())
	}
	return input{text, old, replacement, limit}
}

func encode(c input) string {
	row := fmt.Appendf(nil, "%d %d %d", c.limit, len(c.old), len(c.replacement))
	for _, raw := range []string{c.old, c.replacement, c.text} {
		for i := range len(raw) {
			row = strconv.AppendUint(append(row, ' '), uint64(raw[i]), 10)
		}
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
