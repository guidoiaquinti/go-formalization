package main

import (
	"math"
	"math/rand"
	"slices"
	"strconv"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	dst           []byte
	r             rune
	extraCapacity int
}

func main() {
	var edges []input
	for _, dst := range [][]byte{nil, {}, {0}, {0xff, '\n', 'x'}} {
		for _, r := range []rune{math.MinInt32, -1, 0, 0x7f, 0x80, 0x7ff, 0x800,
			0xd7ff, 0xd800, 0xdfff, 0xe000, 0xfffd, 0xffff, 0x10000,
			0x10ffff, 0x110000, math.MaxInt32} {
			for _, extra := range []int{0, 1, 2, 3, 4} {
				edges = append(edges, input{dst, r, extra})
			}
		}
	}
	check.CompareSamples("godiff-utf8-append-rune", edges, sample, func(in input) string {
		return strconv.FormatInt(int64(in.r), 10) + " " + encode(in.dst)
	}, func(in input) string {
		dst := slices.Clone(in.dst)
		if in.extraCapacity > 0 {
			dst = make([]byte, len(in.dst), len(in.dst)+in.extraCapacity)
			copy(dst, in.dst)
		}
		return encode(utf8.AppendRune(dst, in.r))
	}, func(in input) string {
		// The language's rune-to-string conversion supplies an independent
		// encoding and replacement rule.
		return encode([]byte(string(in.dst) + string(in.r)))
	})
}

func sample(rng *rand.Rand, i int) input {
	dst := make([]byte, rng.Intn(17))
	rng.Read(dst)
	r := rune(rng.Uint32())
	switch i % 4 {
	case 1:
		r = rune(rng.Intn(0x110000))
	case 2:
		r = rune(rng.Intn(0x10000))
	case 3:
		r = rune(rng.Intn(128))
	}
	return input{dst, r, rng.Intn(8)}
}

func encode(bytes []byte) string {
	var row []byte
	for i, b := range bytes {
		if i != 0 {
			row = append(row, ' ')
		}
		row = strconv.AppendUint(row, uint64(b), 10)
	}
	return string(row)
}
