package main

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"slices"
	"sort"
)

type rankedCheck struct {
	name  string
	want  int
	count func() int
}

func main() {
	checks := []rankedCheck{
		{"IntSlice.Len", 0, func() int { return lenFailures(sort.IntSlice.Len) }},
		{"Len counts capacity", 3, func() int { return lenFailures(func(xs sort.IntSlice) int { return cap(xs) }) }},
		{"IntSlice.Less", 0, func() int { return lessFailures(sort.IntSlice.Less) }},
		{"Less accepts equality", 1, func() int { return lessFailures(func(xs sort.IntSlice, i, j int) bool { return xs[i] <= xs[j] }) }},
		{"Less checks second index first", 1, func() int {
			return lessFailures(func(xs sort.IntSlice, i, j int) bool { b := xs[j]; return xs[i] < b })
		}},
		{"IntSlice.Swap", 0, func() int { return swapFailures(sort.IntSlice.Swap) }},
		{"Swap performs only the first write", 2, func() int { return swapFailures(func(xs sort.IntSlice, i, j int) { xs[i] = xs[j] }) }},
		{"Swap writes before validating both indices", 1, func() int {
			return swapFailures(func(xs sort.IntSlice, i, j int) {
				left := xs[i]
				xs[i] = 0
				xs[i] = xs[j]
				xs[j] = left
			})
		}},
		{"IntsAreSorted", 0, func() int { return sortedFailures(sort.IntsAreSorted) }},
		{"sortedness uses a strict comparison", 2, func() int { return sortedFailures(strictlySorted) }},
		{"sortedness checks only the endpoints", 1, func() int {
			return sortedFailures(func(xs []int) bool { return len(xs) < 2 || xs[0] <= xs[len(xs)-1] })
		}},
	}
	for _, check := range checks {
		got := check.count()
		if got != check.want {
			log.Fatalf("%s disagreed on %d fixtures, want %d", check.name, got, check.want)
		}
		fmt.Printf("%s disagreed on %d fixtures\n", check.name, got)
	}
}

func lenFailures(length func(sort.IntSlice) int) int {
	backing := sort.IntSlice{1, 2, 3, 4, 5}
	tests := []struct {
		input sort.IntSlice
		want  int
	}{
		{nil, 0}, {sort.IntSlice{}, 0}, {backing[1:4], 3},
		{backing[1:1], 0}, {backing[:2], 2}, {backing[1:4:4], 3},
	}
	mismatches := 0
	for _, test := range tests {
		if length(test.input) != test.want {
			mismatches++
		}
	}
	return mismatches
}

func lessFailures(less func(sort.IntSlice, int, int) bool) int {
	tests := []struct {
		input     sort.IntSlice
		i, j      int
		want      bool
		panicText string
	}{
		{sort.IntSlice{math.MinInt, math.MaxInt}, 0, 1, true, ""},
		{sort.IntSlice{math.MinInt, math.MaxInt}, 1, 0, false, ""},
		{sort.IntSlice{4, 4}, 0, 1, false, ""},
		{sort.IntSlice{1}, -1, 2, false, "runtime error: index out of range [-1]"},
		{sort.IntSlice{1}, 0, 1, false, "runtime error: index out of range [1] with length 1"},
	}
	mismatches := 0
	for _, test := range tests {
		var got bool
		panicText := observePanic(func() { got = less(test.input, test.i, test.j) })
		if got != test.want || panicText != test.panicText {
			mismatches++
		}
	}
	return mismatches
}

func swapFailures(swap func(sort.IntSlice, int, int)) int {
	tests := []struct {
		input     sort.IntSlice
		i, j      int
		want      sort.IntSlice
		panicText string
	}{
		{sort.IntSlice{6, -10, 40, 2}, 1, 3, sort.IntSlice{6, 2, 40, -10}, ""},
		{sort.IntSlice{9, 4}, 1, 1, sort.IntSlice{9, 4}, ""},
		{sort.IntSlice{2, 8}, 0, 2, sort.IntSlice{2, 8}, "runtime error: index out of range [2] with length 2"},
		{sort.IntSlice{2, 8}, -1, 3, sort.IntSlice{2, 8}, "runtime error: index out of range [-1]"},
	}
	mismatches := 0
	for _, test := range tests {
		got := slices.Clone(test.input)
		panicText := observePanic(func() { swap(got, test.i, test.j) })
		if !slices.Equal(got, test.want) || panicText != test.panicText {
			mismatches++
		}
	}
	return mismatches
}

func observePanic(call func()) (message string) {
	defer func() {
		if value := recover(); value != nil {
			if err, ok := value.(runtime.Error); ok {
				message = err.Error()
			} else {
				message = fmt.Sprintf("unexpected panic %T", value)
			}
		}
	}()
	call()
	return
}

func sortedFailures(sorted func([]int) bool) int {
	tests := []struct {
		input []int
		want  bool
	}{
		{nil, true}, {[]int{}, true}, {[]int{1}, true}, {[]int{1, 1}, true},
		{[]int{math.MinInt, 0, 0, math.MaxInt}, true}, {[]int{0, 2, 1, 3}, false},
		{[]int{1, 0}, false},
	}
	mismatches := 0
	for _, test := range tests {
		if sorted(test.input) != test.want {
			mismatches++
		}
	}
	return mismatches
}

func strictlySorted(xs []int) bool {
	for i := 1; i < len(xs); i++ {
		if xs[i-1] >= xs[i] {
			return false
		}
	}
	return true
}
