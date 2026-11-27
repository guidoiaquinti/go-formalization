package main

import (
	"bytes"
	"log"
	"math"
	"math/rand"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	parts []string
	sep   string
}

func main() {
	if bytes.Join([][]byte{nil}, nil) != nil {
		log.Fatal("bytes.Join of one nil slice must return nil")
	}
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
	check.CompareSamples("godiff-bytes-join", edges, sample, encodeInput, actual, expected)
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

func actual(c input) string {
	var parts [][]byte
	if c.parts != nil {
		parts = make([][]byte, len(c.parts))
		for i, part := range c.parts {
			if part != "" || i%2 == 0 {
				parts[i] = make([]byte, len(part), len(part)+4)
				copy(parts[i], part)
			}
		}
	}
	sep := []byte(c.sep)
	out := bytes.Join(parts, sep)
	if out == nil {
		return "nil"
	}
	snapshot := string(out)
	if len(parts) > 1 && cap(out) != len(out) {
		log.Fatalf("bytes.Join result capacity %d, want %d", cap(out), len(out))
	}

	// Mutate in both directions to observe whether the result aliases an input.
	if len(out) > 0 {
		out[0] ^= 0xff
		for i, part := range parts {
			if string(part) != c.parts[i] {
				log.Fatalf("bytes.Join result aliases input element %d", i)
			}
		}
		if string(sep) != c.sep {
			log.Fatal("bytes.Join result aliases separator")
		}
		out[0] ^= 0xff
	}
	for _, part := range parts {
		for i := range part {
			part[i] ^= 0xff
		}
	}
	for i := range sep {
		sep[i] ^= 0xff
	}
	if string(out) != snapshot {
		log.Fatal("mutating bytes.Join inputs modified the result")
	}
	return describe(snapshot)
}

func expected(c input) string {
	if len(c.parts) == 1 && len(c.parts[0]) == 0 {
		return "nil"
	}
	var raw []byte
	for i, part := range c.parts {
		if i > 0 {
			raw = append(raw, c.sep...)
		}
		raw = append(raw, part...)
	}
	return describe(string(raw))
}
