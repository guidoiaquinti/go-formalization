package main

import (
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	parts []string
	sep   string
}

func main() {
	var edges []input
	for _, sep := range []string{"", ",", "🚀", "\x00\xff"} {
		for _, parts := range [][]string{nil, {}, {""}, {"abc"}, {"", ""}, {"abc", ""}, {"", "é"},
			{"\xff", "\x00"}, {"", "a", "", "b", ""}, {"🚀", "é", "\xff"}} {
			edges = append(edges, input{parts, sep})
		}
		for b := range 256 {
			raw := string([]byte{byte(b)})
			edges = append(edges, input{[]string{raw, "", raw}, sep}, input{[]string{"u", "v"}, raw})
		}
	}
	check.CompareSamples("godiff-strings-join", edges, sample, encodeInput, actual, expected)
}

func sample(rng *rand.Rand, i int) input {
	n := rng.Intn(14)
	if i%4 < 3 {
		n = i % 4
	}
	parts := make([]string, n)
	for j := range parts {
		parts[j] = randomString(rng, 32)
	}
	return input{parts, randomString(rng, 8)}
}

func randomString(rng *rand.Rand, max int) string {
	raw := make([]byte, rng.Intn(max+1))
	rng.Read(raw)
	return string(raw)
}

func encodeInput(c input) string {
	row := strconv.AppendInt(nil, int64(math.MaxInt), 10)
	row = strconv.AppendInt(append(row, ' '), int64(len(c.parts)), 10)
	row = strconv.AppendInt(append(row, ' '), int64(len(c.sep)), 10)
	for _, part := range c.parts {
		row = strconv.AppendInt(append(row, ' '), int64(len(part)), 10)
	}
	for i := range len(c.sep) {
		row = strconv.AppendUint(append(row, ' '), uint64(c.sep[i]), 10)
	}
	for _, part := range c.parts {
		for i := range len(part) {
			row = strconv.AppendUint(append(row, ' '), uint64(part[i]), 10)
		}
	}
	return string(row)
}

func describe(s string) string {
	var row []byte
	for i := range len(s) {
		if i > 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(s[i]), 10)
	}
	return "bytes" + string(row)
}

func actual(c input) string { return describe(strings.Join(c.parts, c.sep)) }

func expected(c input) string {
	var raw []byte
	for i, part := range c.parts {
		if i > 0 {
			raw = append(raw, c.sep...)
		}
		raw = append(raw, part...)
	}
	return describe(string(raw))
}
