package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type replacement struct {
	pairs []string
	text  string
}

type input struct {
	mode          int
	pairs         []string
	first, second string
}

func main() {
	edges := []input{{}, {0, nil, "raw", "second"}, {1, nil, "", "a\xff"}, {1, nil, "\x00", ""}}
	for _, f := range boundaries() {
		edges = append(edges, input{2, f.pairs, f.text, f.text + "\xff"})
	}
	for b := range 256 {
		raw := string([]byte{byte(b)})
		edges = append(edges, input{2, []string{raw, "X", raw, "Y"}, raw + raw, "\x00" + raw},
			input{2, []string{raw, "", raw, "Y"}, raw + raw, "\x00" + raw})
	}
	compare("boundaries", edges)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]input, 5000)
			for i := range cases {
				f := sample(rng, i)
				second := sample(rng, i).text
				if i%2 == 0 {
					second = f.text + second
				}
				cases[i] = input{2, f.pairs, f.text, second}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases)
		}
	}
}

func boundaries() []replacement {
	var cases []replacement
	for _, pairs := range [][]string{
		nil, {"", ""}, {"", "|"}, {"", "|", "a", "X"}, {"a", "X", "", "|"},
		{"", "first", "", "second"}, {"ab", "L", "a", "S"}, {"a", "", "ab", "L"},
		{"a", "first", "a", "second"}, {"\x80", "\xff\x00", "\xff", "", "\xe2\x82", "!"},
		{"ǆ", "z", "", "."}, {"aa", "A", "aba", "B", "abaaba", "a"},
	} {
		for _, text := range []string{"", "a", "ab", "abc", "abba", "xabab", "aaaaa", "ababa", "z", "ǆ🚀", "\xff\x80\xe2\x82", "\x00\xff"} {
			cases = append(cases, replacement{pairs, text})
		}
	}
	for b := range 256 {
		key := string([]byte{byte(b)})
		text := key + "!" + key + key
		cases = append(cases, replacement{[]string{"", "", key + key, "\x00\xff", key, ""}, text},
			replacement{[]string{key, "L", key + key, "S", "", "\x00"}, text})
	}
	return cases
}

func sample(rng *rand.Rand, index int) replacement {
	alphabet := [...]int{2, 7, 256}[index%3]
	raw := func(limit int) string {
		value := make([]byte, rng.Intn(limit+1))
		for i := range value {
			value[i] = byte(rng.Intn(alphabet))
		}
		return string(value)
	}
	var pairs []string
	for range rng.Intn(12) {
		key := raw(12)
		if len(pairs) > 0 {
			base := pairs[2*rng.Intn(len(pairs)/2)]
			switch rng.Intn(4) {
			case 0:
				key = base[:rng.Intn(len(base)+1)]
			case 1:
				key = base + key
			}
		}
		pairs = append(pairs, key, raw(17))
	}
	text := raw(41)
	if len(pairs) > 0 && index%3 != 0 {
		key := pairs[2*rng.Intn(len(pairs)/2)]
		pos := rng.Intn(len(text) + 1)
		text = text[:pos] + key + text[pos:] + key
	}
	return replacement{pairs, text}
}

// Find the first eligible pair in constructor order and write directly
// to the result, consuming raw bytes or emitting one empty hit per offset.
func expected(c replacement) string {
	var out strings.Builder
	for position := 0; position <= len(c.text); {
		emittedEmpty := false
		for {
			match := -1
			for i := 0; i < len(c.pairs); i += 2 {
				key := c.pairs[i]
				if (!emittedEmpty || key != "") && strings.HasPrefix(c.text[position:], key) {
					match = i
					break
				}
			}
			if match < 0 {
				if position < len(c.text) {
					out.WriteByte(c.text[position])
				}
				position++
				break
			}
			out.WriteString(c.pairs[match+1])
			if width := len(c.pairs[match]); width != 0 {
				position += width
				break
			}
			emittedEmpty = true
		}
	}
	return out.String()
}

func observe(c input) (result string) {
	var receiver *strings.Replacer
	switch c.mode {
	case 1:
		receiver = new(strings.Replacer)
	case 2:
		receiver = strings.NewReplacer(c.pairs...)
	}
	defer func() {
		if failure := recover(); failure != nil {
			fault, ok := failure.(runtime.Error)
			if c.mode != 0 || !ok || fault.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(failure)
			}
			result = "panic"
		}
	}()
	return describe(receiver.Replace(c.first)) + "|" + describe(receiver.Replace(c.second))
}

func compare(label string, cases []input) {
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d %d", c.mode, len(c.pairs)/2)
		for i := 0; i < len(c.pairs); i += 2 {
			key, value := c.pairs[i], c.pairs[i+1]
			fmt.Fprintf(&wire, " %d %d", len(key), len(value))
			writeBytes(&wire, key)
			writeBytes(&wire, value)
		}
		fmt.Fprintf(&wire, " %d", len(c.first))
		writeBytes(&wire, c.first)
		writeBytes(&wire, c.second)
		fmt.Fprintln(&wire)
	}
	for i, model := range check.LeanRows("godiff-strings-replacer-replace", len(cases), &wire) {
		c := cases[i]
		want := "panic"
		if c.mode != 0 {
			want = describe(expected(replacement{c.pairs, c.first})) + "|" + describe(expected(replacement{c.pairs, c.second}))
		}
		if got := observe(c); got != want || model != want {
			log.Fatalf("%s row %d pairs=%q inputs=%x,%x Go=%q Lean=%q want=%q", label, i, c.pairs, c.first, c.second, got, model, want)
		}
	}
	fmt.Printf("strings.Replacer.Replace %s: %d cold/warm trace inputs agree.\n", label, len(cases))
}

func writeBytes(out *bytes.Buffer, raw string) {
	for i := range len(raw) {
		fmt.Fprintf(out, " %d", raw[i])
	}
}

func describe(raw string) string {
	values := make([]string, len(raw))
	for i := range len(raw) {
		values[i] = strconv.Itoa(int(raw[i]))
	}
	return strings.Join(values, ",")
}
