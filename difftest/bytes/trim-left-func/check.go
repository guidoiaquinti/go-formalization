package main

import (
	"bytes"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	text        string
	kind, spare int
}

func main() {
	texts := []string{"", "a", "ab", "é", "éa🚀", "�\xffx\xe2\x82", "🚀é🚀", "xéy🚀z", "\x00", "é\xed\xa0\x80🚀"}
	for b := range 256 {
		texts = append(texts, string([]byte{byte(b)}), "é"+string([]byte{byte(b)})+"🚀")
	}
	var edges []input
	for _, text := range texts {
		for kind := range 6 {
			edges = append(edges, input{text, kind, 3})
		}
	}
	check.CompareSamples("godiff-bytes-trim-left-func", edges, sample, encode, actual, expected)
}

func sample(rng *rand.Rand, _ int) input {
	tokens := []string{"é", "🚀", "\x00", "\xff", "\xe2\x82", "�", "abc"}
	var text strings.Builder
	for range rng.Intn(33) {
		if rng.Intn(2) == 0 {
			text.WriteByte(byte(rng.Intn(256)))
		} else {
			text.WriteString(tokens[rng.Intn(len(tokens))])
		}
	}
	return input{text.String(), rng.Intn(6), rng.Intn(8)}
}

func encode(c input) string {
	row := strconv.AppendInt(nil, int64(c.kind), 10)
	row = strconv.AppendInt(append(row, ' '), int64(c.spare), 10)
	for i := range len(c.text) {
		row = strconv.AppendUint(append(row, ' '), uint64(c.text[i]), 10)
	}
	return string(row)
}

func callback(kind int) func(rune) bool {
	switch kind {
	case 0:
		return nil
	case 1:
		return func(rune) bool { return false }
	case 2:
		return func(rune) bool { return true }
	case 3:
		return func(r rune) bool { return r == utf8.RuneError }
	case 4:
		return func(r rune) bool { return r == 'é' || r == '🚀' }
	default:
		return func(r rune) bool { return r%7 == 2 }
	}
}

func actual(c input) (result string) {
	defer func() {
		if value := recover(); value != nil {
			err, ok := value.(runtime.Error)
			if !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(value)
			}
			result = "panic"
		}
	}()
	f := callback(c.kind)
	var calls []rune
	var traced func(rune) bool
	if f != nil {
		traced = func(r rune) bool {
			calls = append(calls, r)
			return f(r)
		}
	}
	var input []byte
	if len(c.text)+c.spare != 0 {
		input = make([]byte, len(c.text), len(c.text)+c.spare)
		copy(input, c.text)
	}
	trimmed := bytes.TrimLeftFunc(input, traced)
	capacity := cap(trimmed)
	if trimmed == nil {
		capacity = -1
	}
	return describe(string(trimmed), capacity, calls)
}

func expected(c input) string {
	if c.text == "" {
		return "nil|"
	}
	f := callback(c.kind)
	if f == nil {
		return "panic"
	}
	var calls []rune
	for offset, r := range c.text {
		calls = append(calls, r)
		if !f(r) {
			return describe(c.text[offset:], len(c.text)+c.spare-offset, calls)
		}
	}
	return describe("", -1, calls)
}

func describe(text string, capacity int, calls []rune) string {
	row := []byte("nil")
	if capacity >= 0 {
		row = strconv.AppendInt(nil, int64(capacity), 10)
		row = append(row, ':')
	}
	for i := range len(text) {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(text[i]), 10)
	}
	row = append(row, '|')
	for i, r := range calls {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendInt(row, int64(r), 10)
	}
	return string(row)
}
