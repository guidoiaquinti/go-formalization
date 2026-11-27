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
	nan := math.Float64frombits(0x7ff0000000000001)
	edges := [][]float64{nil, {}, {nan}, {nan, math.NaN(), math.Inf(-1), 0, math.Inf(1)}, {0, nan}, {0, math.Copysign(0, -1)}, {math.Inf(1), math.Inf(-1)}, {-math.SmallestNonzeroFloat64, 0, math.SmallestNonzeroFloat64}, {1, 0, 2, 3}}

	check.CompareSamples("godiff-float64s-are-sorted", edges, func(r *rand.Rand, i int) []float64 {
		xs := make([]float64, r.Intn(33))
		for j := range xs {
			if i%3 == 0 {
				xs[j] = []float64{nan, -1, 0, 1, math.NaN()}[r.Intn(5)]
			} else {
				xs[j] = math.Float64frombits(r.Uint64())
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
	}, func(xs []float64) string {
		var row strings.Builder
		row.WriteByte('0')
		for _, x := range xs {
			fmt.Fprintf(&row, " %d", math.Float64bits(x))
		}
		return row.String()
	}, func(xs []float64) string {
		return strconv.FormatBool(sort.Float64sAreSorted(xs))
	}, func(xs []float64) string {
		for i, x := range xs {
			for _, y := range xs[i+1:] {
				if !math.IsNaN(x) && (math.IsNaN(y) || x > y) {
					return "false"
				}
			}
		}
		return "true"
	})
}
