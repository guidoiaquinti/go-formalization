package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"slices"
	"strings"
)

var panicMarker = new(int)

type searchResult struct {
	index int
	found bool
}

func result(index int, found bool) searchResult {
	return searchResult{index, found}
}

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		log.Fatalf("Go version %s, want %s", runtime.Version(), want)
	}
	checkIsSorted()
	checkIsSortedFunc()
	checkBinarySearch()
	checkBinarySearchFunc()
	fmt.Printf("%s %s/%s: sortedness and lower-bound controls passed and rejected their mutants\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

func checkIsSorted() {
	for _, c := range []struct {
		values []int
		want   bool
	}{
		{nil, true},
		{[]int{1, 1, 3}, true},
		{[]int{1, 3, 2}, false},
	} {
		if got := slices.IsSorted(c.values); got != c.want {
			log.Fatalf("IsSorted(%v) = %t, want %t", c.values, got, c.want)
		}
	}
	values := []int{1, 1, 3}
	if slices.IsSortedFunc(values, func(a, b int) int {
		if a <= b {
			return -1
		}
		return 1
	}) {
		log.Fatal("strict-increase mutant accepted duplicate elements")
	}
}

func checkIsSortedFunc() {
	values := []int{0, 1, 2, 3}
	observe := func(search func([]int, func(int, int) int) bool) (trace []int, panicked bool, sorted bool) {
		defer func() {
			if p := recover(); p != nil {
				if p != panicMarker {
					panic(p)
				}
				panicked = true
			}
		}()
		sorted = search(values, func(right, left int) int {
			if right != left+1 {
				log.Fatalf("nonadjacent or reversed operands (%d, %d)", right, left)
			}
			trace = append(trace, right)
			if right == 1 {
				panic(panicMarker)
			}
			return right - 3
		})
		return
	}
	trace, panicked, sorted := observe(slices.IsSortedFunc[[]int, int])
	if !slices.Equal(trace, []int{3, 2}) || panicked || sorted {
		log.Fatalf("IsSortedFunc produced trace=%v panic=%t sorted=%t", trace, panicked, sorted)
	}
	trace, panicked, sorted = observe(func(values []int, compare func(int, int) int) bool {
		for i := 1; i < len(values); i++ {
			if compare(values[i], values[i-1]) < 0 {
				return false
			}
		}
		return true
	})
	if !slices.Equal(trace, []int{1}) || !panicked || sorted {
		log.Fatalf("forward mutant produced trace=%v panic=%t sorted=%t", trace, panicked, sorted)
	}
}

func checkBinarySearch() {
	values := []int{1, 4, 4, 4, 7}
	if got := result(slices.BinarySearch(values, 4)); got != (searchResult{1, true}) {
		log.Fatalf("BinarySearch with duplicates returned %+v", got)
	}
	if got := result(slices.BinarySearch(values, 5)); got != (searchResult{4, false}) {
		log.Fatalf("BinarySearch missing key returned %+v", got)
	}
	immediateMatch := func(target int) searchResult {
		lo, hi := 0, len(values)
		for lo < hi {
			mid := lo + (hi-lo)/2
			switch {
			case values[mid] == target:
				return result(mid, true)
			case values[mid] < target:
				lo = mid + 1
			default:
				hi = mid
			}
		}
		return result(lo, false)
	}
	if got := immediateMatch(4); got != (searchResult{2, true}) {
		log.Fatalf("midpoint-equality mutant returned %+v, want later duplicate", got)
	}
	nan := math.NaN()
	if got := result(slices.BinarySearch([]float64{nan, nan, 0}, nan)); got != (searchResult{0, true}) {
		log.Fatalf("BinarySearch NaN key returned %+v", got)
	}
}

func checkBinarySearchFunc() {
	type key struct{ value int }
	compare := func(element int, target key) int {
		switch {
		case element > target.value:
			return -7
		case element < target.value:
			return 13
		default:
			return 0
		}
	}
	values := []int{9, 7, 7, 4}
	if got := result(slices.BinarySearchFunc(values, key{7}, compare)); got != (searchResult{1, true}) {
		log.Fatalf("BinarySearchFunc with duplicates returned %+v", got)
	}
	if got := result(slices.BinarySearchFunc(values, key{6}, compare)); got != (searchResult{3, false}) {
		log.Fatalf("BinarySearchFunc missing key returned %+v", got)
	}
	mutant := func(element int, target key) int {
		order := compare(element, target)
		if order < 0 && order != -1 {
			return 1
		}
		return order
	}
	if got := result(slices.BinarySearchFunc(values, key{7}, mutant)); got != (searchResult{0, false}) {
		log.Fatalf("exact-minus-one mutant returned %+v, want wrong index zero", got)
	}
}
