package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type pair struct{ text, needle string }

func main() {
	values := []string{"", "a", "aaa", "abaa", "é", "\x00", "\xff", "\xc0\x80", "01234567"}
	var edges []pair
	for _, text := range values {
		for _, needle := range values {
			edges = append(edges, pair{text, needle})
		}
	}
	for length := range 8 {
		for bits := range 1 << length {
			var text []byte
			for i := range length {
				text = append(text, byte('a'+((bits>>i)&1)))
			}
			for _, needle := range []string{"", "a", "b", "aa", "ab", "ba", "bb", "aba", "bab", "aaaa"} {
				edges = append(edges, pair{string(text), needle})
			}
		}
	}
	for _, n := range []int{2, 7, 8, 15, 16, 31, 32, 63, 64, 65, 96} {
		pattern := strings.Repeat("a", n) + "b"
		text := strings.Repeat("a", 512)
		edges = append(edges, pair{text, pattern}, pair{text + pattern, pattern},
			pair{pattern + text, pattern}, pair{pattern + pattern, pattern})
	}
	check.CompareSamples("godiff-strings-substring", edges, sample, func(c pair) string {
		out := strconv.AppendInt(nil, int64(len(c.needle)), 10)
		return string(appendBytes(appendBytes(out, c.needle), c.text))
	}, func(c pair) string {
		before, after, found := strings.Cut(c.text, c.needle)
		lastBefore, lastAfter, lastFound := strings.CutLast(c.text, c.needle)
		return result(strings.Index(c.text, c.needle), strings.LastIndex(c.text, c.needle), strings.Contains(c.text, c.needle),
			found, before, after, lastFound, lastBefore, lastAfter)
	}, expected)
}

func sample(rng *rand.Rand, index int) pair {
	text := make([]byte, rng.Intn(129))
	rng.Read(text)
	if index%2 == 0 {
		for i := range text {
			text[i] &= 3
		}
	}
	start := rng.Intn(len(text) + 1)
	length := rng.Intn(len(text) - start + 1)
	needle := append([]byte(nil), text[start:start+length]...)
	switch index % 4 {
	case 1:
		rng.Read(needle)
	case 2:
		needle = append(needle, byte(rng.Intn(256)))
	case 3:
		needle = nil
	}
	return pair{string(text), string(needle)}
}

func expected(c pair) string {
	first, last := -1, -1
	// Scan backward, so the lowest matching offset is saved last.
	for i := len(c.text) - len(c.needle); i >= 0; i-- {
		match := true
		for j := range len(c.needle) {
			match = match && c.text[i+j] == c.needle[j]
		}
		if match {
			last = max(last, i)
			first = i
		}
	}
	if first < 0 {
		return result(-1, -1, false, false, c.text, "", false, c.text, "")
	}
	return result(first, last, true, true, c.text[:first], c.text[first+len(c.needle):],
		true, c.text[:last], c.text[last+len(c.needle):])
}

func result(index, last int, contains, found bool, before, after string, lastFound bool, lastBefore, lastAfter string) string {
	out := fmt.Appendf(nil, "%d %d %t", index, last, contains)
	out = appendSplit(out, found, before, after)
	return string(appendSplit(out, lastFound, lastBefore, lastAfter))
}

func appendSplit(out []byte, found bool, before, after string) []byte {
	out = fmt.Appendf(out, " %t %d", found, len(before))
	out = appendBytes(out, before)
	out = strconv.AppendInt(append(out, ' '), int64(len(after)), 10)
	return appendBytes(out, after)
}

func appendBytes(out []byte, text string) []byte {
	for i := range len(text) {
		out = strconv.AppendUint(append(out, ' '), uint64(text[i]), 10)
	}
	return out
}
