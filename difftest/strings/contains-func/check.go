package main

import (
	"math/rand"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	s    string
	kind int
}

func main() {
	texts := []string{"", "abc", "é🚀🚀", "abc🚀", "�x", "a\xff�", "a\xe2\x82", "\xf4\x90\x80\x80", "\x00"}
	for b := range 256 {
		texts = append(texts, string([]byte{byte(b)}), string([]byte{'a', byte(b), 0x80, 'z'}))
	}
	var edges []input
	for _, s := range texts {
		for kind := range 6 {
			edges = append(edges, input{s, kind})
		}
	}
	check.CompareSamples("godiff-strings-contains-func", edges, sample, encodeInput, actual, expected)
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

func predicate(kind int) func(rune) bool {
	switch kind {
	case 0:
		return nil
	case 1:
		return func(rune) bool { return false }
	case 2:
		return func(rune) bool { return true }
	case 3:
		return func(r rune) bool { return r == '\uFFFD' }
	case 4:
		return func(r rune) bool { return r == '🚀' }
	default:
		return func(r rune) bool { return r%7 == 2 }
	}
}

func encodeInput(c input) string {
	row := strconv.AppendInt(nil, int64(c.kind), 10)
	for i := range len(c.s) {
		row = strconv.AppendUint(append(row, ' '), uint64(c.s[i]), 10)
	}
	return string(row)
}

func describe(found bool, calls []rune) string {
	row := strconv.AppendBool(nil, found)
	row = append(row, '|')
	for i, r := range calls {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendInt(row, int64(r), 10)
	}
	return string(row)
}

func actual(c input) (result string) {
	checkRune := predicate(c.kind)
	var calls []rune
	if checkRune != nil {
		original := checkRune
		checkRune = func(r rune) bool {
			calls = append(calls, r)
			return original(r)
		}
	}
	defer func() {
		if p := recover(); p != nil {
			if _, ok := p.(runtime.Error); c.kind == 0 && ok {
				result = "panic"
			} else {
				panic(p)
			}
		}
	}()
	found := strings.ContainsFunc(c.s, checkRune)
	return describe(found, calls)
}

func expected(c input) string {
	checkRune := predicate(c.kind)
	if checkRune == nil && c.s != "" {
		return "panic"
	}
	var calls []rune
	for _, r := range c.s {
		calls = append(calls, r)
		if checkRune(r) {
			return describe(true, calls)
		}
	}
	return describe(false, calls)
}
