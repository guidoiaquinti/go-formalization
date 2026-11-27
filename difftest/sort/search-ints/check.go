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

type input struct {
	xs  []int
	key int
}

func main() {
	edges := []input{{nil, 0}, {[]int{}, 1}, {[]int{0}, 0}, {[]int{0}, -1}, {[]int{0}, 1}, {[]int{1, 1, 1}, 1}, {[]int{1, 3}, 2}, {[]int{math.MinInt, math.MaxInt}, math.MaxInt}}
	check.CompareSamples("godiff-search-ints", edges, func(r *rand.Rand, i int) input {
		xs := make([]int, r.Intn(65))
		for j := range xs {
			if i%2 == 0 {
				xs[j] = r.Intn(17) - 8
			} else {
				xs[j] = int(r.Uint64())
			}
		}
		slices.Sort(xs)
		key := int(r.Uint64())
		if i%2 == 0 {
			key = r.Intn(21) - 10
		} else if i%4 == 1 && len(xs) != 0 {
			key = xs[r.Intn(len(xs))]
		}
		return input{xs, key}
	}, func(c input) string {
		var row strings.Builder
		fmt.Fprint(&row, c.key)
		for _, v := range c.xs {
			fmt.Fprintf(&row, " %d", v)
		}
		return row.String()
	}, func(c input) string {
		return strconv.Itoa(sort.SearchInts(c.xs, c.key))
	}, func(c input) string {
		for i, v := range c.xs {
			if v >= c.key {
				return strconv.Itoa(i)
			}
		}
		return strconv.Itoa(len(c.xs))
	})
}
