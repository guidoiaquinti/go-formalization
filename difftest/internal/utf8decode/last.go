package utf8decode

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func CompareLast(executable string, includeNil bool, decode func([]byte) (rune, int)) {
	edges := [][]byte{{}}
	if includeNil {
		edges = append(edges, nil)
	}
	for a := range 256 {
		edges = append(edges, []byte{byte(a)})
		for b := range 256 {
			edges = append(edges, []byte{byte(a), byte(b)})
		}
	}
	for _, word := range []string{
		"prefix¢", "prefix€", "prefix🐈", "prefix\xef\xbf\xbd",
		"\x80\x80\x80\x80", "A" + strings.Repeat("\x80", 64),
		"🐈\x80", "€\x80", "A\xe2\x82", "A\xf4\x8f\xbf\xbf",
		"A\xf4\x90\x80\x80", "\x80\x80A", "A\xed\xa0\x80",
	} {
		edges = append(edges, []byte(word))
	}
	check.CompareSamples(executable, edges, sampleLast, encodeLast,
		func(data []byte) string {
			r, n := decode(data)
			return fmt.Sprintf("%d %d", r, n)
		},
		func(data []byte) string {
			r, n := rune(0xfffd), 0
			for len(data) != 0 {
				r, n = Expected(data)
				data = data[n:]
			}
			return fmt.Sprintf("%d %d", r, n)
		})
}

func sampleLast(rng *rand.Rand, i int) []byte {
	raw := make([]byte, rng.Intn(17))
	rng.Read(raw)
	switch i % 4 {
	case 1:
		raw = utf8.AppendRune(raw, rune(rng.Intn(0x110000)))
	case 2:
		raw = append(raw, []byte(strings.Repeat("\x80", 1+rng.Intn(12)))...)
	case 3:
		last := utf8.AppendRune(nil, rune(rng.Intn(0x110000)))
		raw = append(raw, last[:rng.Intn(len(last))]...)
	}
	return raw
}

func encodeLast(s []byte) string {
	row := []byte("0")
	for i := range len(s) {
		row = strconv.AppendUint(append(row, ' '), uint64(s[i]), 10)
	}
	return string(row)
}
