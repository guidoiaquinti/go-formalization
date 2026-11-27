package utf8decode

import (
	"bytes"
	"math/rand"
	"strconv"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func CompareValid(executable string, includeNil bool, valid func([]byte) bool) {
	edges := [][]byte{{}}
	if includeNil {
		edges = append([][]byte{nil}, edges...)
	}
	for a := range 256 {
		edges = append(edges, []byte{byte(a)})
		for b := range 256 {
			edges = append(edges, []byte{byte(a), byte(b)})
		}
	}
	for n := range 81 {
		prefix := bytes.Repeat([]byte{'A'}, n)
		for _, suffix := range []string{"", "¢", "€", "🐈", "\xef\xbf\xbd", "\xc0\x80",
			"\xed\xa0\x80", "\xf4\x90\x80\x80", "\xe2\x82", "\x80"} {
			edges = append(edges, append(bytes.Clone(prefix), suffix...))
		}
	}
	check.CompareSamples(executable, edges, sampleValid, encodeValid, func(data []byte) string {
		return strconv.FormatBool(valid(data))
	}, func(data []byte) string {
		for len(data) != 0 {
			r, width := Expected(data)
			if r == 0xfffd && width == 1 {
				return "false"
			}
			data = data[width:]
		}
		return "true"
	})
}

func sampleValid(r *rand.Rand, i int) []byte {
	if i%4 == 0 {
		data := make([]byte, r.Intn(129))
		r.Read(data)
		return data
	}
	data := bytes.Repeat([]byte{byte(r.Intn(128))}, r.Intn(129))
	for j := r.Intn(9); j > 0; j-- {
		data = utf8.AppendRune(data, rune(r.Intn(0x110000)))
	}
	if i%4 == 2 && len(data) > 0 {
		data[r.Intn(len(data))] = byte(r.Uint32())
	}
	return data
}

func encodeValid(data []byte) string {
	row := strconv.AppendInt(nil, int64(strconv.IntSize/8), 10)
	for _, b := range data {
		row = strconv.AppendUint(append(row, ' '), uint64(b), 10)
	}
	return string(row)
}
