package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"slices"
	"strings"
)

var panicMarker = new(int)

type callbackResult struct {
	value int
	calls int
	panic bool
}

func trace(responses []int, search func([]int, func(int) bool) int) (out callbackResult) {
	xs := make([]int, len(responses))
	for i := range xs {
		xs[i] = i
	}
	defer func() {
		if p := recover(); p != nil {
			if p != panicMarker {
				panic(p)
			}
			out.panic = true
		}
	}()
	out.value = search(xs, func(i int) bool {
		if i != out.calls {
			panic("callback called out of order")
		}
		out.calls++
		if responses[i] == -1 {
			panic(panicMarker)
		}
		return responses[i] == 1
	})
	return
}

func asContains(search func([]int, func(int) bool) int) func([]int, func(int) bool) int {
	return func(xs []int, predicate func(int) bool) int {
		if search(xs, predicate) >= 0 {
			return 1
		}
		return 0
	}
}

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		log.Fatalf("Go version %s, want %s", runtime.Version(), want)
	}

	valueCases := []struct {
		xs   []int
		key  int
		want int
	}{
		{nil, 9, -1},
		{[]int{}, 9, -1},
		{[]int{9}, 9, 0},
		{[]int{1, 9, 0, 9}, 9, 1},
		{[]int{1, 2}, 9, -1},
	}
	var indexFailures, containsFailures int
	for _, c := range valueCases {
		if got := slices.Index(c.xs, c.key); got != c.want {
			log.Fatalf("Index(%v, %d) = %d, want %d", c.xs, c.key, got, c.want)
		}
		if got := slices.Contains(c.xs, c.key); got != (c.want >= 0) {
			log.Fatalf("Contains(%v, %d) = %t, want %t", c.xs, c.key, got, c.want >= 0)
		}
		last := -1
		for i, x := range c.xs {
			if x == c.key {
				last = i
			}
		}
		if last != c.want {
			indexFailures++
		}
		if (slices.Index(c.xs, c.key) > 0) != (c.want >= 0) {
			containsFailures++
		}
	}

	callbackCases := []struct {
		responses       []int
		index, contains callbackResult
	}{
		{nil, callbackResult{-1, 0, false}, callbackResult{0, 0, false}},
		{[]int{0, 0}, callbackResult{-1, 2, false}, callbackResult{0, 2, false}},
		{[]int{0, -1, 1}, callbackResult{0, 2, true}, callbackResult{0, 2, true}},
		{[]int{0, 1, -1}, callbackResult{1, 2, false}, callbackResult{1, 2, false}},
		{[]int{1, 0, 1}, callbackResult{0, 1, false}, callbackResult{1, 1, false}},
	}
	contains := func(xs []int, predicate func(int) bool) int {
		if slices.ContainsFunc(xs, predicate) {
			return 1
		}
		return 0
	}
	eager := func(xs []int, predicate func(int) bool) int {
		first := -1
		for i, x := range xs {
			if predicate(x) && first < 0 {
				first = i
			}
		}
		return first
	}
	var indexFuncFailures, containsFuncFailures int
	for _, c := range callbackCases {
		if got := trace(c.responses, slices.IndexFunc[[]int, int]); got != c.index {
			log.Fatalf("IndexFunc script %v gave %+v, want %+v", c.responses, got, c.index)
		}
		if got := trace(c.responses, contains); got != c.contains {
			log.Fatalf("ContainsFunc script %v gave %+v, want %+v", c.responses, got, c.contains)
		}
		if trace(c.responses, eager) != c.index {
			indexFuncFailures++
		}
		if trace(c.responses, asContains(eager)) != c.contains {
			containsFuncFailures++
		}
	}
	failures := [4]int{indexFailures, containsFailures, indexFuncFailures, containsFuncFailures}
	if want := [4]int{1, 1, 2, 2}; failures != want {
		log.Fatalf("search mutant failure counts %v, want %v", failures, want)
	}
	fmt.Printf("%s %s/%s: five value fixtures and five callback fixtures passed; search mutant failures %v\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, failures)
}
