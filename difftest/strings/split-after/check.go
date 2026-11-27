package main

import (
	"bytes"
	"log"
	"math/rand"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	s, sep string
}

func main() {
	var edges []input
	for _, s := range []string{"", "a", "aaaa", "baaaabaa", "éé\U0001F680", "\xff\xc0\x80\xed\xa0\x80\xe2\x80", "a\x00a"} {
		for _, sep := range []string{"", "a", "aa", "aaa", "é", "\xa9", "\xff", "\x00", "missing"} {
			edges = append(edges, input{s, sep})
		}
	}
	for b := range 256 {
		raw := string([]byte{byte(b)})
		edges = append(edges, input{raw + "a" + raw, ""}, input{raw + "a" + raw, raw})
	}
	check.CompareSamples("godiff-strings-split-after", edges, sample, encodeInput, actual, expected)
}

func sample(rng *rand.Rand, _ int) input {
	raw := make([]byte, rng.Intn(65))
	rng.Read(raw)
	start := rng.Intn(len(raw) + 1)
	end := start + rng.Intn(min(6, len(raw)-start)+1)
	sep := string(raw[start:end])
	if rng.Intn(3) == 0 {
		sep = strings.Repeat("a", rng.Intn(4))
	}
	if rng.Intn(3) == 0 {
		raw = []byte(sep + string(raw) + sep + sep)
	}
	return input{string(raw), sep}
}

func encodeInput(c input) string {
	row := strconv.AppendInt(nil, int64(len(c.sep)), 10)
	for _, s := range []string{c.sep, c.s} {
		for i := range len(s) {
			row = strconv.AppendUint(append(row, ' '), uint64(s[i]), 10)
		}
	}
	return string(row)
}

func describe(parts []string) string {
	if parts == nil {
		return "nil"
	}
	fields := make([]string, len(parts))
	for i, part := range parts {
		bytes := make([]string, len(part))
		for j := range len(part) {
			bytes[j] = strconv.Itoa(int(part[j]))
		}
		fields[i] = "f" + strings.Join(bytes, ",")
	}
	return "slice" + strings.Join(fields, ";")
}

func actual(c input) string {
	got := strings.SplitAfter(c.s, c.sep)
	if got == nil {
		log.Fatalf("strings.SplitAfter(%q, %q) returned nil", c.s, c.sep)
	}
	if rebuilt := strings.Join(got, ""); rebuilt != c.s {
		log.Fatalf("strings.SplitAfter(%q, %q) reconstructed %q", c.s, c.sep, rebuilt)
	}
	if c.sep == "" {
		split := strings.Split(c.s, c.sep)
		if split == nil || !slices.Equal(got, split) {
			log.Fatalf("strings.SplitAfter(%q, %q)=%q, Split=%q", c.s, c.sep, got, split)
		}
	}
	return describe(got)
}

func expected(c input) string {
	pieces := bytes.SplitAfter([]byte(c.s), []byte(c.sep))
	parts := make([]string, len(pieces))
	for i, piece := range pieces {
		parts[i] = string(piece)
	}
	return describe(parts)
}
