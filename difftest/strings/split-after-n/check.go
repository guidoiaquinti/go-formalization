package main

import (
	"bytes"
	"log"
	"math"
	"math/rand"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	s, sep string
	limit  int
}

func main() {
	var edges []input
	add := func(s, sep string) {
		for _, n := range []int{math.MinInt, -4, -1, 0, 1, 2, 3, len(s), len(s) + 1, math.MaxInt} {
			edges = append(edges, input{s, sep, n})
		}
	}
	for _, s := range []string{"", "a", "aaaa", "baaaabaa", "éé\U0001F680", "\xff\xc0\x80\xed\xa0\x80\xe2\x80", "a\x00a"} {
		for _, sep := range []string{"", "a", "aa", "aaa", "é", "\xa9", "\xff", "\x00", "missing"} {
			add(s, sep)
		}
	}
	for b := range 256 {
		raw := string([]byte{byte(b)})
		add(raw+"a"+raw, "")
		add(raw+"a"+raw, raw)
	}
	check.CompareSamples("godiff-strings-split-after-n", edges, sample, encodeInput, actual, expected)
}

func sample(rng *rand.Rand, i int) input {
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
	limits := [...]int{-1, 0, 1, 2, rng.Intn(len(raw) + 3), math.MaxInt}
	return input{string(raw), sep, limits[i%len(limits)]}
}

func encodeInput(c input) string {
	row := strconv.AppendInt(nil, int64(c.limit), 10)
	row = strconv.AppendInt(append(row, ' '), int64(len(c.sep)), 10)
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
	got := strings.SplitAfterN(c.s, c.sep, c.limit)
	if c.limit != 0 && (got == nil || strings.Join(got, "") != c.s) {
		log.Fatalf("strings.SplitAfterN(%q, %q, %d) returned %q without reconstructing the input", c.s, c.sep, c.limit, got)
	}
	if c.sep == "" {
		split := strings.SplitN(c.s, "", c.limit)
		if (split == nil) != (got == nil) || !slices.Equal(got, split) {
			log.Fatalf("strings.SplitAfterN(%q, empty, %d)=%q, SplitN=%q", c.s, c.limit, got, split)
		}
	}
	if c.limit == -1 {
		unbounded := strings.SplitAfter(c.s, c.sep)
		if unbounded == nil || !slices.Equal(got, unbounded) {
			log.Fatalf("strings.SplitAfterN(%q, %q, -1)=%q, SplitAfter=%q", c.s, c.sep, got, unbounded)
		}
	}
	return describe(got)
}

func expected(c input) string {
	if c.limit == 0 {
		return "nil"
	}
	pieces := bytes.SplitAfter([]byte(c.s), []byte(c.sep))
	parts := make([]string, len(pieces))
	for i, piece := range pieces {
		parts[i] = string(piece)
	}
	if c.limit > 0 && c.limit < len(parts) {
		parts[c.limit-1] = strings.Join(parts[c.limit-1:], "")
		parts = parts[:c.limit]
	}
	return describe(parts)
}
