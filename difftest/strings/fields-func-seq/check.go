package main

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	s            string
	kind, budget int
}

type event struct {
	yield bool
	rune  rune
	field string
}

func main() {
	values := []string{"", "a", " alpha  beta gamma ", "alpha beta", "\u00a0\u3000", "a\u00a0b\u3000c", "\xff\uFFFD\xff", "a\xe2\x80", "a\U0001f680b"}
	for b := range 256 {
		values = append(values, "a"+string([]byte{byte(b)})+"b")
	}
	var edges []input
	for _, s := range values {
		for kind := range 5 {
			for _, budget := range []int{0, 1, 2, len(s) + 1} {
				edges = append(edges, input{s, kind, budget})
			}
		}
	}
	check.CompareSamples("godiff-strings-fields-func-seq", edges, sample, encodeInput, actual, expected)
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
	return input{b.String(), rng.Intn(5), rng.Intn(b.Len() + 2)}
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
	row = strconv.AppendInt(append(row, ' '), int64(c.budget), 10)
	for i := range len(c.s) {
		row = strconv.AppendUint(append(row, ' '), uint64(c.s[i]), 10)
	}
	return string(row)
}

func describe(events []event) string {
	encoded := make([]string, len(events))
	for i, e := range events {
		if !e.yield {
			encoded[i] = "p" + strconv.Itoa(int(e.rune))
			continue
		}
		bytes := make([]string, len(e.field))
		for j := range len(e.field) {
			bytes[j] = strconv.Itoa(int(e.field[j]))
		}
		encoded[i] = "y" + strings.Join(bytes, ",")
	}
	return strings.Join(encoded, ";")
}

func actual(c input) string {
	var events []event
	seq := strings.FieldsFuncSeq(c.s, func(r rune) bool {
		events = append(events, event{rune: r})
		return predicate(c.kind, r)
	})
	left := c.budget
	if left > 0 {
		seq(func(s string) bool {
			events = append(events, event{yield: true, field: s})
			left--
			return left != 0
		})
	}
	limited := describe(events)
	events = nil
	seq(func(s string) bool {
		events = append(events, event{yield: true, field: s})
		return true
	})
	return limited + "|" + describe(events)
}

func expected(c input) string {
	// Range supplies independently decoded runes and byte offsets. Build
	// the complete invocation before selecting the bounded trace.
	var events []event
	start := -1
	for offset, r := range c.s {
		events = append(events, event{rune: r})
		if predicate(c.kind, r) {
			if start >= 0 {
				events = append(events, event{yield: true, field: c.s[start:offset]})
				start = -1
			}
		} else if start < 0 {
			start = offset
		}
	}
	if start >= 0 {
		events = append(events, event{yield: true, field: c.s[start:]})
	}
	end := 0
	for left := c.budget; end < len(events) && left > 0; end++ {
		if events[end].yield {
			left--
		}
	}
	return describe(events[:end]) + "|" + describe(events)
}
