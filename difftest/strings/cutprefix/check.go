package main

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	text, prefix string
}

func main() {
	values := []string{"", "a", "ab", "abc", "aaa", "a\x00", "\x00", "\xff", "\xc0\x80", "é", "🔧"}
	var edges []input
	for _, text := range values {
		for _, prefix := range values {
			edges = append(edges, input{text, prefix})
		}
	}
	for i := range 256 {
		b := string([]byte{byte(i)})
		edges = append(edges, input{b + "!", b})
	}
	check.CompareSamples("godiff-strings-cutprefix", edges, sample, func(c input) string {
		row := appendBytes(strconv.AppendInt(nil, int64(len(c.prefix)), 10), c.prefix)
		return string(appendBytes(row, c.text))
	}, func(c input) string {
		return result(strings.CutPrefix(c.text, c.prefix))
	}, func(c input) string {
		for i := range len(c.prefix) {
			if i >= len(c.text) || c.text[i] != c.prefix[i] {
				return result(c.text, false)
			}
		}
		return result(c.text[len(c.prefix):], true)
	})
}

func sample(rng *rand.Rand, i int) input {
	data := make([]byte, rng.Intn(33))
	rng.Read(data)
	length := rng.Intn(len(data) + 1)
	prefix := append([]byte(nil), data[:length]...)
	switch i % 4 {
	case 1:
		prefix = append(prefix, byte(rng.Intn(256)))
	case 2:
		if length > 0 {
			prefix[rng.Intn(length)] ^= 1
		}
	case 3:
		rng.Read(prefix)
	}
	return input{string(data), string(prefix)}
}

func result(text string, found bool) string {
	return string(appendBytes(strconv.AppendBool(nil, found), text))
}

func appendBytes(out []byte, text string) []byte {
	for i := range len(text) {
		out = strconv.AppendUint(append(out, ' '), uint64(text[i]), 10)
	}
	return out
}
