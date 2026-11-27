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
	panicOn13     bool
	spare, offset int
}

func row(header string, values, trace []int) string {
	var out strings.Builder
	out.WriteString(header)
	for _, v := range values {
		fmt.Fprintf(&out, " %d", v)
	}
	if trace != nil {
		out.WriteString(" |")
		for _, v := range trace {
			fmt.Fprintf(&out, " %d", v)
		}
	}
	return out.String()
}

func callDeleteFunc(view namedSlice, panicOn13 bool) (result namedSlice, panicked bool, trace []int) {
	trace = make([]int, 0, len(view))
	defer func() {
		if p := recover(); p != nil {
			if p != "predicate panic" {
				log.Fatalf("unexpected panic %T: %v", p, p)
			}
			panicked = true
		}
	}()
	result = slices.DeleteFunc(view, func(value int) bool {
		trace = append(trace, value)
		if panicOn13 && value == 13 {
			panic("predicate panic")
		}
		return value%2 == 0
	})
	return result, false, trace
}

func main() {
	edges := []input{
		{}, {[]int{}, true, 1, 1}, {[]int{0}, false, 2, 1},
		{[]int{1, -3}, false, 1, 2}, {[]int{4, -2, 0}, false, 3, 1},
		{[]int{2, 3, 4, 5}, false, 4, 2}, {[]int{13, 1}, true, 1, 1},
		{[]int{1, 13, 2}, true, 1, 1}, {[]int{1, 2, 3, 13, 6}, true, 3, 2},
	}
	check.CompareSamples("godiff-slices-delete-func", edges, func(r *rand.Rand, sample int) input {
		values := make([]int, r.Intn(33))
		for i := range values {
			if r.Intn(7) == 0 {
				values[i] = 13
			} else {
				values[i] = int(r.Uint64())
			}
		}
		return input{values, sample%2 == 0, r.Intn(9), r.Intn(5)}
	}, func(c input) string {
		if c.panicOn13 {
			return row("1", c.values, nil)
		}
		return row("0", c.values, nil)
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
		result, panicked, trace := callDeleteFunc(view, c.panicOn13)
		header := "panic"
		if !panicked {
			if (result == nil) != (view == nil) || cap(result) != capacity || len(result) > len(view) {
				log.Fatal("DeleteFunc changed nilness or capacity, or expanded the length")
			}
			if capacity != 0 && &result[:capacity][0] != &view[:capacity][0] {
				log.Fatal("DeleteFunc moved the slice to a different backing array")
			}
			header = fmt.Sprint(len(result))
		}
		if backing != nil {
			for i, value := range backing {
				if (i < c.offset || i >= c.offset+len(view)) && value != i+1 {
					log.Fatalf("DeleteFunc wrote outside the view at %d", i)
				}
			}
		}
		return row(header, view, trace)
	}, func(c input) string {
		limit := len(c.values)
		if c.panicOn13 {
			if i := slices.Index(c.values, 13); i >= 0 {
				limit = i
			}
		}
		expected := make([]int, len(c.values))
		kept := 0
		for _, value := range c.values[:limit] {
			if value%2 != 0 {
				expected[kept] = value
				kept++
			}
		}
		if limit < len(c.values) {
			copy(expected[kept:], c.values[kept:])
			return row("panic", expected, c.values[:limit+1])
		}
		trace := append(make([]int, 0, len(c.values)), c.values...)
		return row(fmt.Sprint(kept), expected, trace)
	})
}
