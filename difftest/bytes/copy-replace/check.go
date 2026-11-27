package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	text, old, replacement string
	limit                  int
	nonNil                 bool
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
					edges = append(edges, input{text, old, replacement, limit, true})
				}
			}
		}
	}
	for text := range 63 {
		for old := range 63 {
			for _, limit := range []int{-1, 0, 1, 2, 5} {
				edges = append(edges, input{binaryBytes(text), binaryBytes(old), "\xfe\x00", limit, true})
			}
		}
	}
	for _, old := range olds {
		for _, replacement := range replacements {
			for _, limit := range limits {
				edges = append(edges, input{"", old, replacement, limit, false})
			}
		}
	}
	check.CompareSamples("godiff-bytes-copy-replace", edges, sample, encode, actual, expected)

}

func binaryBytes(index int) string {
	var raw []byte
	for bits := index + 1; bits > 1; bits >>= 1 {
		raw = append(raw, byte(bits&1)*255)
	}
	return string(raw)
}

func matchingStarts(text, old string) []int {
	var starts []int
	if old == "" {
		for start := range text {
			starts = append(starts, start)
		}
		return append(starts, len(text))
	}
	width := len(old)
	for i := 0; i+width <= len(text); {
		if text[i:i+width] == old {
			starts = append(starts, i)
			i += width
		} else {
			i++
		}
	}
	return starts
}

func apply(c input, starts []int, limit int) (text string, nonNil bool) {
	if limit >= 0 {
		starts = starts[:min(limit, len(starts))]
	}
	var output strings.Builder
	consumed := 0
	for _, offset := range starts {
		output.WriteString(c.text[consumed:offset])
		output.WriteString(c.replacement)
		consumed = offset + len(c.old)
	}
	output.WriteString(c.text[consumed:])
	return output.String(), len(starts) > 0 || len(c.text) > 0
}

func describe(text string, nonNil bool, appendAllocation bool) string {
	tag := "nil"
	if nonNil {
		tag = "slice"
	}
	capacity := strconv.Itoa(len(text))
	if appendAllocation && len(text) > 0 {
		capacity = "append"
	}
	data := make([]string, len(text))
	for i := range data {
		data[i] = strconv.Itoa(int(text[i]))
	}
	return fmt.Sprintf("%s:%s:%s", tag, capacity, strings.Join(data, ","))
}

func expected(c input) string {
	starts := matchingStarts(c.text, c.old)
	replaced, nonNil := apply(c, starts, c.limit)
	all, allNonNil := apply(c, starts, -1)
	return strings.Join([]string{
		describe(c.text, c.nonNil, true),
		describe(replaced, nonNil, len(starts) == 0 || c.limit == 0),
		describe(all, allNonNil, len(starts) == 0),
	}, "|")
}

func actual(c input) string {
	var data []byte
	if c.nonNil {
		data = make([]byte, len(c.text), len(c.text)+2)
		copy(data, c.text)
		copy(data[len(data):cap(data)], []byte{0xa5, 0x59})
	}
	old, replacement := []byte(c.old), []byte(c.replacement)
	// Exercise argument views that alias the source whenever their bytes fit.
	asView := func(want string, other []byte) []byte {
		if i := strings.Index(c.text, want); i >= 0 && len(want) > 0 {
			return data[i : i+len(want)]
		}
		return other
	}
	old, replacement = asView(c.old, old), asView(c.replacement, replacement)
	beforeOld, beforeReplacement := slices.Clone(old), slices.Clone(replacement)
	inspect := func(out []byte, appendAllocation bool) string {
		result := describe(string(out), out != nil, appendAllocation)
		if !appendAllocation || len(out) == 0 {
			if cap(out) != len(out) {
				panic("unexpected exact capacity")
			}
		}
		if len(out) > 0 {
			for i := range out {
				out[i] ^= 0xff
			}
		}
		if string(data) != c.text || !bytes.Equal(old, beforeOld) || !bytes.Equal(replacement, beforeReplacement) ||
			(c.nonNil && !bytes.Equal(data[len(data):cap(data)], []byte{0xa5, 0x59})) {
			panic("replacement or clone aliases its input")
		}
		return result
	}
	starts := matchingStarts(c.text, c.old)
	clone := inspect(bytes.Clone(data), true)
	replaced := inspect(bytes.Replace(data, old, replacement, c.limit), len(starts) == 0 || c.limit == 0)
	all := inspect(bytes.ReplaceAll(data, old, replacement), len(starts) == 0)
	return strings.Join([]string{clone, replaced, all}, "|")
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
	return input{text, old, replacement, limit, len(text) > 0 || index%2 != 0}
}

func encode(c input) string {
	tag := 0
	if c.nonNil {
		tag = 1
	}
	row := fmt.Appendf(nil, "%d %d %d %d %d", tag, len(c.text), len(c.old), len(c.replacement), c.limit)
	for _, part := range []string{c.text, c.old, c.replacement} {
		for i := range len(part) {
			row = strconv.AppendUint(append(row, ' '), uint64(part[i]), 10)
		}
	}
	return string(row)
}
