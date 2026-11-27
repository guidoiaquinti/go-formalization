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
	key int
}

func main() {
	checkFloat64()
	edges := []input{{nil, 0}, {[]int{}, 1}, {[]int{0}, 0}, {[]int{0}, -1}, {[]int{0}, 1}, {[]int{1, 1, 1}, 1}, {[]int{1, 3}, 2}, {[]int{math.MinInt, math.MaxInt}, math.MaxInt}}
	check.CompareSamples("godiff-slices-binary-search", edges, func(r *rand.Rand, i int) input {
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
		fmt.Fprintf(&row, "0 %d", c.key)
		for _, v := range c.xs {
			fmt.Fprintf(&row, " %d", v)
		}
		return row.String()
	}, func(c input) string {
		i, found := slices.BinarySearch(c.xs, c.key)
		return fmt.Sprintf("%d %t", i, found)
	}, func(c input) string {
		for i, v := range c.xs {
			if v >= c.key {
				return fmt.Sprintf("%d %t", i, v == c.key)
			}
		}
		return fmt.Sprintf("%d false", len(c.xs))
	})
}

type floatInput struct {
	xs  []float64
	key float64
}

func checkFloat64() {
	nan := math.Float64frombits(0x7ff0000000000001)
	edges := []floatInput{{nil, nan}, {[]float64{}, 1}, {[]float64{nan}, nan}, {[]float64{nan}, 0}, {[]float64{nan, 0, 1}, nan}, {[]float64{nan, 0, math.Copysign(0, -1), 1}, 0}, {[]float64{nan, math.Inf(-1), math.Inf(1)}, math.Inf(-1)}, {[]float64{nan, math.Inf(-1), math.Inf(1)}, math.Inf(1)}, {[]float64{1, 3}, 2}, {[]float64{1, 1, 1}, 1}}

	check.CompareSamples("godiff-slices-binary-search", edges, func(r *rand.Rand, i int) floatInput {
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
		return floatInput{xs, key}
	}, func(c floatInput) string {
		var row strings.Builder
		fmt.Fprintf(&row, "1 %d", math.Float64bits(c.key))
		for _, v := range c.xs {
			fmt.Fprintf(&row, " %d", math.Float64bits(v))
		}
		return row.String()
	}, func(c floatInput) string {
		i, found := slices.BinarySearch(c.xs, c.key)
		return fmt.Sprintf("%d %t", i, found)
	}, func(c floatInput) string {
		for i, v := range c.xs {
			if math.IsNaN(c.key) || !math.IsNaN(v) && v >= c.key {
				return fmt.Sprintf("%d %t", i, v == c.key || math.IsNaN(v) && math.IsNaN(c.key))
			}
		}
		return fmt.Sprintf("%d false", len(c.xs))
	})
}
