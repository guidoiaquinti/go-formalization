package main

import (
	"encoding/binary"
	"fmt"
	"math/bits"
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
		for size := range 12 {
			backing := make([]byte, 1+size+3)
			for i := range backing {
				backing[i] = 0x55
			}
			edges = append(edges, scenario{word, 1, size, backing})
		}
	}
	check.CompareSamples("godiff-binary-varint-put", edges, sample, encode, native, reference)
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
	return call(c, binary.PutUvarint) + "|" + call(c, func(b []byte, word uint64) int {
		return binary.PutVarint(b, int64(word))
	})
}

func call(c scenario, put func([]byte, uint64) int) (result string) {
	backing := slices.Clone(c.backing)
	defer func() {
		if failure := recover(); failure != nil {
			e, ok := failure.(runtime.Error)
			if !ok || !strings.HasPrefix(e.Error(), "runtime error: index out of range") {
				panic(failure)
			}
			result = wire("bounds", backing)
		}
	}()
	count := put(backing[c.offset:c.offset+c.size], c.word)
	return wire(strconv.Itoa(count), backing)
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
	copied := copy(backing[c.offset:end], payload)
	if copied < width {
		return wire("bounds", backing)
	}
	return wire(strconv.Itoa(width), backing)
}

func encode(c scenario) string {
	row := fmt.Appendf(nil, "%d %d %d", c.word, c.offset, c.size)
	for _, b := range c.backing {
		row = strconv.AppendUint(append(row, ' '), uint64(b), 10)
	}
	return string(row)
}

func wire(count string, backing []byte) string {
	fields := make([]string, len(backing))
	for i, v := range backing {
		fields[i] = strconv.Itoa(int(v))
	}
	return count + ":" + strings.Join(fields, ",")
}
