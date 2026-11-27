package main

import (
	"fmt"
	"log"
	"math/rand"
	"slices"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type namedSlice []int
type input struct {
	values        []int
	spare, offset int
}

func row(header string, values []int) string {
	var out strings.Builder
	out.WriteString(header)
	for _, v := range values {
		fmt.Fprintf(&out, " %d", v)
	}
	return out.String()
}

func main() {
	edges := []input{{}, {[]int{}, 0, 0}, {[]int{}, 1, 1}, {[]int{1}, 0, 0}, {[]int{1, 2}, 7, 3}, {[]int{1, 2, 3}, 2, 1}}
	check.CompareSamples("godiff-slices-reverse", edges, func(r *rand.Rand, _ int) input {
		values := make([]int, r.Intn(33))
		for i := range values {
			values[i] = int(r.Uint64())
		}
		return input{values, r.Intn(9), r.Intn(5)}
	}, func(c input) string {
		return row("0", c.values)
	}, func(c input) string {
		var view, backing namedSlice
		if c.values != nil {
			backing = make(namedSlice, c.offset+len(c.values)+c.spare)
			for i := range backing {
				backing[i] = i + 1
			}
			view = backing[c.offset : c.offset+len(c.values)]
			copy(view, c.values)
		}
		capacity := cap(view)
		slices.Reverse(view)
		if (view == nil) != (c.values == nil) || len(view) != len(c.values) || cap(view) != capacity {
			log.Fatal("Reverse changed the slice header")
		}
		if backing != nil {
			for i, value := range backing {
				if (i < c.offset || i >= c.offset+len(view)) && value != i+1 {
					log.Fatalf("Reverse wrote outside the view at %d", i)
				}
			}
		}
		return row("ok", view)
	}, func(c input) string {
		expected := make([]int, len(c.values))
		for i := range expected {
			expected[i] = c.values[len(expected)-1-i]
		}
		return row("ok", expected)
	})
}
