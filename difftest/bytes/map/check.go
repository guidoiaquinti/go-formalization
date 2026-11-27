package main

import (
	"bytes"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	s    string
	kind int
}

func main() {
	values := []string{"", "same ascii", "é€🚀", "\x00", "\xff�\x80", "valid�", "x\xc2\xa0", "a\xe2\x82", "\xed\xa0\x80",
		"\xf4\x90\x80\x80", strings.Repeat("a", 64) + "\xff", "\xff" + strings.Repeat("🚀", 32)}
	for b := range 256 {
		values = append(values, string([]byte{byte(b)}), string([]byte{'a', byte(b), 0x80, 'z'}))
	}
	var edges []input
	for _, s := range values {
		for kind := range 6 {
			edges = append(edges, input{s, kind})
		}
	}
	check.CompareSamples("godiff-bytes-map", edges, sample, encodeInput, actual, expected)
}

func sample(rng *rand.Rand, _ int) input {
	chunks := []string{"é", "🚀", "�", "\xf0\x9f", "\xff", "\x80", "abc"}
	var s strings.Builder
	for range rng.Intn(49) {
		if rng.Intn(2) == 0 {
			s.WriteByte(byte(rng.Intn(256)))
		} else {
			s.WriteString(chunks[rng.Intn(len(chunks))])
		}
	}
	return input{s.String(), rng.Intn(6)}
}

func mapping(kind int, r rune) rune {
	switch kind {
	case 0:
		return r
	case 1:
		return -1
	case 2:
		if r == '\uFFFD' {
			return 'x'
		}
		return r
	case 3:
		if r%2 == 0 {
			return 0x110000
		}
		return 0xd800
	case 4:
		if r%2 == 0 {
			return '🚀'
		}
		return 0
	default:
		if r%3 == 0 {
			return math.MinInt32
		}
		return math.MaxInt32
	}
}

func encodeInput(c input) string {
	row := strconv.AppendInt(nil, int64(c.kind), 10)
	for i := range len(c.s) {
		row = strconv.AppendUint(append(row, ' '), uint64(c.s[i]), 10)
	}
	return string(row)
}

func describe(calls []rune, out string) string {
	row := []byte{}
	for i, r := range calls {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendInt(row, int64(r), 10)
	}
	row = append(row, '|')
	row = append(row, "bytes"...)
	for i := range len(out) {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(out[i]), 10)
	}
	return string(row)
}

func actual(c input) string {
	var calls []rune
	var src []byte
	if c.s != "" {
		src = []byte(c.s)
	}
	result := bytes.Map(func(r rune) rune {
		calls = append(calls, r)
		return mapping(c.kind, r)
	}, src)
	if result == nil {
		return "nil"
	}
	return describe(calls, string(result))
}

func expected(c input) string {
	var calls, output []rune
	for _, r := range c.s {
		calls = append(calls, r)
		if mapped := mapping(c.kind, r); mapped >= 0 {
			output = append(output, mapped)
		}
	}
	return describe(calls, string(output))
}
