package main

import (
	"bytes"
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
	values := []string{"", "a", " alpha  beta gamma ", "alpha beta", "\u00a0\u3000", "a\u00a0b\u3000c", "\xff\uFFFD\xff", "a\xe2\x80", "a\U0001f680b"}
	for b := range 256 {
		values = append(values, "a"+string([]byte{byte(b)})+"b")
	}
	var edges []input
	for _, s := range values {
		for kind := range 5 {
			edges = append(edges, input{s, kind})
		}
	}
	check.CompareSamples("godiff-bytes-fields-func", edges, sample, encodeInput, actual, expected)
}

func sample(rng *rand.Rand, _ int) input {
	chunks := []string{" ", "\u3000", "\u00a0", "\uFFFD", "\U0001f680", "x"}
	var b strings.Builder
	for range rng.Intn(65) {
		if rng.Intn(2) == 0 {
			b.WriteByte(byte(rng.Intn(256)))
		} else {
			b.WriteString(chunks[rng.Intn(len(chunks))])
		}
	}
	return input{b.String(), rng.Intn(5)}
}

func predicate(kind int, r rune) bool {
	switch kind {
	case 0:
		return r == ' ' || r == '\u00a0' || r == '\u3000'
	case 1:
		return false
	case 2:
		return true
	case 3:
		return r == '\uFFFD'
	default:
		return r%5 == 2
	}
}

func encodeInput(c input) string {
	row := strconv.AppendInt(nil, int64(c.kind), 10)
	for i := range len(c.s) {
		row = strconv.AppendUint(append(row, ' '), uint64(c.s[i]), 10)
	}
	return string(row)
}

type fieldView struct {
	text     string
	capacity int
}

func describe(calls []rune, fields []fieldView) string {
	encodedCalls := make([]string, len(calls))
	for i, r := range calls {
		encodedCalls[i] = strconv.Itoa(int(r))
	}
	encodedFields := make([]string, len(fields))
	for i, field := range fields {
		bytes := make([]string, len(field.text))
		for j := range len(field.text) {
			bytes[j] = strconv.Itoa(int(field.text[j]))
		}
		encodedFields[i] = "f" + strconv.Itoa(field.capacity) + ":" + strings.Join(bytes, ",")
	}
	return strings.Join(encodedCalls, ",") + "|" + strings.Join(encodedFields, ";")
}

func actual(c input) string {
	var calls []rune
	input := make([]byte, len(c.s), len(c.s)+4)
	copy(input, c.s)
	fields := bytes.FieldsFunc(input, func(r rune) bool {
		calls = append(calls, r)
		return predicate(c.kind, r)
	})
	views := make([]fieldView, len(fields))
	for i, field := range fields {
		views[i] = fieldView{string(field), cap(field)}
	}
	return describe(calls, views)
}

func expected(c input) string {
	// Build substrings directly from rune-range offsets instead of saving spans.
	var calls []rune
	var fields []fieldView
	start := -1
	for offset, r := range c.s {
		calls = append(calls, r)
		if predicate(c.kind, r) {
			if start >= 0 {
				fields = append(fields, fieldView{c.s[start:offset], offset - start})
				start = -1
			}
		} else if start < 0 {
			start = offset
		}
	}
	if start >= 0 {
		fields = append(fields, fieldView{c.s[start:], len(c.s) - start})
	}
	return describe(calls, fields)
}
