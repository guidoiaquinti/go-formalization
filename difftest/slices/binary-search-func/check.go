package main

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	xs  []int
	key searchKey
}

type searchKey struct {
	value      int
	descending bool
}

func compareElement(element int, key searchKey) int {
	if element == key.value {
		return 0
	}
	if (element < key.value) != key.descending {
		return -7
	}
	return 13
}

func main() {
	edges := []input{
		{nil, searchKey{0, false}},
		{[]int{}, searchKey{1, true}},
		{[]int{0}, searchKey{0, false}},
		{[]int{0}, searchKey{-1, false}},
		{[]int{0}, searchKey{1, false}},
		{[]int{3, 2, 2, 1}, searchKey{2, true}},
		{[]int{3, 1}, searchKey{2, true}},
		{[]int{math.MinInt, math.MaxInt}, searchKey{math.MaxInt, false}},
		{[]int{math.MaxInt, math.MinInt}, searchKey{math.MinInt, true}},
	}
	check.CompareSamples("godiff-slices-binary-search-func", edges, func(r *rand.Rand, i int) input {
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
		descending := i%4 >= 2
		if descending {
			slices.Reverse(xs)
		}
		return input{xs, searchKey{key, descending}}
	}, func(c input) string {
		var row strings.Builder
		if c.key.descending {
			row.WriteByte('1')
		} else {
			row.WriteByte('0')
		}
		fmt.Fprintf(&row, " %d", c.key.value)
		for _, v := range c.xs {
			fmt.Fprintf(&row, " %d", v)
		}
		return row.String()
	}, func(c input) string {
		i, found := slices.BinarySearchFunc(c.xs, c.key, compareElement)
		return fmt.Sprintf("%d %t", i, found)
	}, func(c input) string {
		for i, v := range c.xs {
			if (!c.key.descending && v >= c.key.value) || (c.key.descending && v <= c.key.value) {
				return fmt.Sprintf("%d %t", i, v == c.key.value)
			}
		}
		return fmt.Sprintf("%d false", len(c.xs))
	})
}
