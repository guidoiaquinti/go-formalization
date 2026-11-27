package main

import (
	"bytes"
	"math/rand"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	text, sep string
}

func main() {
	values := []string{"", "a", "aa", "aaa", "aaaa", "aba", "ababa", "a\x00a",
		"é", "🔧", "\x80", "\xc0\x80", "\xc2", "\xe0\xa0", "\xed\xa0\x80", "\xf4\x90\x80\x80",
		"\xe0\xa0\x80", "\xed\x9f\xbf", "\xf0\x90\x80\x80", "\xf4\x8f\xbf\xbf"}
	var edges []input
	for _, text := range values {
		for _, sep := range values {
			edges = append(edges, input{text, sep})
		}
	}
	for a := range 256 {
		for b := range 256 {
			edges = append(edges, input{string([]byte{byte(a), byte(b)}), ""})
		}
		for _, sep := range []string{string([]byte{byte(a)}), string([]byte{byte(a), byte(a)})} {
			edges = append(edges, input{sep + sep + sep, sep})
		}
	}
	for _, n := range []int{32, 64, 128} {
		text := strings.Repeat("a", 4*n) + "b"
		edges = append(edges, input{text, "a"}, input{text, strings.Repeat("a", n) + "b"},
			input{text, strings.Repeat("a", n) + "c"})
	}
	check.CompareSamples("godiff-bytes-count", edges, sample, encode, func(c input) string {
		return strconv.Itoa(bytes.Count([]byte(c.text), []byte(c.sep)))
	}, func(c input) string { return strconv.Itoa(expected(c)) })
}

func sample(rng *rand.Rand, i int) input {
	raw := make([]byte, rng.Intn(65))
	rng.Read(raw)
	switch i % 5 {
	case 0:
		return input{string(raw), ""}
	case 1:
		return input{string(raw), string([]byte{byte(rng.Intn(256))})}
	case 2:
		sep := string(raw[:rng.Intn(min(len(raw), 8)+1)])
		return input{sep + string(raw) + sep, sep}
	case 3:
		return input{strings.Repeat("aba", rng.Intn(32)), "abaaba"}
	default:
		return input{string(raw), "abc"}
	}
}

func expected(c input) int {
	if len(c.sep) == 0 {
		n := 1
		for text := c.text; len(text) != 0; n++ {
			_, width := utf8.DecodeRuneInString(text)
			text = text[width:]
		}
		return n
	}
	count := 0
	for i := 0; i+len(c.sep) <= len(c.text); {
		if c.text[i:i+len(c.sep)] == c.sep {
			count++
			i += len(c.sep)
		} else {
			i++
		}
	}
	return count
}

func encode(c input) string {
	row := strconv.AppendInt(nil, int64(len(c.sep)), 10)
	for _, s := range []string{c.sep, c.text} {
		for i := range len(s) {
			row = strconv.AppendUint(append(row, ' '), uint64(s[i]), 10)
		}
	}
	return string(row)
}
