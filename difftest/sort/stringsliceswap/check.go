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
	xs   sort.StringSlice
	i, j int
}

func observe(status string, xs sort.StringSlice) string {
	var out strings.Builder
	out.WriteString(status + " [")
	for i, v := range xs {
		if i > 0 {
			out.WriteByte(' ')
		}
		out.WriteByte('[')
		for j := range len(v) {
			if j > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(&out, v[j])
		}
		out.WriteByte(']')
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
		{sort.StringSlice{}, -1, 0},
		{sort.StringSlice{"\x00\xff", "\xff"}, 0, 1},
		{sort.StringSlice{"\x00\xff", "\xff"}, 1, 0},
		{sort.StringSlice{"same", "same"}, 0, 1},
		{sort.StringSlice{""}, 0, 0},
		{sort.StringSlice{"a", "ab"}, 0, 1},
		{sort.StringSlice{"a", "ab"}, 1, 0},
		{sort.StringSlice{""}, 0, 1},
		{sort.StringSlice{""}, math.MinInt, math.MaxInt},
		{sort.StringSlice{""}, math.MaxInt, -1},
	}
	check.CompareSamples("godiff-stringsliceswap", edges, func(r *rand.Rand, _ int) input {
		xs := make(sort.StringSlice, r.Intn(17))
		for k := range xs {
			raw := make([]byte, r.Intn(9))
			r.Read(raw)
			xs[k] = string(raw)
		}
		return input{xs, r.Intn(len(xs)+2) - 1, r.Intn(len(xs)+2) - 1}
	}, func(c input) string {
		var s strings.Builder
		fmt.Fprintf(&s, "%d %d", c.i, c.j)
		for _, value := range c.xs {
			s.WriteString(" -1")
			for i := range len(value) {
				fmt.Fprintf(&s, " %d", value[i])
			}
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
		want := make(sort.StringSlice, len(c.xs))
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
