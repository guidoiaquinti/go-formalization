package main

import (
	"cmp"
	"fmt"
	"log"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func require(ok bool, label string) {
	if !ok {
		log.Fatal(label)
	}
}
func catch(call func()) (result any) {
	defer func() { result = recover() }()
	call()
	return
}

func main() {
	// A non-NaN and three distinct NaN encodings also exercise exact movement.
	words := []uint64{0x8000000000000000, 0x7ff0000000000001, 0, 0xfff8000000000002, 0x7ff8000000000003, 0x3ff0000000000000}
	data := make([]float64, len(words))
	for i, v := range words {
		data[i] = math.Float64frombits(v)
	}
	slices.Sort(data)
	want := []uint64{words[1], words[3], words[4], words[0], words[2], words[5]}
	for i, v := range data {
		require(math.Float64bits(v) == want[i], "exact NaN and zero ordering")
	}

	compareMarker := new(int)
	values := []int{3, 2, 1}
	var calls [][2]int
	p := catch(func() {
		slices.SortFunc(values, func(a, b int) int {
			calls = append(calls, [2]int{a, b})
			if len(calls) == 2 {
				panic(compareMarker)
			}
			return cmp.Compare(a, b)
		})
	})
	require(p == compareMarker && slices.Equal(values, []int{2, 3, 1}) && slices.Equal(calls, [][2]int{{2, 3}, {1, 3}}), "comparator payload, trace and retained prior swap")
	require(catch(func() { slices.SortFunc([]int{2}, nil) }) == nil, "singleton invokes nil comparator")
	require(catch(func() { slices.SortFunc([]int{2, 1}, nil) }) != nil, "pair skips nil comparator")

	producerMarker := new(byte)
	acknowledged, comparisons := 0, 0
	p = catch(func() {
		slices.SortedFunc(func(yield func(int) bool) {
			for _, v := range []int{3, 2, 1} {
				require(yield(v), "early producer stop")
				acknowledged++
			}
			panic(producerMarker)
		}, func(a, b int) int { comparisons++; return cmp.Compare(a, b) })
	})
	require(p == producerMarker && acknowledged == 3 && comparisons == 0, "sorting started before producer panic")
	finished := false
	output := slices.SortedFunc(func(yield func(int) bool) {
		for _, v := range []int{3, 1, 2} {
			require(yield(v), "early producer stop")
		}
		finished = true
	}, func(a, b int) int { require(finished, "comparison preceded producer return"); return cmp.Compare(a, b) })
	require(slices.Equal(output, []int{1, 2, 3}), "successful collect and sort")
	require(slices.Sorted(func(func(int) bool) {}) == nil, "empty producer yielded nonnil")
	source := []int{3, 1, 2}
	sorted := slices.Sorted(slices.Values(source))
	require(slices.Equal(sorted, []int{1, 2, 3}) && slices.Equal(source, []int{3, 1, 2}),
		"Sorted failed to collect before sorting into separate storage")
	sorted[0] = -1
	require(source[0] == 3, "Sorted result aliases source")

	cases := []struct{ row, expected string }{
		{"0 1 1 1 1 0 3 2 1", "ok/slice:1,2,3/0/|compare:1/slice:2,3,1/0/2,3;1,3|ok/slice:1,2,3/3/|compare:1/nil/3/2,3;1,3"},
		{"0 1 2 1 -1 0 3 2 1", "ok/slice:1,2,3/0/|ok/slice:1,2,3/0/2,3;1,3;1,2|producer:7/nil/3/|producer:7/nil/3/"},
		{"0 0 1 0 -1 0", "ok/nil/0/|ok/nil/0/|ok/nil/0/|ok/nil/0/"},
		{"0 1 0 0 -1 0 2 1", "ok/slice:1,2/0/|nilCompare/slice:2,1/0/|nilSequence/nil/0/|nilSequence/nil/0/"},
	}
	var wire strings.Builder
	for _, c := range cases {
		fmt.Fprintln(&wire, strconv.IntSize, c.row)
	}
	for i, got := range check.LeanRows("godiff-slices-public-sort", len(cases), strings.NewReader(wire.String())) {
		require(got == cases[i].expected, fmt.Sprintf("row %d: Lean %q, want %q", i, got, cases[i].expected))
	}
	fmt.Printf("Public slice sort: %d hand-computed Lean rows and pinned native mutation controls\n", len(cases))
}
