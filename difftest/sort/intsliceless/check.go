package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	xs   sort.IntSlice
	i, j int
}

func actual(c input) (out string) {
	defer func() {
		if p := recover(); p != nil {
			if e, ok := p.(runtime.Error); ok {
				out = e.Error()
			} else {
				out = fmt.Sprintf("unexpected panic %T: %v", p, p)
			}
		}
	}()
	return strconv.FormatBool(c.xs.Less(c.i, c.j))
}

func main() {
	edges := []input{
		{},
		{sort.IntSlice{}, -1, 0},
		{sort.IntSlice{math.MinInt, math.MaxInt}, 0, 1},
		{sort.IntSlice{math.MinInt, math.MaxInt}, 1, 0},
		{sort.IntSlice{1, 1}, 0, 1},
		{sort.IntSlice{0}, 0, 0},
		{sort.IntSlice{0}, 0, 1},
		{sort.IntSlice{0}, math.MinInt, math.MaxInt},
		{sort.IntSlice{0}, math.MaxInt, -1},
	}
	check.CompareSamples("godiff-intsliceless", edges, func(r *rand.Rand, _ int) input {
		xs := make(sort.IntSlice, r.Intn(17))
		for k := range xs {
			xs[k] = int(r.Uint64())
		}
		return input{xs, r.Intn(len(xs)+2) - 1, r.Intn(len(xs)+2) - 1}
	}, func(c input) string {
		var s strings.Builder
		fmt.Fprintf(&s, "%d %d", c.i, c.j)
		for _, value := range c.xs {
			fmt.Fprintf(&s, " %d", value)
		}
		return s.String()
	}, actual, func(c input) string {
		for _, index := range []int{c.i, c.j} {
			if index < 0 {
				return fmt.Sprintf("runtime error: index out of range [%d]", index)
			}
			if index >= len(c.xs) {
				return fmt.Sprintf("runtime error: index out of range [%d] with length %d", index, len(c.xs))
			}
		}
		return strconv.FormatBool(c.xs[c.i] < c.xs[c.j])
	})
}
