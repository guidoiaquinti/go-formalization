package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"unsafe"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type words []int

type input struct {
	cells                  []int
	length, lo, hi, offset int
	alias, count           int
	independent            []int
}

func row(header string, cells []int) string {
	var out strings.Builder
	out.WriteString(header)
	for _, v := range cells {
		fmt.Fprintf(&out, " %d", v)
	}
	return out.String()
}

func actual(c input) (out string) {
	var base, view words
	if c.cells != nil {
		base = make(words, c.offset+len(c.cells))
		for i := range base[:c.offset] {
			base[i] = i + 101
		}
		copy(base[c.offset:], c.cells)
		view = base[c.offset : c.offset+c.length]
	}
	old := view[:cap(view)]
	src := c.independent
	if c.alias >= 0 {
		src = old[c.alias : c.alias+c.count]
	}
	checkGuards := func() {
		for i, v := range base[:c.offset] {
			if v != i+101 {
				log.Fatalf("Replace changed guard cell %d: %d", i, v)
			}
		}
	}
	defer func() {
		if p := recover(); p != nil {
			checkGuards()
			if _, ok := p.(runtime.Error); !ok || (c.lo >= 0 && c.lo <= c.hi && c.hi <= c.length) {
				panic(p)
			}
			if !slices.Equal(old, c.cells) {
				log.Fatal("Replace wrote to the backing region before rejecting the index")
			}
			out = "bounds panic"
		}
	}()
	result := slices.Replace(view, c.lo, c.hi, src...)
	checkGuards()
	allocated := unsafe.SliceData(view) != unsafe.SliceData(result)
	resultCells := result
	if allocated {
		if !slices.Equal(old, c.cells) {
			log.Fatal("Replace changed old storage on the allocating path")
		}
	} else {
		if cap(result) != cap(view) {
			log.Fatal("Replace changed capacity on the in-place path")
		}
		resultCells = result[:cap(result)]
	}
	return row(fmt.Sprintf("ok %d %t %t", len(result), result == nil, allocated), resultCells)
}

func expected(c input) string {
	if c.lo < 0 || c.hi < c.lo || c.hi > c.length {
		return "bounds panic"
	}
	source := c.independent
	if c.alias >= 0 {
		source = c.cells[c.alias : c.alias+c.count]
	}
	length := c.length - (c.hi - c.lo) + len(source)
	isNil := c.cells == nil && length == 0
	allocated := length > len(c.cells)
	cells := make([]int, max(length, len(c.cells)))
	copy(cells, c.cells)
	copy(cells[c.lo:], source)
	copy(cells[c.lo+len(source):], c.cells[c.hi:c.length])
	if length < c.length {
		clear(cells[length:c.length])
	}
	return row(fmt.Sprintf("ok %d %t %t", length, isNil, allocated), cells)
}

func main() {
	edges := []input{
		{alias: -1}, {lo: -1, alias: -1}, {hi: 1, alias: -1},
		{alias: -1, independent: []int{1}},
		{cells: []int{}, offset: 2, alias: -1},
		{cells: []int{1, 2, 3, 4}, length: 2, lo: 0, hi: 3, offset: 1, alias: -1},
		{cells: []int{1, 2, 3, 4}, length: 2, lo: 1, hi: 0, offset: 1, alias: -1},
		{cells: []int{1, 2, 3, 4}, length: 3, lo: 1, hi: 1, offset: 1, alias: 1, count: 1},
		{cells: []int{1, 2, 3, 4}, length: 3, lo: 1, hi: 3, offset: 1, alias: 2, count: 1},
		{cells: []int{1, 2, 3, 4}, length: 3, lo: 1, hi: 3, offset: 1, alias: 0, count: 4},
		{cells: []int{1, 2, 3, 4}, length: 3, lo: 0, hi: 1, offset: 1, alias: 1, count: 3},
		{cells: []int{1, 2, 3, 4, 5, 6}, length: 4, lo: 0, hi: 2, offset: 1, alias: 3, count: 1},
		{cells: []int{1, 2, 3, 4, 5, 6}, length: 4, lo: 0, hi: 2, offset: 1, alias: 3},
		{cells: []int{1, 2, 3}, length: 2, lo: 0, hi: 1, offset: 1, alias: 0, count: 2},
		{cells: []int{1, 2, 3}, length: 2, lo: 0, hi: 1, offset: 1, alias: 1, count: 2},
		{cells: []int{1, 2, 3, 4, 5}, length: 4, lo: 0, hi: 2, offset: 1, alias: 1, count: 3},
		{cells: []int{1, 2, 3, 4}, length: 3, lo: 0, hi: 2, offset: 1, alias: 1, count: 3},
		{cells: []int{1, 2, 3}, length: 3, lo: 0, hi: 3, alias: -1},
	}
	check.CompareSamples("godiff-slices-replace", edges, func(r *rand.Rand, i int) input {
		if i%43 == 0 {
			return input{hi: i%4 - 1, alias: -1, independent: []int{1}}
		}
		capacity := r.Intn(35)
		length := r.Intn(capacity + 1)
		cells := make([]int, capacity)
		for j := range cells {
			cells[j] = r.Intn(101) - 50
		}
		lo := r.Intn(length + 1)
		hi := lo + r.Intn(length-lo+1)
		c := input{cells: cells, length: length, lo: lo, hi: hi, offset: r.Intn(5), alias: -1}
		if i%7 == 0 {
			c.hi = length + 1
		}
		if i%2 == 0 {
			c.alias = r.Intn(capacity + 1)
			c.count = r.Intn(capacity - c.alias + 1)
		} else {
			c.independent = make([]int, r.Intn(17))
			for j := range c.independent {
				c.independent[j] = r.Intn(101) - 50
			}
		}
		return c
	}, func(c input) string {
		nils := 0
		if c.cells == nil {
			nils = 1
		}
		header := fmt.Sprintf("%d %d %d %d %d %d %d %d", strconv.IntSize, c.length, c.lo, c.hi, nils, len(c.cells), c.alias, c.count)
		return row(row(header, c.cells), c.independent)
	}, actual, expected)
}
