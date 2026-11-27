package main

import (
	"fmt"
	"log"
	"slices"
	"sort"
)

type call struct {
	method string
	i, j   int
}

type recorder struct {
	length  int
	trace   []call
	compare func(int, int) bool
	panicOn string
	marker  *int
}

func (r *recorder) Len() int {
	r.log("Len", 0, 0)
	return r.length
}

func (r *recorder) Less(i, j int) bool {
	r.log("Less", i, j)
	return r.compare(i, j)
}

func (r *recorder) Swap(i, j int) {
	r.log("Swap", i, j)
}

func (r *recorder) log(method string, i, j int) {
	r.trace = append(r.trace, call{method, i, j})
	if r.panicOn == method {
		panic(r.marker)
	}
}

func main() {
	checkSorted()
	checkReverse()
}

func checkSorted() {
	marker := new(int)
	cases := []struct {
		responses []int
		n         int
		want      bool
		panic     bool
		trace     []call
	}{
		{nil, 0, true, false, []call{{"Len", 0, 0}}},
		{[]int{0, 0, 0}, 4, true, false, []call{{"Len", 0, 0}, {"Less", 3, 2}, {"Less", 2, 1}, {"Less", 1, 0}}},
		{[]int{2, 1, 0}, 4, false, false, []call{{"Len", 0, 0}, {"Less", 3, 2}, {"Less", 2, 1}}},
		{[]int{1, 2, 0}, 4, false, true, []call{{"Len", 0, 0}, {"Less", 3, 2}, {"Less", 2, 1}}},
	}
	implementations := []struct {
		name string
		run  func(sort.Interface) bool
		want int
	}{
		{"IsSorted", sort.IsSorted, 0},
		{"IsSorted scans forward", forward, 3},
		{"IsSorted continues after inversion", exhaustive, 1},
	}
	for _, candidate := range implementations {
		failures := 0
		for _, c := range cases {
			r := &recorder{length: c.n, marker: marker, compare: func(i, j int) bool {
				if i != j+1 || j < 0 || j >= len(c.responses) {
					panic("invalid Less pair")
				}
				switch c.responses[j] {
				case 2:
					panic(marker)
				case 1:
					return true
				default:
					return false
				}
			}}
			var result bool
			panicked := observe(func() { result = candidate.run(r) })
			var wantPanic any
			if c.panic {
				wantPanic = marker
			}
			if result != c.want || panicked != wantPanic || !slices.Equal(r.trace, c.trace) {
				failures++
			}
		}
		assertCount(candidate.name, failures, candidate.want, len(cases))
	}
}

func forward(data sort.Interface) bool {
	n := data.Len()
	for i := 1; i < n; i++ {
		if data.Less(i, i-1) {
			return false
		}
	}
	return true
}

func exhaustive(data sort.Interface) bool {
	sorted := true
	for i := data.Len() - 1; i > 0; i-- {
		if data.Less(i, i-1) {
			sorted = false
		}
	}
	return sorted
}

func checkReverse() {
	marker := new(int)
	implementations := []struct {
		name string
		run  func(sort.Interface) sort.Interface
		want int
	}{
		{"Reverse", sort.Reverse, 0},
		{"Reverse leaves Less indices in order", func(data sort.Interface) sort.Interface { return data }, 3},
		{"Reverse swaps the Swap indices", func(data sort.Interface) sort.Interface { return reversedSwap{sort.Reverse(data)} }, 2},
		{"Reverse calls Len during construction", func(data sort.Interface) sort.Interface { data.Len(); return sort.Reverse(data) }, 4},
	}
	for _, candidate := range implementations {
		failures := 0
		for _, method := range []string{"", "Len", "Less", "Swap"} {
			r := &recorder{length: 5, marker: marker, panicOn: method, compare: func(i, j int) bool { return i < j }}
			want := []call{{"Len", 0, 0}, {"Less", 1, 3}, {"Swap", 3, 1}}
			var wantPanic any
			for i, item := range want {
				if item.method == method {
					wantPanic = marker
					want = want[:i+1]
					break
				}
			}
			var length int
			var less, constructed bool
			panicked := observe(func() {
				reversed := candidate.run(r)
				constructed = true
				length = reversed.Len()
				less = reversed.Less(3, 1)
				reversed.Swap(3, 1)
			})
			if !constructed || panicked != wantPanic || !slices.Equal(r.trace, want) || method == "" && (length != 5 || !less) {
				failures++
			}
		}
		assertCount(candidate.name, failures, candidate.want, 4)
	}
}

type reversedSwap struct{ sort.Interface }

func (r reversedSwap) Swap(i, j int) { r.Interface.Swap(j, i) }

func observe(run func()) (caught any) {
	defer func() { caught = recover() }()
	run()
	return nil
}

func assertCount(name string, got, want, total int) {
	if got != want {
		log.Fatalf("%s disagreed on %d fixtures, want %d", name, got, want)
	}
	fmt.Printf("%s disagreed on %d of %d fixtures\n", name, got, total)
}
