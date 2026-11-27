package main

import (
	"encoding/binary"
	"math/rand"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	edges := [][]byte{nil, {}, {1}, {0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0xff}}
	// Exhaust the complete 16-bit word domain using sufficient input for
	// the wider accessors as well, with distinct trailing bytes.
	for word := range 65536 {
		edges = append(edges, []byte{byte(word >> 8), byte(word), 0x81, 0x27, 0xe4, 0x5a, 0xc0, 0x6c, 0x13})
	}
	for width := range 10 {
		for b := range 256 {
			word := make([]byte, width)
			for i := range word {
				word[i] = byte(b + i)
			}
			edges = append(edges, word)
		}
	}
	check.CompareSamples("godiff-binary-endian-read", edges, sample, encode, native, expected)
}

func sample(r *rand.Rand, i int) []byte {
	size := r.Intn(25)
	if i%4 == 0 {
		size = i % 9
	}
	offset := r.Intn(16)
	raw := make([]byte, size+offset+8)
	r.Read(raw)
	return raw[offset : offset+size]
}

func encode(input []byte) string {
	row := []byte{'0'}
	for _, b := range input {
		row = strconv.AppendUint(append(row, ' '), uint64(b), 10)
	}
	return string(row)
}

func native(input []byte) string {
	return strings.Join([]string{
		call(func() uint64 { return uint64(binary.BigEndian.Uint16(input)) }),
		call(func() uint64 { return uint64(binary.BigEndian.Uint32(input)) }),
		call(func() uint64 { return binary.BigEndian.Uint64(input) }),
		call(func() uint64 { return uint64(binary.LittleEndian.Uint16(input)) }),
		call(func() uint64 { return uint64(binary.LittleEndian.Uint32(input)) }),
		call(func() uint64 { return binary.LittleEndian.Uint64(input) }),
	}, "|")
}

func call(read func() uint64) (result string) {
	defer func() {
		if failure := recover(); failure != nil {
			if _, ok := failure.(runtime.Error); !ok {
				panic(failure)
			}
			result = "panic"
		}
	}()
	return strconv.FormatUint(read(), 10)
}

func expected(input []byte) string {
	var result []string
	for _, big := range []bool{true, false} {
		for _, width := range []int{2, 4, 8} {
			if len(input) < width {
				result = append(result, "panic")
				continue
			}
			var word uint64
			for i := range width {
				shift := 8 * i
				if big {
					shift = 8 * (width - i - 1)
				}
				word |= uint64(input[i]) << shift
			}
			result = append(result, strconv.FormatUint(word, 10))
		}
	}
	return strings.Join(result, "|")
}
