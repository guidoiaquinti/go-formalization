package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"runtime"
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
	edges := []scenario{{}, {0x1234, 1, 0, []byte{1, 2, 3}}, {0x89abcdef, 2, 2, []byte{3, 4, 5, 6, 7}}}
	for i := range 65536 {
		edges = append(edges, scenario{0x89abcdef01590000 | uint64(i), 1, 8, []byte{0xa5, 0, 1, 2, 3, 4, 5, 6, 7, 0xb6}})
	}
	for size := range 11 {
		for offset := range 16 {
			edges = append(edges, scenario{^uint64(0), offset, size, make([]byte, offset+size+5)})
		}
	}
	check.CompareSamples("godiff-binary-endian-put", edges, sample, encode, native, reference)
}

func sample(r *rand.Rand, i int) scenario {
	offset, size := r.Intn(16), r.Intn(20)
	if i%2 == 0 {
		size = i % 10
	}
	b := make([]byte, offset+size+r.Intn(8))
	r.Read(b)
	return scenario{r.Uint64(), offset, size, b}
}

func native(c scenario) string {
	type writeFunc func([]byte, uint64)
	writers := []writeFunc{
		func(b []byte, x uint64) { binary.BigEndian.PutUint16(b, uint16(x)) },
		func(b []byte, x uint64) { binary.BigEndian.PutUint32(b, uint32(x)) },
		binary.BigEndian.PutUint64,
		func(b []byte, x uint64) { binary.LittleEndian.PutUint16(b, uint16(x)) },
		func(b []byte, x uint64) { binary.LittleEndian.PutUint32(b, uint32(x)) },
		binary.LittleEndian.PutUint64,
	}
	results := make([]string, len(writers))
	for i, writer := range writers {
		b := slices.Clone(c.backing)
		status := call(writer, b[c.offset:c.offset+c.size], c.word)
		results[i] = wire(status, b)
	}
	return strings.Join(results, "|")
}

func call(write func([]byte, uint64), buf []byte, value uint64) (status string) {
	defer func() {
		if failure := recover(); failure != nil {
			if _, ok := failure.(runtime.Error); !ok {
				panic(failure)
			}
			status = "panic:"
		}
	}()
	write(buf, value)
	return "ok:"
}

func reference(c scenario) string {
	var results []string
	for _, big := range []bool{true, false} {
		for _, width := range []int{2, 4, 8} {
			b := slices.Clone(c.backing)
			status := "panic:"
			if c.size >= width {
				status = "ok:"
				for i := range width {
					index := i
					if big {
						index = width - 1 - i
					}
					b[c.offset+index] = byte(c.word >> (8 * i))
				}
			}
			results = append(results, wire(status, b))
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
func wire(label string, b []byte) string {
	result := make([]string, len(b))
	for i, v := range b {
		result[i] = strconv.Itoa(int(v))
	}
	return label + strings.Join(result, ",")
}
