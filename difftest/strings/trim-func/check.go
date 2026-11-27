package main

import (
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	text string
	kind int
}

func main() {
	texts := []string{"", "a", "ab", "é", "éa🚀", "�\xffx\xe2\x82", "🚀é🚀", "xéy🚀z", "\x00", "é\xed\xa0\x80🚀"}
	for b := range 256 {
		texts = append(texts, string([]byte{byte(b)}), "é"+string([]byte{byte(b)})+"🚀")
	}
	var edges []input
	for _, text := range texts {
		for kind := range 6 {
			edges = append(edges, input{text, kind})
		}
	}
	check.CompareSamples("godiff-strings-trim-func", edges, sample, encode, actual, expected)
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
	return input{text.String(), rng.Intn(6)}
}

func encode(c input) string {
	row := strconv.AppendInt(nil, int64(c.kind), 10)
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
	trimmed := strings.TrimFunc(c.text, traced)
	return describe(trimmed, calls)
}

func expected(c input) string {
	if c.text == "" {
		return "|"
	}
	f := callback(c.kind)
	if f == nil {
		return "panic"
	}
	var positions []int
	var runes, calls []rune
	for i, r := range c.text {
		positions = append(positions, i)
		runes = append(runes, r)
	}
	positions = append(positions, len(c.text))
	for left, r := range runes {
		calls = append(calls, r)
		if f(r) {
			continue
		}
		for right := len(runes) - 1; right >= left; right-- {
			calls = append(calls, runes[right])
			if !f(runes[right]) {
				return describe(c.text[positions[left]:positions[right+1]], calls)
			}
		}
	}
	return describe("", calls)
}

func describe(text string, calls []rune) string {
	var row []byte
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
