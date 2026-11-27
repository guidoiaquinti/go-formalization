package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"slices"
	"sort"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	xs   sort.IntSlice
	i, j int
}

func observe(status string, xs sort.IntSlice) string {
	var out strings.Builder
	out.WriteString(status + " [")
	for i, v := range xs {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(&out, v)
	}
	out.WriteByte(']')
	return out.String()
}

func actual(c input) (out string) {
	xs := slices.Clone(c.xs)
	status := "ok"
	defer func() {
		if p := recover(); p != nil {
			if e, ok := p.(runtime.Error); ok {
				status = e.Error()
			} else {
				status = fmt.Sprintf("unexpected panic %T: %v", p, p)
			}
		}
		out = observe(status, xs)
	}()
	xs.Swap(c.i, c.j)
	return
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
	check.CompareSamples("godiff-intsliceswap", edges, func(r *rand.Rand, _ int) input {
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
				return observe(fmt.Sprintf("runtime error: index out of range [%d]", index), c.xs)
			}
			if index >= len(c.xs) {
				return observe(fmt.Sprintf("runtime error: index out of range [%d] with length %d", index, len(c.xs)), c.xs)
			}
		}
		want := make(sort.IntSlice, len(c.xs))
		for k := range want {
			source := k
			if k == c.i {
				source = c.j
			} else if k == c.j {
				source = c.i
			}
			want[k] = c.xs[source]
		}
		return observe("ok", want)
	})
}
