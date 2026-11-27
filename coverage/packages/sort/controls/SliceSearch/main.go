package main

import (
	"cmp"
	"fmt"
	"log"
	"math"
	"sort"
)

type fixture[T any] struct {
	values []T
	key    T
	want   int
}

func main() {
	ints := []fixture[int]{
		{nil, 0, 0},
		{[]int{-4, 0, 0, 8}, 0, 1},
		{[]int{-4, 0, 0, 8}, 1, 3},
		{[]int{-4, 0, 0, 8}, 9, 4},
		{[]int{math.MinInt, math.MaxInt}, math.MinInt, 0},
		{[]int{math.MinInt, math.MaxInt}, math.MaxInt, 1},
	}
	run("IntSlice.Search", func(values []int, key int) int { return sort.IntSlice(values).Search(key) }, ints, 0)
	run("IntSlice.Search uses strict comparison", linear(func(value, key int) bool { return value > key }), ints, 3)

	strs := []fixture[string]{
		{nil, "x", 0},
		{[]string{"", "a", "a", "ab", "b"}, "a", 1},
		{[]string{"", "a", "a", "ab", "b"}, "aa", 3},
		{[]string{"a", "ab"}, "ab", 1},
		{[]string{"a", "ab"}, "abc", 2},
		{[]string{"\x00", "\x80", "\xff"}, "\x80", 1},
	}
	run("StringSlice.Search", func(values []string, key string) int { return sort.StringSlice(values).Search(key) }, strs, 0)
	run("StringSlice.Search uses strict comparison", linear(func(value, key string) bool { return value > key }), strs, 3)
	run("StringSlice.Search compares lengths", linear(func(value, key string) bool { return len(value) >= len(key) }), strs, 1)

	signaling := math.Float64frombits(0x7ff0000000000001)
	quiet := math.Float64frombits(0xfff80000000000ab)
	negativeZero := math.Copysign(0, -1)
	floats := []fixture[float64]{
		{nil, quiet, 0},
		{[]float64{signaling, quiet, -1, negativeZero, 0, 1}, signaling, 6},
		{[]float64{signaling, quiet, -1, negativeZero, 0, 1}, 0, 3},
		{[]float64{signaling, quiet, math.Inf(-1), math.Inf(1)}, math.Inf(-1), 2},
		{[]float64{signaling, quiet, math.Inf(-1), math.Inf(1)}, math.Inf(1), 3},
		{[]float64{1, 1, 2}, 1, 0},
	}
	run("Float64Slice.Search", func(values []float64, key float64) int { return sort.Float64Slice(values).Search(key) }, floats, 0)
	run("Float64Slice.Search uses strict comparison", linear(func(value, key float64) bool { return value > key }), floats, 4)
	run("Float64Slice.Search uses NaN-first ordering",
		linear(func(value, key float64) bool { return !cmp.Less(value, key) }), floats, 1)
}

func run[T any](name string, search func([]T, T) int, tests []fixture[T], want int) {
	mismatches := 0
	for _, test := range tests {
		if result := search(test.values, test.key); result != test.want {
			mismatches++
			if want == 0 {
				log.Fatalf("%s(%v, %v) = %d, want %d", name, test.values, test.key, result, test.want)
			}
		}
	}
	if mismatches != want {
		log.Fatalf("%s disagreed on %d fixtures, want %d", name, mismatches, want)
	}
	fmt.Printf("%s disagreed on %d of %d fixtures\n", name, mismatches, len(tests))
}

func linear[T any](qualifies func(T, T) bool) func([]T, T) int {
	return func(values []T, key T) int {
		for i, value := range values {
			if qualifies(value, key) {
				return i
			}
		}
		return len(values)
	}
}
