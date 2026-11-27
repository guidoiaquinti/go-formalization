package main

import (
	"math"
	"math/rand"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	dst []byte
	r   rune
}

func main() {
	var edges []input
	for _, r := range []rune{math.MinInt32, -1, 0, 0x7f, 0x80, 0x7ff, 0x800,
		0xd7ff, 0xd800, 0xdfff, 0xe000, 0xfffd, 0xffff, 0x10000,
		0x10ffff, 0x110000, math.MaxInt32} {
		edges = append(edges, input{nil, r})
		for n := range 9 {
			dst := make([]byte, n)
			for i := range dst {
				dst[i] = 0xff - byte(i)
			}
			edges = append(edges, input{dst, r})
		}
	}
	check.CompareSamples("godiff-utf8-encode-rune", edges, sample, func(in input) string {
		return field(strconv.FormatInt(int64(in.r), 10), in.dst)
	}, actual, expected)
}

func sample(rng *rand.Rand, i int) input {
	dst := make([]byte, rng.Intn(9))
	rng.Read(dst)
	r := rune(rng.Uint32())
	switch i % 4 {
	case 1:
		r = rune(rng.Intn(0x110000))
	case 2:
		r = rune(rng.Intn(0x10000))
	case 3:
		r = rune(rng.Intn(0x800))
	}
	return input{dst, r}
}

func actual(in input) (out string) {
	dst := slices.Clone(in.dst)
	defer func() {
		if p := recover(); p != nil {
			if bounds, ok := p.(runtime.Error); ok && strings.Contains(bounds.Error(), "index out of range") {
				out = field("bounds", dst)
				return
			}
			panic(p)
		}
	}()
	written := utf8.EncodeRune(dst, in.r)
	return field(strconv.Itoa(written), dst)
}

func expected(in input) string {
	dst := slices.Clone(in.dst)
	encoded := string(in.r)
	if len(encoded) > len(dst) {
		return field("bounds", dst)
	}
	copy(dst, encoded)
	return field(strconv.Itoa(len(encoded)), dst)
}

func field(status string, bytes []byte) string {
	row := []byte(status)
	for _, b := range bytes {
		row = strconv.AppendUint(append(row, ' '), uint64(b), 10)
	}
	return string(row)
}
