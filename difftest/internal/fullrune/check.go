package fullrune

import (
	"math/rand"
	"strconv"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func Compare(executable string, actual func([]byte) bool) {
	edges := [][]byte{nil, {}}
	for a := range 256 {
		edges = append(edges, []byte{byte(a)})
		for b := range 256 {
			edges = append(edges, []byte{byte(a), byte(b)})
		}
	}
	continuations := []byte{0, 0x7f, 0x80, 0x8f, 0x90, 0x9f, 0xa0, 0xbf, 0xc0, 0xff}
	for a := byte(0xf0); a <= 0xf4; a++ {
		for _, b := range continuations {
			for _, c := range continuations {
				edges = append(edges, []byte{a, b, c})
			}
		}
	}
	edges = append(edges, []byte{0xf0, 0x90, 0x80, 0x80}, []byte{0xf4, 0x90, 0x80, 0x80})
	check.CompareSamples(executable, edges, sample, encode, func(data []byte) string {
		return strconv.FormatBool(actual(data))
	}, func(data []byte) string {
		return strconv.FormatBool(expected(data))
	})
}

func sample(rng *rand.Rand, i int) []byte {
	if i%2 == 0 {
		data := make([]byte, rng.Intn(9))
		rng.Read(data)
		return data
	}
	data := utf8.AppendRune(nil, rune(rng.Intn(0x110000)))
	data = data[:rng.Intn(len(data)+1)]
	if len(data) > 1 && i%4 == 1 {
		data[1+rng.Intn(len(data)-1)] = byte(rng.Uint32())
	}
	return data
}

func encode(data []byte) string {
	row := []byte("0")
	for _, b := range data {
		row = strconv.AppendUint(append(row, ' '), uint64(b), 10)
	}
	return string(row)
}

// A nonempty input is short only when its known bits can be completed to a
// scalar of the width encoded by the lead byte.
func expected(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	lead := data[0]
	var width, value, minimum int
	switch {
	case lead&0xe0 == 0xc0:
		width, value, minimum = 2, int(lead&0x1f), 0x80
	case lead&0xf0 == 0xe0:
		width, value, minimum = 3, int(lead&0xf), 0x800
	case lead&0xf8 == 0xf0:
		width, value, minimum = 4, int(lead&7), 0x10000
	default:
		return true
	}
	if len(data) >= width {
		return true
	}
	for _, b := range data[1:] {
		if b&0xc0 != 0x80 {
			return true
		}
		value = value*64 + int(b&0x3f)
	}
	scale := 1 << (6 * (width - len(data)))
	lo, hi := value*scale, (value+1)*scale-1
	canComplete := max(lo, minimum) <= min(hi, 0xd7ff) ||
		max(lo, minimum, 0xe000) <= min(hi, 0x10ffff)
	return !canComplete
}
