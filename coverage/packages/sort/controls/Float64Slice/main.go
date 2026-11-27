package main

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"sort"
)

const (
	negativeZero = 0x8000000000000000
	signalingNaN = 0x7ff0000000000001
	quietNaN     = 0xfff8000000000042
)

func main() {
	checkLen()
	checkLess()
	checkSwap()
	checkSorted()
}

func checkLen() {
	backing := values(0, negativeZero, signalingNaN, quietNaN, 0x7ff0000000000000)
	cases := []struct {
		input  sort.Float64Slice
		length int
	}{
		{nil, 0}, {sort.Float64Slice{}, 0}, {backing[1:4], 3},
		{backing[1:1], 0}, {backing[:2], 2}, {backing[1:4:4], 3},
	}
	for _, implementation := range []struct {
		name string
		run  func(sort.Float64Slice) int
		want int
	}{
		{"Float64Slice.Len", sort.Float64Slice.Len, 0},
		{"Len returns capacity", func(x sort.Float64Slice) int { return cap(x) }, 3},
	} {
		mismatches := 0
		for _, c := range cases {
			if implementation.run(c.input) != c.length {
				mismatches++
			}
		}
		assertCount(implementation.name, mismatches, implementation.want, len(cases))
	}
}

func checkLess() {
	cases := []struct {
		input     sort.Float64Slice
		i, j      int
		want      bool
		panicText string
	}{
		{values(signalingNaN, 0), 0, 1, true, ""},
		{values(0, quietNaN), 0, 1, false, ""},
		{values(quietNaN, signalingNaN), 0, 1, false, ""},
		{values(0, negativeZero), 0, 1, false, ""},
		{sort.Float64Slice{math.Inf(-1), math.Inf(1)}, 0, 1, true, ""},
		{values(0), -1, 2, false, "runtime error: index out of range [-1]"},
	}
	for _, implementation := range []struct {
		name string
		run  func(sort.Float64Slice, int, int) bool
		want int
	}{
		{"Float64Slice.Less", sort.Float64Slice.Less, 0},
		{"Less omits NaN handling", func(x sort.Float64Slice, i, j int) bool { return x[i] < x[j] }, 1},
		{"Less puts all NaNs before one another", func(x sort.Float64Slice, i, j int) bool { return x[i] < x[j] || math.IsNaN(x[i]) }, 1},
		{"Less treats equality as strict ordering", func(x sort.Float64Slice, i, j int) bool { return x[i] <= x[j] || math.IsNaN(x[i]) && !math.IsNaN(x[j]) }, 1},
		{"Less checks second index first", func(x sort.Float64Slice, i, j int) bool { _ = x[j]; return x.Less(i, j) }, 1},
	} {
		mismatches := 0
		for _, c := range cases {
			var got bool
			msg := observePanic(func() { got = implementation.run(c.input, c.i, c.j) })
			if got != c.want || msg != c.panicText {
				mismatches++
			}
		}
		assertCount(implementation.name, mismatches, implementation.want, len(cases))
	}
}

func checkSwap() {
	cases := []struct {
		input     []uint64
		i, j      int
		want      []uint64
		panicText string
	}{
		{[]uint64{signalingNaN, negativeZero, quietNaN}, 0, 2, []uint64{quietNaN, negativeZero, signalingNaN}, ""},
		{[]uint64{0, negativeZero}, 0, 1, []uint64{negativeZero, 0}, ""},
		{[]uint64{quietNaN}, 0, 0, []uint64{quietNaN}, ""},
		{[]uint64{signalingNaN, negativeZero}, 0, 2, []uint64{signalingNaN, negativeZero}, "runtime error: index out of range [2] with length 2"},
	}
	for _, implementation := range []struct {
		name string
		run  func(sort.Float64Slice, int, int)
		want int
	}{
		{"Float64Slice.Swap", sort.Float64Slice.Swap, 0},
		{"Swap canonicalizes NaN payloads", func(x sort.Float64Slice, i, j int) {
			x.Swap(i, j)
			if math.IsNaN(x[i]) {
				x[i] = math.NaN()
			}
			if math.IsNaN(x[j]) {
				x[j] = math.NaN()
			}
		}, 2},
		{"Swap writes before verifying the second index", func(x sort.Float64Slice, i, j int) {
			left := x[i]
			x[i] = 0
			x[i] = x[j]
			x[j] = left
		}, 1},
	} {
		mismatches := 0
		for _, c := range cases {
			got := values(c.input...)
			msg := observePanic(func() { implementation.run(got, c.i, c.j) })
			matches := msg == c.panicText && len(got) == len(c.want)
			for i := range c.want {
				matches = matches && math.Float64bits(got[i]) == c.want[i]
			}
			if !matches {
				mismatches++
			}
		}
		assertCount(implementation.name, mismatches, implementation.want, len(cases))
	}
}

func checkSorted() {
	cases := []struct {
		input []float64
		want  bool
	}{
		{nil, true},
		{values(signalingNaN, quietNaN, negativeZero, 0), true},
		{values(0, quietNaN), false},
		{values(quietNaN, signalingNaN), true},
		{[]float64{0, 2, 1, 3}, false},
		{values(0, negativeZero), true},
	}
	for _, implementation := range []struct {
		name string
		run  func([]float64) bool
		want int
	}{
		{"Float64sAreSorted", sort.Float64sAreSorted, 0},
		{"Sortedness omits NaN ordering", func(x []float64) bool {
			for i := 1; i < len(x); i++ {
				if x[i] < x[i-1] {
					return false
				}
			}
			return true
		}, 1},
		{"Sortedness rejects adjacent NaNs", func(x []float64) bool {
			for i := 1; i < len(x); i++ {
				if math.IsNaN(x[i]) && math.IsNaN(x[i-1]) {
					return false
				}
			}
			return sort.Float64sAreSorted(x)
		}, 2},
	} {
		mismatches := 0
		for _, c := range cases {
			if implementation.run(c.input) != c.want {
				mismatches++
			}
		}
		assertCount(implementation.name, mismatches, implementation.want, len(cases))
	}
}

func values(words ...uint64) sort.Float64Slice {
	out := make(sort.Float64Slice, len(words))
	for i, w := range words {
		out[i] = math.Float64frombits(w)
	}
	return out
}

func observePanic(call func()) (message string) {
	defer func() {
		if value := recover(); value != nil {
			if e, ok := value.(runtime.Error); ok {
				message = e.Error()
			} else {
				message = fmt.Sprintf("unexpected panic %T", value)
			}
		}
	}()
	call()
	return
}

func assertCount(name string, got, want, total int) {
	if got != want {
		log.Fatalf("%s disagreed on %d fixtures, want %d", name, got, want)
	}
	fmt.Printf("%s disagreed on %d of %d fixtures\n", name, got, total)
}
