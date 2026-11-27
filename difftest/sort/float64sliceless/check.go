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
	xs   sort.Float64Slice
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
		{sort.Float64Slice{}, -1, 0},
		{sort.Float64Slice{math.Inf(-1), math.Inf(1)}, 0, 1},
		{sort.Float64Slice{math.Inf(-1), math.Inf(1)}, 1, 0},
		{sort.Float64Slice{1, 1}, 0, 1},
		{sort.Float64Slice{0}, 0, 0},
		{sort.Float64Slice{0}, 0, 1},
		{sort.Float64Slice{0}, math.MinInt, math.MaxInt},
		{sort.Float64Slice{0}, math.MaxInt, -1},
	}
	special := []float64{0, math.Copysign(0, -1), math.Inf(-1), math.Inf(1), math.SmallestNonzeroFloat64, -math.SmallestNonzeroFloat64, math.MaxFloat64, math.NaN(), math.Float64frombits(0x7ff0000000000001), math.Float64frombits(0xfff8000000000001)}
	for i := range special {
		for j := range special {
			edges = append(edges, input{special, i, j})
		}
	}
	check.CompareSamples("godiff-float64sliceless", edges, func(r *rand.Rand, _ int) input {
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
				return fmt.Sprintf("runtime error: index out of range [%d]", index)
			}
			if index >= len(c.xs) {
				return fmt.Sprintf("runtime error: index out of range [%d] with length %d", index, len(c.xs))
			}
		}
		a, b := c.xs[c.i], c.xs[c.j]
		if math.IsNaN(b) {
			return "false"
		}
		return strconv.FormatBool(math.IsNaN(a) || a < b)
	})
}
