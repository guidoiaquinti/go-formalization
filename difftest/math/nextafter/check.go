package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type pair struct{ x, y uint64 }

const sign uint64 = 1 << 63
const infinity uint64 = 0x7ff0000000000000

func rank(word uint64) int64 {
	magnitude := int64(word % sign)
	if word >= sign {
		return -magnitude
	}
	return magnitude
}

func expected(p pair) string {
	a, b := rank(p.x), rank(p.y)
	if p.x%sign > infinity || p.y%sign > infinity {
		return strconv.FormatUint(0x7ff8000000000001, 10)
	}
	if a == b {
		return strconv.FormatUint(p.x, 10)
	}
	if a < b {
		a++
	} else {
		a--
	}
	word := uint64(a)
	if a < 0 || a == 0 && p.x >= sign {
		word = sign + uint64(-a)
	}
	return strconv.FormatUint(word, 10)
}

func sample(rng *rand.Rand, i int) pair {
	x, y := rng.Uint64(), rng.Uint64()
	switch i % 4 {
	case 0:
		y = x
	case 1:
		y = x + 1
	case 2:
		y = x - 1
	}
	return pair{x, y}
}

func main() {
	bounds := binary64.Boundaries()
	var edges []pair
	for _, x := range bounds {
		for _, y := range bounds {
			edges = append(edges, pair{x, y})
		}
	}
	check.CompareSamples("godiff-nextafter", edges, sample,
		func(p pair) string { return fmt.Sprintf("%d %d", p.x, p.y) },
		func(p pair) string {
			return strconv.FormatUint(math.Float64bits(math.Nextafter(binary64.Store(p.x), binary64.Store(p.y))), 10)
		}, expected)
}
