package main

import (
	"fmt"
	"log"
	"runtime"
	"slices"
	"unsafe"
)

type words []int
type insert func(words, int, ...int) words

type fixture struct {
	length, index, sourceStart, sourceEnd int
	independent, want                     []int
	allocated, panic                      bool
}

func main() {
	tests := []fixture{
		{4, 1, 4, 6, nil, []int{1, 5, 6, 2, 3, 4, 7}, false, false},
		{4, 1, 0, 2, nil, []int{1, 1, 2, 2, 3, 4, 7}, false, false},
		{4, 4, 3, 5, nil, []int{1, 2, 3, 4, 4, 5, 7}, false, false},
		{7, 5, 1, 3, nil, []int{1, 2, 3, 4, 5, 2, 3, 6, 7}, true, false},
		{4, -1, -1, 0, nil, nil, false, true},
		{4, 5, -1, 0, nil, nil, false, true},
	}
	for _, candidate := range []struct {
		name string
		run  insert
		want int
	}{
		{"Insert", slices.Insert[words], 0},
		{"Shift before copying an overlapping source", shiftFirst, 1},
		{"Allocate for every nonempty insert", alwaysAllocate, 3},
		{"Ignore bounds when the source is empty", lateBounds, 2},
	} {
		mismatches := 0
		for i, c := range tests {
			if !check(c, candidate.run) {
				mismatches++
				if candidate.want == 0 {
					log.Fatalf("Go disagreed on fixture %d", i)
				}
			}
		}
		if mismatches != candidate.want {
			log.Fatalf("%s disagreed on %d fixtures, want %d", candidate.name, mismatches, candidate.want)
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", candidate.name, mismatches, len(tests))
	}
}

func check(c fixture, run insert) bool {
	backing := words{100, 1, 2, 3, 4, 5, 6, 7}
	original := slices.Clone(backing)
	view := backing[1 : 1+c.length]
	source := c.independent
	if c.sourceStart >= 0 {
		source = view[:cap(view)][c.sourceStart:c.sourceEnd]
	}
	var out words
	var caught any
	func() {
		defer func() { caught = recover() }()
		out = run(view, c.index, source...)
	}()
	if c.panic {
		_, bounds := caught.(runtime.Error)
		return bounds && slices.Equal(backing, original)
	}
	if caught != nil || len(out) != c.length+len(source) || backing[0] != original[0] ||
		(unsafe.SliceData(view) != unsafe.SliceData(out)) != c.allocated {
		return false
	}
	if c.allocated {
		return slices.Equal(backing, original) && slices.Equal(out, c.want)
	}
	return cap(out) == cap(view) && slices.Equal(out[:cap(out)], c.want)
}

func shiftFirst(s words, i int, v ...int) words {
	_ = s[i:]
	if len(s)+len(v) > cap(s) {
		return slices.Insert(s, i, v...)
	}
	result := s[:len(s)+len(v)]
	copy(result[i+len(v):], s[i:])
	copy(result[i:], v)
	return result
}

func alwaysAllocate(s words, i int, v ...int) words {
	_ = s[i:]
	if len(v) == 0 {
		return s
	}
	result := make(words, 0, len(s)+len(v))
	result = append(result, s[:i]...)
	result = append(result, v...)
	return append(result, s[i:]...)
}

func lateBounds(s words, i int, v ...int) words {
	if len(v) == 0 {
		return s
	}
	return slices.Insert(s, i, v...)
}
