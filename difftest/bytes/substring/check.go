package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type pair struct{ data, needle []byte }

func main() {
	values := [][]byte{nil, {}, {0}, {0, 1, 0, 1, 0}, {0xff}, {0x80, 0xc0}, []byte("é")}
	var edges []pair
	for _, text := range values {
		for _, needle := range values {
			edges = append(edges, pair{text, needle})
			padded := make([]byte, len(text), len(text)+8)
			copy(padded, text)
			edges = append(edges, pair{padded, needle})
		}
	}
	for length := range 8 {
		for flags := range 1 << length {
			data := make([]byte, length, length+3)
			for i := range data {
				data[i] = byte((flags >> i) & 1)
			}
			for _, needle := range [][]byte{{}, {0}, {1}, {0, 0}, {1, 1}, {0, 1, 0}, {1, 0, 1}} {
				edges = append(edges, pair{data, needle})
			}
		}
	}
	for _, length := range []int{2, 7, 8, 15, 16, 31, 32, 63, 64, 65, 96} {
		needle := append(bytes.Repeat([]byte("a"), length), 'b')
		data := bytes.Repeat([]byte("a"), 512)
		edges = append(edges, pair{data, needle}, pair{append(bytes.Clone(data), needle...), needle},
			pair{bytes.Repeat(needle, 2), needle})
	}
	check.CompareSamples("godiff-bytes-substring", edges, sample, encode, actual, expected)
}

func sample(rng *rand.Rand, i int) pair {
	length := rng.Intn(129)
	data := make([]byte, length, length+rng.Intn(8))
	rng.Read(data)
	if i%2 == 0 {
		for j := range data {
			data[j] &= 3
		}
	}
	begin := rng.Intn(length + 1)
	n := rng.Intn(length - begin + 1)
	needle := bytes.Clone(data[begin : begin+n])
	switch i % 4 {
	case 1:
		rng.Read(needle)
	case 2:
		needle = append(needle, byte(rng.Uint32()))
	case 3:
		needle = nil
	}
	return pair{data, needle}
}

func actual(c pair) string {
	before, after, found := bytes.Cut(c.data, c.needle)
	lastBefore, lastAfter, lastFound := bytes.CutLast(c.data, c.needle)
	return row(bytes.Index(c.data, c.needle), bytes.LastIndex(c.data, c.needle), bytes.Contains(c.data, c.needle),
		split(found, before, after), split(lastFound, lastBefore, lastAfter))
}

func expected(c pair) string {
	var hits []int
	for i := 0; i <= len(c.data)-len(c.needle); i++ {
		match := true
		for j, b := range c.needle {
			match = match && b == c.data[i+j]
		}
		if match {
			hits = append(hits, i)
		}
	}
	if len(hits) == 0 {
		missing := split(false, c.data, nil)
		return row(-1, -1, false, missing, missing)
	}
	first, last := hits[0], hits[len(hits)-1]
	return row(first, last, true, split(true, c.data[:first], c.data[first+len(c.needle):]),
		split(true, c.data[:last], c.data[last+len(c.needle):]))
}

func row(first, last int, present bool, cut, cutLast string) string {
	return fmt.Sprintf("%d|%d|%t|%s|%s", first, last, present, cut, cutLast)
}

func split(found bool, before, after []byte) string {
	return strings.Join([]string{strconv.FormatBool(found), describe(before), describe(after)}, "|")
}

func describe(data []byte) string {
	if data == nil {
		return "nil"
	}
	return string(appendBytes(strconv.AppendInt(nil, int64(cap(data)), 10), data))
}

func encode(c pair) string {
	capacity, length := cap(c.data), len(c.needle)
	if c.data == nil {
		capacity = -1
	}
	if c.needle == nil {
		length = -1
	}
	out := fmt.Appendf(nil, "%d %d", capacity, length)
	return string(appendBytes(appendBytes(out, c.needle), c.data))
}

func appendBytes(out, data []byte) []byte {
	for _, b := range data {
		out = strconv.AppendUint(append(out, ' '), uint64(b), 10)
	}
	return out
}
