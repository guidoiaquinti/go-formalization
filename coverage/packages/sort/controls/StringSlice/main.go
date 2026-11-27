package main

import (
	"fmt"
	"log"
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
		{"StringSlice.Len", 0, func() int { return lenFailures(sort.StringSlice.Len) }},
		{"Len counts capacity", 3, func() int { return lenFailures(func(xs sort.StringSlice) int { return cap(xs) }) }},
		{"StringSlice.Less", 0, func() int { return lessFailures(sort.StringSlice.Less) }},
		{"Less compares bytes as signed", 2, func() int { return lessFailures(signedByteLess) }},
		{"Less accepts equality", 1, func() int { return lessFailures(func(xs sort.StringSlice, i, j int) bool { return xs[i] <= xs[j] }) }},
		{"Less checks second index first", 1, func() int {
			return lessFailures(func(xs sort.StringSlice, i, j int) bool { b := xs[j]; return xs[i] < b })
		}},
		{"StringSlice.Swap", 0, func() int { return swapFailures(sort.StringSlice.Swap) }},
		{"Swap performs only the first write", 2, func() int { return swapFailures(func(xs sort.StringSlice, i, j int) { xs[i] = xs[j] }) }},
		{"Swap writes before validating both indices", 1, func() int {
			return swapFailures(func(xs sort.StringSlice, i, j int) {
				left := xs[i]
				xs[i] = ""
				xs[i] = xs[j]
				xs[j] = left
			})
		}},
		{"StringsAreSorted", 0, func() int { return sortedFailures(sort.StringsAreSorted) }},
		{"sortedness uses a strict comparison", 2, func() int { return sortedFailures(strictlySorted) }},
		{"sortedness checks only the endpoints", 1, func() int {
			return sortedFailures(func(xs []string) bool { return len(xs) < 2 || xs[0] <= xs[len(xs)-1] })
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

func lenFailures(length func(sort.StringSlice) int) int {
	backing := sort.StringSlice{"", "a", "ab", "\xff", "abc"}
	tests := []struct {
		input sort.StringSlice
		want  int
	}{
		{nil, 0}, {sort.StringSlice{}, 0}, {backing[1:4], 3},
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

func lessFailures(less func(sort.StringSlice, int, int) bool) int {
	tests := []struct {
		input     sort.StringSlice
		i, j      int
		want      bool
		panicText string
	}{
		{sort.StringSlice{"\x00", "\xff"}, 0, 1, true, ""},
		{sort.StringSlice{"\x00", "\xff"}, 1, 0, false, ""},
		{sort.StringSlice{"same", "same"}, 0, 1, false, ""},
		{sort.StringSlice{"same"}, -1, 2, false, "runtime error: index out of range [-1]"},
		{sort.StringSlice{"same"}, 0, 1, false, "runtime error: index out of range [1] with length 1"},
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

func swapFailures(swap func(sort.StringSlice, int, int)) int {
	tests := []struct {
		input     sort.StringSlice
		i, j      int
		want      sort.StringSlice
		panicText string
	}{
		{sort.StringSlice{"a", "\x00", "b", "\xff"}, 1, 3, sort.StringSlice{"a", "\xff", "b", "\x00"}, ""},
		{sort.StringSlice{"abc", "def"}, 1, 1, sort.StringSlice{"abc", "def"}, ""},
		{sort.StringSlice{"\x00", "\xff"}, 0, 2, sort.StringSlice{"\x00", "\xff"}, "runtime error: index out of range [2] with length 2"},
		{sort.StringSlice{"\x00", "\xff"}, -1, 3, sort.StringSlice{"\x00", "\xff"}, "runtime error: index out of range [-1]"},
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

func sortedFailures(sorted func([]string) bool) int {
	tests := []struct {
		input []string
		want  bool
	}{
		{nil, true}, {[]string{}, true}, {[]string{"a"}, true}, {[]string{"a", "a"}, true},
		{[]string{"", "a", "a", "\xff"}, true}, {[]string{"", "ab", "a", "\xff"}, false},
		{[]string{"\xff", "\x00"}, false},
	}
	mismatches := 0
	for _, test := range tests {
		if sorted(test.input) != test.want {
			mismatches++
		}
	}
	return mismatches
}

func strictlySorted(xs []string) bool {
	for i := 1; i < len(xs); i++ {
		if xs[i-1] >= xs[i] {
			return false
		}
	}
	return true
}

func signedByteLess(xs sort.StringSlice, i, j int) bool {
	a, b := xs[i], xs[j]
	for k := range min(len(a), len(b)) {
		if a[k] != b[k] {
			return int8(a[k]) < int8(b[k])
		}
	}
	return len(a) < len(b)
}
