// Package byteorder exercises lexicographic order on pairs of byte sequences.
package byteorder

import (
	"bytes"
	"math/rand"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type pair struct {
	left, right []byte
}

func Compare(binary string, actual func(a, b []byte) string, expected func(order int) string) {
	values := [][]byte{nil, {}, {0}, {0, 1}, {0, 0}, {1, 0}, {127}, {128}, {255}, {0xc0, 0x80}}
	var edges []pair
	for _, left := range values {
		for _, right := range values {
			edges = append(edges, pair{left, right})
		}
	}
	for a := range 256 {
		for b := range 256 {
			edges = append(edges, pair{[]byte{byte(a)}, []byte{byte(b)}})
		}
	}
	check.CompareSamples(binary, edges, sample, encode, func(c pair) string {
		return actual(c.left, c.right)
	}, func(c pair) string {
		return expected(order(c))
	})
}

func sample(rng *rand.Rand, i int) pair {
	a := make([]byte, rng.Intn(129))
	rng.Read(a)
	b := make([]byte, rng.Intn(129))
	rng.Read(b)
	switch i % 5 {
	case 0:
		b = a
	case 1:
		b = a[:rng.Intn(len(a)+1)]
	case 2:
		b = bytes.Clone(a)
	case 3:
		b = bytes.Clone(a)
		if len(b) != 0 {
			b[rng.Intn(len(b))] ^= 1
		}
	}
	return pair{a, b}
}

func order(c pair) int {
	// Fold from the end so the first differing byte has the final word.
	order := 0
	if len(c.left) < len(c.right) {
		order = -1
	} else if len(c.left) > len(c.right) {
		order = 1
	}
	for i := min(len(c.left), len(c.right)) - 1; i >= 0; i-- {
		if c.left[i] < c.right[i] {
			order = -1
		} else if c.left[i] > c.right[i] {
			order = 1
		}
	}
	return order
}

func encode(c pair) string {
	var row []byte
	for _, data := range [][]byte{c.left, c.right} {
		length := len(data)
		if data == nil {
			length = -1
		}
		if len(row) != 0 {
			row = append(row, ' ')
		}
		row = strconv.AppendInt(row, int64(length), 10)
	}
	for _, data := range [][]byte{c.left, c.right} {
		for _, word := range data {
			row = strconv.AppendUint(append(row, ' '), uint64(word), 10)
		}
	}
	return string(row)
}
