package main

import (
	"bytes"
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
	xs   sort.StringSlice
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
	check.CompareSamples("godiff-stringsliceless", edges, func(r *rand.Rand, _ int) input {
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
				return fmt.Sprintf("runtime error: index out of range [%d]", index)
			}
			if index >= len(c.xs) {
				return fmt.Sprintf("runtime error: index out of range [%d] with length %d", index, len(c.xs))
			}
		}
		return strconv.FormatBool(bytes.Compare([]byte(c.xs[c.i]), []byte(c.xs[c.j])) < 0)
	})
}
