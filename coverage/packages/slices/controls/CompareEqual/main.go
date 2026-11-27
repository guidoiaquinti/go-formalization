package main

import (
	"cmp"
	"fmt"
	"log"
	"os"
	"runtime"
	"slices"
	"strings"
)

var panicMarker = new(int)

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		log.Fatalf("Go version %s, want %s", runtime.Version(), want)
	}

	valueCases := []struct {
		xs, ys []int
		want   int
	}{
		{nil, []int{}, 0},
		{[]int{1}, []int{1, 2}, -1},
		{[]int{3}, []int{1, 2}, 1},
		{[]int{3, 1}, []int{3, 2}, -1},
		{[]int{3, 1}, []int{3, 1}, 0},
	}
	var prefixEqualFailures, lengthFirstFailures int
	for _, c := range valueCases {
		if got := slices.Equal(c.xs, c.ys); got != (c.want == 0) {
			log.Fatalf("Equal(%v, %v) = %t, want %t", c.xs, c.ys, got, c.want == 0)
		}
		if got := slices.Compare(c.xs, c.ys); got != c.want {
			log.Fatalf("Compare(%v, %v) = %d, want %d", c.xs, c.ys, got, c.want)
		}
		overlap := min(len(c.xs), len(c.ys))
		if slices.Equal(c.xs[:overlap], c.ys[:overlap]) != (c.want == 0) {
			prefixEqualFailures++
		}
		lengthFirst := cmp.Compare(len(c.xs), len(c.ys))
		if lengthFirst == 0 {
			lengthFirst = slices.Compare(c.xs, c.ys)
		}
		if lengthFirst != c.want {
			lengthFirstFailures++
		}
	}

	equalFuncFailures := checkEqualFunc()
	compareFuncFailures := checkCompareFunc()
	failures := [4]int{prefixEqualFailures, lengthFirstFailures, equalFuncFailures, compareFuncFailures}
	if want := [4]int{1, 1, 1, 2}; failures != want {
		log.Fatalf("mutant failure counts %v, want %v", failures, want)
	}
	fmt.Printf("%s %s/%s: value and callback comparisons passed; mutant failures %v\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, failures)
}

func checkEqualFunc() int {
	xs, ys := []int{2, 5}, []string{"two"}
	calls := 0
	predicate := func(int, string) bool {
		calls++
		return true
	}
	if got := slices.EqualFunc(xs, ys, predicate); got || calls != 0 {
		log.Fatalf("unequal lengths produced %t and %d calls, want false and no calls", got, calls)
	}
	overlap := min(len(xs), len(ys))
	slices.EqualFunc(xs[:overlap], ys[:overlap], predicate)
	mutantFailures := 0
	if calls != 0 {
		mutantFailures++
	}
	calls = 0
	if got := slices.EqualFunc(xs, []string{"two", "five"}, func(a int, b string) bool {
		calls++
		if a != 2 || b != "two" {
			log.Fatalf("unexpected equality callback arguments (%d, %q)", a, b)
		}
		return false
	}); got || calls != 1 {
		log.Fatalf("first false comparison produced %t and %d calls, want false and one call", got, calls)
	}
	checkPanic(func() { slices.EqualFunc(xs, []string{"two", "five"}, func(int, string) bool { panic(panicMarker) }) })
	return mutantFailures
}

func checkCompareFunc() int {
	var normalizedFailures int
	for _, want := range []int{-17, 42} {
		calls := 0
		got := slices.CompareFunc([]int{2, 5, 7}, []string{"two", "five"}, func(a int, b string) int {
			calls++
			if calls != 1 || a != 2 || b != "two" {
				log.Fatalf("unexpected comparison callback arguments (%d, %q) on call %d", a, b, calls)
			}
			return want
		})
		if got != want || calls != 1 {
			log.Fatalf("nonzero comparison produced %d and %d calls, want %d and one call", got, calls, want)
		}
		if cmp.Compare(got, 0) != want {
			normalizedFailures++
		}
	}
	calls := 0
	if got := slices.CompareFunc([]int{2}, []string{"two", "five"}, func(int, string) int {
		calls++
		return 0
	}); got != -1 || calls != 1 {
		log.Fatalf("all-tie comparison produced %d and %d calls, want -1 and one call", got, calls)
	}
	checkPanic(func() {
		slices.CompareFunc([]int{2}, []string{"two"}, func(int, string) int { panic(panicMarker) })
	})
	return normalizedFailures
}

func checkPanic(call func()) {
	defer func() {
		if got := recover(); got != panicMarker {
			log.Fatalf("callback panic = %v, want original payload", got)
		}
	}()
	call()
}
