package main

import (
	"fmt"
	"log"
	"runtime"
	"slices"
	"unsafe"
)

type words []int
type replace func(words, int, int, ...int) words

type fixture struct {
	length, lo, hi, from, to int
	source, want             []int
	allocated, pans          bool
}

func main() {
	cases := []fixture{
		{3, 0, 2, 1, 4, nil, []int{2, 3, 4, 3, 5, 6}, false, false},
		{4, 1, 3, -1, 0, []int{8}, []int{1, 8, 4, 0, 5, 6}, false, false},
		{4, 0, 2, 3, 4, nil, []int{4, 3, 4, 0, 5, 6}, false, false},
		{4, 2, 2, -1, 0, []int{9}, []int{1, 2, 9, 3, 4, 6}, false, false},
		{6, 2, 3, -1, 0, []int{7, 8, 9}, []int{1, 2, 7, 8, 9, 4, 5, 6}, true, false},
		{3, 0, 4, -1, 0, []int{8}, nil, false, true},
		{3, 4, 4, -1, 0, nil, nil, false, true},
	}
	for _, candidate := range []struct {
		name string
		run  replace
		want int
	}{
		{"Replace", slices.Replace[words], 0},
		{"Shift before copying an overlapping source", shiftFirst, 1},
		{"Leave obsolete cells uncleared", noClear, 2},
		{"Return early for an empty edit", skipBounds, 1},
	} {
		mismatches := 0
		for i, c := range cases {
			if !matches(c, candidate.run) {
				mismatches++
				if candidate.want == 0 {
					log.Fatalf("Go disagreed on fixture %d", i)
				}
			}
		}
		if mismatches != candidate.want {
			log.Fatalf("%s disagreed on %d fixtures, want %d", candidate.name, mismatches, candidate.want)
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", candidate.name, mismatches, len(cases))
	}
}

func matches(c fixture, run replace) bool {
	backing := words{19, 1, 2, 3, 4, 5, 6}
	before := slices.Clone(backing)
	view := backing[1 : 1+c.length]
	src := c.source
	if c.from >= 0 {
		src = view[:cap(view)][c.from:c.to]
	}
	var out words
	var caught any
	func() {
		defer func() { caught = recover() }()
		out = run(view, c.lo, c.hi, src...)
	}()
	if c.pans {
		_, bounds := caught.(runtime.Error)
		return bounds && slices.Equal(backing, before)
	}
	if caught != nil || backing[0] != before[0] || len(out) != c.length-(c.hi-c.lo)+len(src) ||
		(unsafe.SliceData(out) != unsafe.SliceData(view)) != c.allocated {
		return false
	}
	if c.allocated {
		return slices.Equal(out, c.want) && slices.Equal(backing, before)
	}
	return cap(out) == cap(view) && slices.Equal(out[:cap(out)], c.want)
}

func shiftFirst(s words, lo, hi int, v ...int) words {
	if lo >= 0 && lo < hi && hi <= len(s) && len(v) > hi-lo && len(s)-(hi-lo)+len(v) <= cap(s) {
		out := s[:len(s)-(hi-lo)+len(v)]
		copy(out[lo+len(v):], s[hi:])
		copy(out[lo:], v)
		return out
	}
	return slices.Replace(s, lo, hi, v...)
}

func noClear(s words, lo, hi int, v ...int) words {
	before := slices.Clone(s)
	out := slices.Replace(s, lo, hi, v...)
	if len(out) < len(s) {
		copy(s[len(out):], before[len(out):])
	}
	return out
}

func skipBounds(s words, lo, hi int, v ...int) words {
	if lo == hi && len(v) == 0 {
		return s
	}
	return slices.Replace(s, lo, hi, v...)
}
