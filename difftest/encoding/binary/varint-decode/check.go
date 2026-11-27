package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"slices"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	boundaries := [][]byte{nil, {}}
	for value := range 256 {
		boundaries = append(boundaries, []byte{byte(value)})
		for _, fill := range []byte{0x80, 0xff} {
			for n := 8; n <= 10; n++ {
				candidate := append(slices.Repeat([]byte{fill}, n), byte(value))
				boundaries = append(boundaries, candidate, append(slices.Clone(candidate), 0))
			}
		}
	}
	check.CompareSamples("godiff-binary-varint-decode", boundaries, sample, encode, native, oracle)
}

func sample(r *rand.Rand, trial int) []byte {
	if trial%4 == 0 {
		input := binary.AppendUvarint(nil, r.Uint64()>>uint(trial%64))
		return append(input, byte(r.Uint32()))
	}
	input := make([]byte, r.Intn(18))
	r.Read(input)
	if trial%4 == 1 {
		for i := range input {
			input[i] |= 0x80
		}
	}
	return input
}

func encode(input []byte) string {
	row := strconv.AppendInt(nil, int64(len(input)), 10)
	for _, b := range input {
		row = strconv.AppendUint(append(row, ' '), uint64(b), 10)
	}
	return string(row)
}

func native(input []byte) string {
	u, n := binary.Uvarint(input)
	s, m := binary.Varint(input)
	return fmt.Sprintf("%d:%d|%d:%d", u, n, s, m)
}

func oracle(input []byte) string {
	last := slices.IndexFunc(input[:min(10, len(input))], func(b byte) bool { return b < 128 })
	var count int
	var value uint64
	switch {
	case last < 0 && len(input) > 10:
		count = -11
	case last < 0:
		count = 0
	case last == 9 && input[last] >= 2:
		count = -10
	default:
		count = last + 1
		for index := last; index >= 0; index-- {
			value = value*128 + uint64(input[index]&127)
		}
	}
	signed := int64(value / 2)
	signed = signed - int64(value%2)*(2*signed+1)
	return fmt.Sprintf("%d:%d|%d:%d", value, count, signed, count)
}
