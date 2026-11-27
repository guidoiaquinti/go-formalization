package main

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	text, suffix string
}

func main() {
	values := []string{"", "a", "ab", "abc", "aaa", "a\x00", "\x00", "\xff", "\xc0\x80", "é", "🔧"}
	var edges []input
	for _, text := range values {
		for _, suffix := range values {
			edges = append(edges, input{text, suffix})
		}
	}
	for i := range 256 {
		b := string([]byte{byte(i)})
		edges = append(edges, input{"!" + b, b})
	}
	check.CompareSamples("godiff-strings-cutsuffix", edges, sample, func(c input) string {
		row := appendBytes(strconv.AppendInt(nil, int64(len(c.suffix)), 10), c.suffix)
		return string(appendBytes(row, c.text))
	}, func(c input) string {
		return result(strings.CutSuffix(c.text, c.suffix))
	}, func(c input) string {
		before := len(c.text) - len(c.suffix)
		for i := 1; i <= len(c.suffix); i++ {
			if i > len(c.text) || c.text[len(c.text)-i] != c.suffix[len(c.suffix)-i] {
				return result(c.text, false)
			}
		}
		return result(c.text[:before], true)
	})
}

func sample(rng *rand.Rand, i int) input {
	data := make([]byte, rng.Intn(33))
	rng.Read(data)
	length := rng.Intn(len(data) + 1)
	suffix := append([]byte(nil), data[len(data)-length:]...)
	switch i % 4 {
	case 1:
		suffix = append([]byte{byte(rng.Intn(256))}, suffix...)
	case 2:
		if length > 0 {
			suffix[rng.Intn(length)] ^= 1
		}
	case 3:
		rng.Read(suffix)
	}
	return input{string(data), string(suffix)}
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
