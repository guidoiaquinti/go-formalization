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
	xs   sort.Float64Slice
	i, j int
}

func observe(status string, xs sort.Float64Slice) string {
	var out strings.Builder
	out.WriteString(status + " [")
	for i, v := range xs {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(&out, math.Float64bits(v))
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
		{sort.Float64Slice{}, -1, 0},
		{sort.Float64Slice{math.Float64frombits(0x7ff0000000000001), math.Float64frombits(0xfff8000000000002)}, 0, 1},
		{sort.Float64Slice{math.Float64frombits(0x7ff0000000000001), math.Float64frombits(0xfff8000000000002)}, 1, 0},
		{sort.Float64Slice{0, math.Copysign(0, -1)}, 0, 1},
		{sort.Float64Slice{0}, 0, 0},
		{sort.Float64Slice{0}, 0, 1},
		{sort.Float64Slice{0}, math.MinInt, math.MaxInt},
		{sort.Float64Slice{0}, math.MaxInt, -1},
	}
	check.CompareSamples("godiff-float64sliceswap", edges, func(r *rand.Rand, _ int) input {
		xs := make(sort.Float64Slice, r.Intn(17))
		for k := range xs {
			xs[k] = math.Float64frombits(r.Uint64())
		}
		return input{xs, r.Intn(len(xs)+2) - 1, r.Intn(len(xs)+2) - 1}
	}, func(c input) string {
		var s strings.Builder
		fmt.Fprintf(&s, "%d %d", c.i, c.j)
		for _, value := range c.xs {
			fmt.Fprintf(&s, " %d", math.Float64bits(value))
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
		want := make(sort.Float64Slice, len(c.xs))
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
