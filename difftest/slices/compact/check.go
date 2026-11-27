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
	edges := []input{{}, {[]int{}, 0, 0}, {[]int{}, 1, 1}, {[]int{1}, 0, 0}, {[]int{1, 2}, 7, 3}, {[]int{1, 2, 3}, 2, 1}, {[]int{2, 2, 2}, 1, 2}, {[]int{1, 1, 2, 1, 1, 0, 0, 3}, 4, 2}, {[]int{0, 0}, 2, 1}}
	check.CompareSamples("godiff-slices-compact", edges, func(r *rand.Rand, sample int) input {
		values := make([]int, r.Intn(33))
		for i := range values {
			if sample%2 == 0 {
				values[i] = r.Intn(7) - 3
			} else if i != 0 && r.Intn(3) == 0 {
				values[i] = values[i-1]
			} else {
				values[i] = int(r.Uint64())
			}
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
		result := slices.Compact(view)
		if (result == nil) != (view == nil) || cap(result) != capacity || len(result) > len(view) {
			log.Fatal("Compact changed nilness or capacity, or expanded the length")
		}
		if capacity != 0 && &result[:capacity][0] != &view[:capacity][0] {
			log.Fatal("Compact moved the slice to a different backing array")
		}
		if backing != nil {
			for i, value := range backing {
				if (i < c.offset || i >= c.offset+len(view)) && value != i+1 {
					log.Fatalf("Compact wrote outside the view at %d", i)
				}
			}
		}
		return row(fmt.Sprint(len(result)), view)
	}, func(c input) string {
		expected := make([]int, len(c.values))
		kept := 0
		for i, value := range c.values {
			if i == 0 || c.values[i-1] != value {
				expected[kept] = value
				kept++
			}
		}
		return row(fmt.Sprint(kept), expected)
	})
}
