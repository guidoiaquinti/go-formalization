package main

import (
	"encoding/binary"
	"fmt"
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
	edges := []scenario{{}, {0, 0, 0, []byte{}}, {0x1234, 1, 0, []byte{1, 2, 3}}}
	for i := range 65536 {
		edges = append(edges, scenario{0x89abcdef01590000 | uint64(i), 1, 1, []byte{0xa5, 0, 1, 2, 3, 4, 5, 6, 7, 0xb6}})
	}
	for size := range 11 {
		for spare := range 10 {
			edges = append(edges, scenario{^uint64(0), 1, size, make([]byte, 1+size+spare)})
		}
	}
	check.CompareSamples("godiff-binary-endian-append", edges, sample, encode, native, reference)
}

func sample(r *rand.Rand, i int) scenario {
	offset, size := r.Intn(16), r.Intn(20)
	spare := r.Intn(14)
	if i%2 == 0 {
		spare = i % 10
	}
	backing := make([]byte, offset+size+spare)
	r.Read(backing)
	return scenario{r.Uint64(), offset, size, backing}
}

func native(c scenario) string {
	type appendFunc func([]byte, uint64) []byte
	writers := []appendFunc{
		func(b []byte, x uint64) []byte { return binary.BigEndian.AppendUint16(b, uint16(x)) },
		func(b []byte, x uint64) []byte { return binary.BigEndian.AppendUint32(b, uint32(x)) },
		binary.BigEndian.AppendUint64,
		func(b []byte, x uint64) []byte { return binary.LittleEndian.AppendUint16(b, uint16(x)) },
		func(b []byte, x uint64) []byte { return binary.LittleEndian.AppendUint32(b, uint32(x)) },
		binary.LittleEndian.AppendUint64,
	}
	results := make([]string, len(writers))
	for i, writer := range writers {
		b := slices.Clone(c.backing)
		out := writer(b[c.offset:c.offset+c.size:len(b)], c.word)
		reused := len(b) > c.offset && &out[0] == &b[c.offset]
		results[i] = wire(reused, out, b)
	}
	return strings.Join(results, "|")
}

func reference(c scenario) string {
	var results []string
	for _, big := range []bool{true, false} {
		for _, width := range []int{2, 4, 8} {
			b := slices.Clone(c.backing)
			end := c.offset + c.size
			word := make([]byte, width)
			for i := range width {
				index := i
				if big {
					index = width - 1 - i
				}
				word[index] = byte(c.word >> (8 * i))
			}
			reuse := len(b)-end >= width
			if reuse {
				copy(b[end:end+width], word)
			}
			out := append(slices.Clone(c.backing[c.offset:end]), word...)
			results = append(results, wire(reuse, out, b))
		}
	}
	return strings.Join(results, "|")
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
