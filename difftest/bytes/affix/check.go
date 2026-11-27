package main

import (
	"bytes"
	"math/rand"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	src, affix []byte
}

func main() {
	values := [][]byte{nil, {}, {0}, {1}, {0, 1}, {1, 0}, {0, 0}, {0xff}, {0xc0, 0x80}}
	var edges []input
	for _, src := range values {
		for _, affix := range values {
			edges = append(edges, input{src, affix})
			padded := make([]byte, len(src), len(src)+3)
			copy(padded, src)
			edges = append(edges, input{padded, affix})
		}
	}
	for b := range 256 {
		edges = append(edges, input{[]byte{byte(b), 0, byte(b)}, []byte{byte(b)}})
	}
	check.CompareSamples("godiff-bytes-affix", edges, sample, encode, actual, expected)
}

func sample(rng *rand.Rand, i int) input {
	length := rng.Intn(33)
	src := make([]byte, length, length+rng.Intn(9))
	rng.Read(src)
	n := rng.Intn(length + 1)
	offset := 0
	if i%2 != 0 {
		offset = length - n
	}
	affix := append([]byte(nil), src[offset:offset+n]...)
	switch (i / 2) % 4 {
	case 1:
		affix = append(affix, byte(rng.Intn(256)))
	case 2:
		if n > 0 {
			affix[rng.Intn(n)] ^= 1
		}
	case 3:
		rng.Read(affix)
	}
	return input{src, affix}
}

func encode(c input) string {
	capacity := cap(c.src)
	if c.src == nil {
		capacity = -1
	}
	row := strconv.AppendInt(nil, int64(capacity), 10)
	row = strconv.AppendInt(append(row, ' '), int64(len(c.affix)), 10)
	return string(appendBytes(appendBytes(row, c.affix), c.src))
}

func actual(c input) string {
	pre, post := bytes.HasPrefix(c.src, c.affix), bytes.HasSuffix(c.src, c.affix)
	leading, leadingFound := bytes.CutPrefix(c.src, c.affix)
	trailing, trailingFound := bytes.CutSuffix(c.src, c.affix)
	return results(pre, post, leadingFound, trailingFound, leading, trailing,
		bytes.TrimPrefix(c.src, c.affix), bytes.TrimSuffix(c.src, c.affix))
}

func expected(c input) string {
	src, affix := string(c.src), string(c.affix)
	var pre, post bool
	leading, trailing := c.src, c.src
	for n := range len(src) + 1 {
		if affix+src[n:] == src {
			pre = true
			leading = c.src[n:]
		}
		if src[:n]+affix == src {
			post = true
			trailing = c.src[:n]
		}
	}
	return results(pre, post, pre, post, leading, trailing, leading, trailing)
}

func results(pre, post, leadingFound, trailingFound bool, leading, trailing, trimPre, trimPost []byte) string {
	return strings.Join([]string{
		strconv.FormatBool(pre), strconv.FormatBool(post),
		strconv.FormatBool(leadingFound) + " " + describe(leading),
		strconv.FormatBool(trailingFound) + " " + describe(trailing),
		describe(trimPre), describe(trimPost),
	}, "|")
}

func describe(data []byte) string {
	if data == nil {
		return "nil"
	}
	return string(appendBytes(strconv.AppendInt(nil, int64(cap(data)), 10), data))
}

func appendBytes(out, data []byte) []byte {
	for _, b := range data {
		out = strconv.AppendUint(append(out, ' '), uint64(b), 10)
	}
	return out
}
