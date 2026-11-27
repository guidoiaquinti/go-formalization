package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"slices"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	values []int
	count  int
}

func row(header string, values []int) string {
	var out strings.Builder
	out.WriteString(header)
	for _, v := range values {
		fmt.Fprintf(&out, " %d", v)
	}
	return out.String()
}

func actual(c input) (result string) {
	defer func() {
		if p := recover(); p != nil {
			if message, ok := p.(string); ok {
				result = message
			} else {
				result = fmt.Sprintf("unexpected panic %T: %v", p, p)
			}
		}
	}()
	got := slices.Repeat(c.values, c.count)
	if got == nil || cap(got) != len(got) {
		log.Fatalf("Repeat: invalid result header for count %d", c.count)
	}
	return row("ok", got)
}

func main() {
	edges := []input{{nil, 0}, {nil, math.MaxInt}, {[]int{}, 2}, {nil, -1}, {[]int{1}, 0}, {[]int{1, 2}, 1}, {[]int{1, 2, 3}, 5}, {[]int{1, 2}, math.MaxInt}, {[]int{1, 2, 3}, math.MaxInt}, {[]int{1}, math.MinInt}}
	check.CompareSamples("godiff-slices-repeat", edges, func(r *rand.Rand, _ int) input {
		xs := make([]int, r.Intn(13))
		for i := range xs {
			xs[i] = int(r.Uint64())
		}
		return input{xs, r.Intn(17) - 2}
	}, func(c input) string {
		return row(fmt.Sprintf("%d %d", math.MaxInt, c.count), c.values)
	}, actual, func(c input) string {
		if c.count < 0 {
			return "cannot be negative"
		}
		if len(c.values) > 0 && c.count > math.MaxInt/len(c.values) {
			return "the result of (len(x) * count) overflows"
		}
		expected := make([]int, len(c.values)*c.count)
		for i := range expected {
			expected[i] = c.values[i%len(c.values)]
		}
		return row("ok", expected)
	})
}
