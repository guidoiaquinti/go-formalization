package main

import (
	"fmt"
	"log"
	"math"
	"slices"
	"sort"
)

func main() {
	checkSearch()
	checkFind()
}

func checkSearch() {
	cases := []struct {
		n, cut int
		trace  []int
	}{
		{-1, 0, nil},
		{0, 0, nil},
		{1, 0, []int{0}},
		{5, 2, []int{2, 1}},
		{5, 5, []int{2, 4}},
		{math.MaxInt, math.MaxInt, nil},
	}
	candidates := []struct {
		name string
		run  func(int, func(int) bool) int
		want int
	}{
		{"Search", sort.Search, 0},
		{"Search computes a signed midpoint", signedMidpoint, 1},
		{"Search returns the last false index", lastFalse, 3},
	}
	for _, candidate := range candidates {
		mismatches := 0
		for _, c := range cases {
			var trace []int
			got, completed := safe(func() int {
				return candidate.run(c.n, func(i int) bool {
					if i < 0 || i >= c.n {
						panic("invalid callback index")
					}
					trace = append(trace, i)
					return i >= c.cut
				})
			})
			if !completed || got != c.cut || c.n != math.MaxInt && !slices.Equal(trace, c.trace) {
				mismatches++
			}
		}
		assertCount(candidate.name, mismatches, candidate.want, len(cases))
	}
}

func checkFind() {
	cases := []struct {
		n, cut, endMatch int
		trace            []int
	}{
		{-1, 0, 0, nil},
		{0, 0, 0, nil},
		{1, 0, 1, []int{0, 0}},
		{5, 1, 4, []int{2, 1, 0, 1}},
		{5, 2, 2, []int{2, 1, 2}},
		{5, 5, 5, []int{2, 4}},
	}
	candidates := []struct {
		name string
		run  func(int, func(int) int) (int, bool)
		want int
	}{
		{"Find", sort.Find, 0},
		{"Find skips the final equality check", skipEquality, 3},
		{"Find compares a position at the end", compareAtEnd, 3},
	}
	for _, candidate := range candidates {
		mismatches := 0
		for _, c := range cases {
			var trace []int
			var found bool
			got, completed := safe(func() int {
				var index int
				index, found = candidate.run(c.n, func(i int) int {
					if i < 0 || i >= c.n {
						panic("invalid callback index")
					}
					trace = append(trace, i)
					if i < c.cut {
						return math.MaxInt
					}
					if i < c.endMatch {
						return 0
					}
					return math.MinInt
				})
				return index
			})
			if !completed || got != c.cut || found != (c.cut < c.endMatch) || !slices.Equal(trace, c.trace) {
				mismatches++
			}
		}
		assertCount(candidate.name, mismatches, candidate.want, len(cases))
	}
}

func signedMidpoint(n int, predicate func(int) bool) int {
	lo, hi := 0, n
	for lo < hi {
		middle := (lo + hi) / 2
		if predicate(middle) {
			hi = middle
		} else {
			lo = middle + 1
		}
	}
	return lo
}

func lastFalse(n int, predicate func(int) bool) int {
	return max(0, sort.Search(n, predicate)-1)
}

func skipEquality(n int, compare func(int) int) (int, bool) {
	index := sort.Search(n, func(i int) bool { return compare(i) <= 0 })
	return index, index < n
}

func compareAtEnd(n int, compare func(int) int) (int, bool) {
	index := sort.Search(n, func(i int) bool { return compare(i) <= 0 })
	return index, compare(index) == 0
}

func safe(run func() int) (result int, completed bool) {
	defer func() {
		if caught := recover(); caught != nil {
			if caught != "invalid callback index" {
				panic(caught)
			}
			completed = false
		}
	}()
	return run(), true
}

func assertCount(name string, got, want, total int) {
	if got != want {
		log.Fatalf("%s disagreed on %d fixtures, want %d", name, got, want)
	}
	fmt.Printf("%s disagreed on %d of %d fixtures\n", name, got, total)
}
