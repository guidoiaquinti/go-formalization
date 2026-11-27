package main

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

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
	return strconv.Itoa(slices.Min(xs))
}

func main() {
	edges := [][]int{nil, {}, {0}, {1, 1}, {math.MinInt, math.MaxInt}, {math.MaxInt, math.MinInt}, {1, 0, 2, 3}, {1, 2, 3, 0}}
	check.CompareSamples("godiff-slices-min", edges, func(r *rand.Rand, i int) []int {
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
			return "slices.Min: empty list"
		}
		ordered := slices.Clone(xs)
		slices.Sort(ordered)
		return strconv.Itoa(ordered[0])
	})
}
