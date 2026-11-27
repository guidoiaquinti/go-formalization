package main

import (
	"encoding/binary"
	"fmt"
	"math/bits"
	"math/rand"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type scenario struct {
	word         uint64
	offset, size int
	backing      []byte
}

func main() {
	edges := []scenario{{}, {0, 0, 0, []byte{}}}
	words := []uint64{0, 1, 63, 64, 127, 128, 1 << 63, ^uint64(0)}
	for bit := 6; bit < 64; bit += 7 {
		boundary := uint64(1) << bit
		words = append(words, boundary-1, boundary, boundary+1, -boundary-1, -boundary, -boundary+1)
		if bit < 63 {
			boundary <<= 1
			words = append(words, boundary-1, boundary, boundary+1)
		}
	}
	for _, word := range words {
		for spare := range 12 {
			edges = append(edges, scenario{word, 1, 1, make([]byte, 2+spare)})
		}
	}
	check.CompareSamples("godiff-binary-varint-append", edges, sample, encode, native, reference)
}

func sample(r *rand.Rand, i int) scenario {
	offset, size, spare := r.Intn(8), r.Intn(12), i%12
	backing := make([]byte, offset+size+spare)
	r.Read(backing)
	word := r.Uint64()
	// Evenly exercise each encoded length, as well as all bits of uint64.
	if i%3 != 0 {
		word >>= uint(i % 64)
		if i%3 == 1 {
			word = -word
		}
	}
	return scenario{word, offset, size, backing}
}

func native(c scenario) string {
	b := slices.Clone(c.backing)
	u := binary.AppendUvarint(b[c.offset:c.offset+c.size:len(b)], c.word)
	unsigned := wire(len(b) > c.offset && &u[0] == &b[c.offset], u, b)
	b = slices.Clone(c.backing)
	s := binary.AppendVarint(b[c.offset:c.offset+c.size:len(b)], int64(c.word))
	signed := wire(len(b) > c.offset && &s[0] == &b[c.offset], s, b)
	return unsigned + "|" + signed
}

func reference(c scenario) string {
	signed := int64(c.word)
	var folded uint64
	if signed < 0 {
		folded = 2*uint64(-(signed+1)) + 1
	} else {
		folded = 2 * c.word
	}
	return referenceWord(c, c.word) + "|" + referenceWord(c, folded)
}

func referenceWord(c scenario, word uint64) string {
	width := max(1, (bits.Len64(word)+6)/7)
	payload := make([]byte, width)
	for i := range payload {
		payload[i] = byte(word>>(7*i)) & 127
		if i+1 < width {
			payload[i] += 128
		}
	}
	end := c.offset + c.size
	backing := slices.Clone(c.backing)
	copied := copy(backing[end:], payload)
	out := append(slices.Clone(c.backing[c.offset:end]), payload...)
	return wire(copied == width, out, backing)
}

func encode(c scenario) string {
	row := fmt.Appendf(nil, "%d %d %d", c.word, c.offset, c.size)
	for _, b := range c.backing {
		row = strconv.AppendUint(append(row, ' '), uint64(b), 10)
	}
	return string(row)
}

func wire(reuse bool, out, old []byte) string {
	render := func(b []byte) string {
		fields := make([]string, len(b))
		for i, v := range b {
			fields[i] = strconv.Itoa(int(v))
		}
		return strings.Join(fields, ",")
	}
	return fmt.Sprintf("%t:%s:%s", reuse, render(out), render(old))
}
