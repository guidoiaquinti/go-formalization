package main

import (
	"fmt"
	"log"
	"math/rand"
	"slices"
	"strings"
	"unsafe"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type words []int

type input struct {
	length, offset, growth int
	cells                  []int
}

func output(header string, cells []int) string {
	var row strings.Builder
	row.WriteString(header)
	for _, value := range cells {
		fmt.Fprintf(&row, " %d", value)
	}
	return row.String()
}

func actual(c input) (row string) {
	var old, backing words
	if c.cells != nil {
		backing = make(words, c.offset+len(c.cells))
		for i := range backing[:c.offset] {
			backing[i] = i + 1
		}
		copy(backing[c.offset:], c.cells)
		old = backing[c.offset : c.offset+c.length]
	}
	checkOld := func() {
		for i, got := range backing[:c.offset] {
			if got != i+1 {
				log.Fatalf("Grow changed guard cell %d: %d", i, got)
			}
		}
		if !slices.Equal(old[:cap(old)], c.cells) {
			log.Fatal("Grow changed the original backing region")
		}
	}
	defer func() {
		if p := recover(); p != nil {
			checkOld()
			row = fmt.Sprint("panic ", p)
			if msg, ok := p.(string); !ok || msg != "cannot be negative" {
				log.Fatalf("unexpected Grow panic %T: %v", p, p)
			}
		}
	}()
	result := slices.Grow(old, c.growth)
	checkOld()
	if c.growth < 0 {
		log.Fatal("Grow unexpectedly returned for a negative growth count")
	}
	required := max(cap(old), c.length+c.growth)
	if len(result) != c.length || cap(result) < required {
		log.Fatalf("Grow returned len=%d cap=%d, expected len=%d cap≥%d", len(result), cap(result), c.length, required)
	}
	allocated := unsafe.SliceData(old) != unsafe.SliceData(result)
	if required == cap(old) && cap(result) != cap(old) {
		log.Fatal("Grow changed capacity with enough spare storage")
	}
	return output(fmt.Sprintf("ok %d %t %t", len(result), result == nil, allocated), result[:required])
}

func expected(c input) string {
	if c.growth < 0 {
		return "panic cannot be negative"
	}
	capacity := c.length + c.growth
	if capacity <= len(c.cells) {
		return output(fmt.Sprintf("ok %d %t false", c.length, c.cells == nil), c.cells)
	}
	cells := make([]int, capacity)
	copy(cells, c.cells)
	return output(fmt.Sprintf("ok %d false true", c.length), cells)
}

func main() {
	edges := []input{{0, 0, 0, nil}, {0, 0, 1, nil}, {0, 0, -1, nil}, {0, 0, 0, []int{}}, {0, 0, 1, []int{}},
		{0, 3, 0, []int{}}, {0, 3, 1, []int{}},
		{0, 2, 2, []int{11, 22}}, {0, 2, 3, []int{11, 22}},
		{2, 3, 0, []int{1, 2, -3, -4}}, {2, 3, 1, []int{1, 2, -3, -4}},
		{2, 3, 2, []int{1, 2, -3, -4}}, {2, 3, 3, []int{1, 2, -3, -4}},
		{2, 3, -1, []int{1, 2, -3, -4}}, {3, 2, 1, []int{1, 2, 3}},
		{1, 0, 4096, []int{1, 2, 3}},
	}
	check.CompareSamples("godiff-slices-grow", edges, func(r *rand.Rand, i int) input {
		if i%41 == 0 {
			return input{0, 0, i%3 - 1, nil}
		}
		capacity := r.Intn(65)
		length := r.Intn(capacity + 1)
		cells := make([]int, capacity)
		for j := range cells {
			cells[j] = r.Intn(2001) - 1000
		}
		slack := capacity - length
		growth := slack + i%4 - 1
		if i%7 == 0 {
			growth = r.Intn(129) - 4
		}
		return input{length, r.Intn(5), growth, cells}
	}, func(c input) string {
		isNil := 0
		if c.cells == nil {
			isNil = 1
		}
		return output(fmt.Sprintf("%d %d %d", c.growth, c.length, isNil), c.cells)
	}, actual, expected)
}
