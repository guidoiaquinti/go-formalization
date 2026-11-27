package main

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	edges := [][]int{nil, {}, {0}, {1, 1}, {math.MinInt, math.MaxInt}, {math.MaxInt, math.MinInt}, {1, 0, 2, 3}, {1, 2, 3, 0}, {1, 2, 2, 3}}
	check.CompareSamples("godiff-ints-are-sorted", edges, func(r *rand.Rand, i int) []int {
		xs := make([]int, r.Intn(33))
		for j := range xs {
			if i%3 == 0 {
				xs[j] = r.Intn(9) - 4
			} else {
				xs[j] = int(r.Uint64())
			}
		}
		if i%3 != 2 {
			slices.Sort(xs)
			if i%3 == 1 && len(xs) > 1 {
				j := r.Intn(len(xs) - 1)
				xs[j], xs[j+1] = xs[j+1], xs[j]
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
	}, func(xs []int) string {
		return strconv.FormatBool(sort.IntsAreSorted(xs))
	}, func(xs []int) string {
		for i, x := range xs {
			for _, y := range xs[i+1:] {
				if x > y {
					return "false"
				}
			}
		}
		return "true"
	})
}
