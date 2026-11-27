package main

import (
	"fmt"
	"log"
	"slices"
	"unsafe"
)

type named []int

type test struct {
	makeInput            func() named
	growth               int
	length               int
	nilResult, allocated bool
	cells                []int
	wantPanic            string
}

func main() {
	nilInput := func() named { return nil }
	empty := func() named { return named{} }
	view := func() named { return named{31, 11, 22, 77, 88}[1:3] }
	tests := []test{
		{nilInput, 0, 0, true, false, nil, ""},
		{nilInput, 1, 0, false, true, []int{0}, ""},
		{empty, 0, 0, false, false, nil, ""},
		{view, 2, 2, false, false, []int{11, 22, 77, 88}, ""},
		{view, 3, 2, false, true, []int{11, 22, 77, 88, 0}, ""},
		{view, -1, 0, false, false, nil, "cannot be negative"},
	}
	candidates := []struct {
		name string
		run  func(named, int) named
		want int
	}{
		{"Grow", slices.Grow[named], 0},
		{"Copy only visible cells", visibleOnly, 1},
		{"Allocate even when the existing capacity fits", alwaysAllocate, 1},
		{"Accept a negative growth count", acceptNegative, 1},
	}
	for _, candidate := range candidates {
		failures := 0
		for i, c := range tests {
			if !matches(candidate.run, c) {
				failures++
				if candidate.want == 0 {
					log.Fatalf("Go disagreed on fixture %d", i)
				}
			}
		}
		if failures != candidate.want {
			log.Fatalf("%s disagreed on %d fixtures, want %d", candidate.name, failures, candidate.want)
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", candidate.name, failures, len(tests))
	}
	fmt.Printf("Grow(nil, 1) returned nil: %t\n", slices.Grow(named(nil), 1) == nil)
}

func matches(grow func(named, int) named, c test) bool {
	input := c.makeInput()
	before := slices.Clone(input[:cap(input)])
	var caught any
	var out named
	func() {
		defer func() { caught = recover() }()
		out = grow(input, c.growth)
	}()
	if !slices.Equal(input[:cap(input)], before) {
		return false
	}
	if c.wantPanic != "" {
		return caught == c.wantPanic
	}
	if caught != nil || len(out) != c.length || (out == nil) != c.nilResult || cap(out) < len(c.cells) {
		return false
	}
	if (unsafe.SliceData(out) != unsafe.SliceData(input)) != c.allocated {
		return false
	}
	return slices.Equal(out[:len(c.cells)], c.cells)
}

func visibleOnly(input named, n int) named {
	if n < 0 || n <= cap(input)-len(input) {
		return slices.Grow(input, n)
	}
	out := make(named, len(input), len(input)+n)
	copy(out, input)
	return out
}

func alwaysAllocate(input named, n int) named {
	if n <= 0 || n > cap(input)-len(input) {
		return slices.Grow(input, n)
	}
	out := make(named, cap(input))
	copy(out, input[:cap(input)])
	return out[:len(input)]
}

func acceptNegative(input named, n int) named {
	return slices.Grow(input, max(n, 0))
}
