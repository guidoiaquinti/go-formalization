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
	xs  []float64
	key float64
}

func main() {
	nan := math.Float64frombits(0x7ff0000000000001)
	edges := []input{{nil, nan}, {[]float64{}, 1}, {[]float64{nan}, nan}, {[]float64{nan}, 0}, {[]float64{nan, 0, 1}, nan}, {[]float64{nan, 0, math.Copysign(0, -1), 1}, 0}, {[]float64{nan, math.Inf(-1), math.Inf(1)}, math.Inf(-1)}, {[]float64{nan, math.Inf(-1), math.Inf(1)}, math.Inf(1)}, {[]float64{1, 3}, 2}, {[]float64{1, 1, 1}, 1}}

	check.CompareSamples("godiff-float64slice-search", edges, func(r *rand.Rand, i int) input {
		xs := make([]float64, r.Intn(65))
		for j := range xs {
			if i%2 == 0 {
				xs[j] = []float64{nan, -1, 0, 1, math.NaN()}[r.Intn(5)]
			} else {
				xs[j] = math.Float64frombits(r.Uint64())
			}
		}
		slices.Sort(xs)
		key := math.Float64frombits(r.Uint64())
		if i%2 == 0 {
			key = []float64{nan, -2, -1, 0, 1, 2}[r.Intn(6)]
		} else if i%4 == 1 && len(xs) != 0 {
			key = xs[r.Intn(len(xs))]
		}
		return input{xs, key}
	}, func(c input) string {
		var row strings.Builder
		fmt.Fprint(&row, math.Float64bits(c.key))
		for _, v := range c.xs {
			fmt.Fprintf(&row, " %d", math.Float64bits(v))
		}
		return row.String()
	}, func(c input) string {
		return strconv.Itoa(sort.Float64Slice(c.xs).Search(c.key))
	}, func(c input) string {
		for i, v := range c.xs {
			if v >= c.key {
				return strconv.Itoa(i)
			}
		}
		return strconv.Itoa(len(c.xs))
	})
}
