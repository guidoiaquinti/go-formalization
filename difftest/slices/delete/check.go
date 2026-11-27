package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"slices"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type namedSlice []int
type input struct {
	values        []int
	i, j          int
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

func callDelete(view namedSlice, i, j int) (result namedSlice, panicked bool) {
	defer func() {
		if p := recover(); p != nil {
			if _, ok := p.(runtime.Error); !ok {
				log.Fatalf("unexpected panic %T: %v", p, p)
			}
			panicked = true
		}
	}()
	return slices.Delete(view, i, j), false
}

func main() {
	edges := []input{
		{},
		{[]int{}, 0, 0, 1, 1},
		{[]int{1}, 0, 0, 0, 0},
		{[]int{1}, 0, 1, 2, 1},
		{[]int{1, 2, 3, 4}, 1, 3, 7, 3},
		{[]int{1, 2, 3}, 2, 3, 2, 1},
		{[]int{1}, -1, 0, 2, 1},
		{[]int{1}, 1, 0, 2, 1},
		{[]int{1}, 0, 2, 2, 1},
	}
	check.CompareSamples("godiff-slices-delete", edges, func(r *rand.Rand, sample int) input {
		values := make([]int, r.Intn(33))
		for k := range values {
			values[k] = int(r.Uint64())
		}
		i := r.Intn(len(values) + 1)
		j := i + r.Intn(len(values)-i+1)
		if sample%2 == 0 {
			i = r.Intn(len(values)+9) - 4
			j = r.Intn(len(values)+9) - 4
		}
		return input{values, i, j, r.Intn(9), r.Intn(5)}
	}, func(c input) string {
		return row(fmt.Sprintf("%d %d", c.i, c.j), c.values)
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
		result, panicked := callDelete(view, c.i, c.j)
		header := "panic"
		if !panicked {
			if (result == nil) != (view == nil) || cap(result) != capacity || len(result) > len(view) {
				log.Fatal("Delete changed nilness or capacity, or expanded the length")
			}
			if capacity != 0 && &result[:capacity][0] != &view[:capacity][0] {
				log.Fatal("Delete moved the slice to a different backing array")
			}
			header = fmt.Sprint(len(result))
		}
		if backing != nil {
			for i, value := range backing {
				if (i < c.offset || i >= c.offset+len(view)) && value != i+1 {
					log.Fatalf("Delete wrote outside the view at %d", i)
				}
			}
		}
		return row(header, view)
	}, func(c input) string {
		if c.i < 0 || c.i > c.j || c.j > len(c.values) {
			return row("panic", c.values)
		}
		expected := make([]int, len(c.values))
		for k := 0; k < c.i; k++ {
			expected[k] = c.values[k]
		}
		for k := c.j; k < len(c.values); k++ {
			expected[k-(c.j-c.i)] = c.values[k]
		}
		return row(fmt.Sprint(len(c.values)-(c.j-c.i)), expected)
	})
}
