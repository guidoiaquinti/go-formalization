// Package bytefind exercises forward and backward byte searches.
package bytefind

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	data []byte
	key  byte
}

func Compare(binary string, actual func([]byte, byte) (int, int)) {
	var edges []input
	for a := range 256 {
		for c := range 256 {
			edges = append(edges, input{[]byte{byte(a)}, byte(c)})
		}
	}
	for _, length := range []int{0, 1, 2, 7, 8, 15, 16, 31, 32, 63, 64, 127, 128, 129, 256} {
		for offset := range 8 {
			backing := make([]byte, offset+length)
			data := backing[offset:]
			edges = append(edges, input{data, 1}, input{data, 0})
			if length > 0 {
				ends := make([]byte, offset+length)
				pattern := ends[offset:]
				pattern[0], pattern[len(pattern)-1] = 1, 1
				edges = append(edges, input{pattern, 1})
			}
		}
	}
	edges = append(edges, input{nil, 0})
	check.CompareSamples(binary, edges, sample, func(c input) string {
		out := strconv.AppendUint(nil, uint64(c.key), 10)
		for _, b := range c.data {
			out = strconv.AppendUint(append(out, ' '), uint64(b), 10)
		}
		return string(out)
	}, func(c input) string {
		first, last := actual(c.data, c.key)
		return fmt.Sprintf("%d %d", first, last)
	}, expected)
}

func expected(c input) string {
	var positions []int
	for i, b := range c.data {
		if b == c.key {
			positions = append(positions, i)
		}
	}
	if len(positions) == 0 {
		return "-1 -1"
	}
	return fmt.Sprintf("%d %d", positions[0], positions[len(positions)-1])
}

func sample(rng *rand.Rand, i int) input {
	offset := rng.Intn(8)
	data := make([]byte, offset+rng.Intn(129))[offset:]
	rng.Read(data)
	key := byte(rng.Uint32())
	if len(data) > 0 && i%3 != 0 {
		data[rng.Intn(len(data))] = key
		if i%3 == 2 {
			data[rng.Intn(len(data))] = key
		}
	}
	return input{data, key}
}
