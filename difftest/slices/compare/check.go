package main

import (
	"cmp"
	"fmt"
	"math/rand"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type pair struct{ left, right []int }

func main() {
	edges := []pair{{}, {nil, []int{}}, {[]int{}, nil}, {nil, []int{0}},
		{[]int{1, 2}, []int{1, 2}}, {[]int{1, 2}, []int{0, 2}}, {[]int{1, 2}, []int{1, 3}},
		{[]int{1, 2, 3}, []int{1, 2}}}
	check.CompareSamples("godiff-slices-compare", edges, func(r *rand.Rand, i int) pair {
		left := make([]int, r.Intn(65))
		for j := range left {
			left[j] = int(r.Uint64())
		}
		right := slices.Clone(left)
		switch i % 4 {
		case 0:
		case 1:
			if len(right) != 0 {
				right[r.Intn(len(right))] ^= 1
			}
		case 2:
			right = append(right, int(r.Uint64()))
		case 3:
			for j := range right {
				right[j] = int(r.Uint64())
			}
		}
		return pair{left, right}
	}, func(p pair) string {
		var row strings.Builder
		fmt.Fprint(&row, len(p.left))
		for _, xs := range [][]int{p.left, p.right} {
			for _, x := range xs {
				fmt.Fprintf(&row, " %d", x)
			}
		}
		return row.String()
	}, func(p pair) string {
		return strconv.Itoa(slices.Compare(p.left, p.right))
	}, func(p pair) string {
		return strconv.Itoa(expected(p))
	})
}

func expected(p pair) int {
	// Start with the prefix tie-breaker, then let each earlier mismatch take precedence.
	order := cmp.Compare(len(p.left), len(p.right))
	for i := min(len(p.left), len(p.right)) - 1; i >= 0; i-- {
		if c := cmp.Compare(p.left[i], p.right[i]); c != 0 {
			order = c
		}
	}
	return order
}
