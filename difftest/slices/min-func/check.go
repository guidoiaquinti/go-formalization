package main

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type entry struct {
	key   int
	index int
}

func compareKey(a, b entry) int {
	if a.key < b.key {
		return -7
	}
	if a.key == b.key {
		return 0
	}
	return 13
}

func actual(xs []int) (result string) {
	defer func() {
		if p := recover(); p != nil {
			if message, ok := p.(string); ok {
				result = message
			} else {
				result = fmt.Sprintf("unexpected panic %T: %v", p, p)
			}
		}
	}()
	entries := make([]entry, len(xs))
	for i, key := range xs {
		entries[i] = entry{key, i}
	}
	resultEntry := slices.MinFunc(entries, compareKey)
	return fmt.Sprintf("%d %d", resultEntry.key, resultEntry.index)
}

func main() {
	edges := [][]int{nil, {}, {0}, {1, 1}, {math.MinInt, math.MaxInt}, {math.MaxInt, math.MinInt}, {1, 0, 2, 3}, {1, 2, 3, 0}, {3, -2, 0, -2, 1}}
	check.CompareSamples("godiff-slices-min-func", edges, func(r *rand.Rand, i int) []int {
		xs := make([]int, r.Intn(33))
		for j := range xs {
			if i%2 == 0 {
				xs[j] = int(r.Uint64())
			} else {
				xs[j] = r.Intn(9) - 4
			}
		}
		return xs
	}, func(xs []int) string {
		var row strings.Builder
		row.WriteByte('0')
		for _, x := range xs {
			fmt.Fprintf(&row, " %d", x)
		}
		return row.String()
	}, actual, func(xs []int) string {
		if len(xs) == 0 {
			return "slices.MinFunc: empty list"
		}
		key := slices.Min(xs)
		return fmt.Sprintf("%d %d", key, slices.Index(xs, key))
	})
}
