package main

import (
	"fmt"
	"math/rand"
	"reflect"
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
	check.CompareSamples("godiff-slices-equal", edges, func(r *rand.Rand, i int) pair {
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
		return strconv.FormatBool(slices.Equal(p.left, p.right))
	}, func(p pair) string {
		return strconv.FormatBool((len(p.left) == 0 && len(p.right) == 0) || reflect.DeepEqual(p.left, p.right))
	})
}
