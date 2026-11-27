package main

import (
	"bytes"
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func runLean(name string, cases [][2]string) {
	var rows bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&rows, c[0])
	}
	for i, row := range check.LeanRows(name, len(cases), &rows) {
		if row != cases[i][1] {
			panic(fmt.Sprintf("%s row %d got %q want %q", name, i, row, cases[i][1]))
		}
	}
}

type callbackFault int

func recovered(action func()) (value any) {
	defer func() { value = recover() }()
	action()
	return
}

func chunkControl() {
	data := []int{5, 6, 7, 8, 9, 90, 91}
	seq := slices.Chunk(data[:5], 2)
	seq(func(first []int) bool {
		if cap(first) != 2 || !slices.Equal(first, []int{5, 6}) || &first[0] != &data[0] {
			panic("unexpected first chunk geometry")
		}
		extended := append(first, 99)
		if &first[0] == &extended[0] || data[2] != 7 {
			panic("appending modified next chunk")
		}
		first[0] = 44
		return false
	})
	var trace []string
	seq(func(part []int) bool {
		trace = append(trace, fmt.Sprintf("%d:%v", cap(part), part))
		return true
	})
	if strings.Join(trace, ";") != "2:[44 6];2:[7 8];1:[9]" ||
		!slices.Equal(data, []int{44, 6, 7, 8, 9, 90, 91}) {
		panic("unexpected restarted chunk sequence")
	}
	if p := recovered(func() { _ = slices.Chunk([]int(nil), 0) }); p != "cannot be less than 1" {
		panic("bad size did not panic at construction")
	}
	if len(slices.Collect(slices.Chunk([]int{}, 2))) != 0 {
		panic("unexpected empty yield")
	}
	runLean("godiff-slices-chunk", [][2]string{
		{"5 2 1 99 1 44 5 6 7 8 9 90 91", "2:5,6/44,6,7,8,9,90,91|2:44,6;2:7,8;1:9/44,6,7,8,9,90,91"},
		{"0 0 0 0 0 0", "cannot be less than 1"},
		{"0 2 2 2 0 0 90", "/90|/90"},
	})
	// A continuation-style iterator would resume at the second chunk and
	// miss the mutation in the first backing cell.
	continuation := fmt.Sprintf("2:%v;1:%v", data[2:4], data[4:5])
	if strings.Join(trace, ";") == continuation || data[0] == 5 {
		panic("chunk mutation survived")
	}
}

func stableControl() {
	before := []int64{513, 256, 512}
	order := func(x, y int64) int { return cmp.Compare(x>>8, y>>8) }
	arr := slices.Clone(before)
	slices.SortStableFunc(arr, order)
	want := []int64{256, 513, 512}
	if !slices.Equal(arr, want) {
		panic("stable order changed")
	}
	beforeCount := 0
	comparisons := 0
	seq := func(yield func(int64) bool) {
		for _, v := range before {
			if !yield(v) {
				panic("collection stopped prematurely")
			}
			beforeCount++
		}
	}
	from := slices.SortedStableFunc(seq, func(x, y int64) int {
		if beforeCount != len(before) {
			panic("early comparison")
		}
		comparisons++
		return order(x, y)
	})
	if comparisons != 2 || !slices.Equal(from, want) {
		panic("collection sort mismatch")
	}
	empty := slices.SortedStableFunc(func(func(int64) bool) {}, (func(int64, int64) int)(nil))
	if empty != nil {
		panic("empty producer returned a nonnil slice")
	}
	arr = slices.Clone(before)
	calls := 0
	p := recovered(func() {
		slices.SortStableFunc(arr, func(x, y int64) int {
			if calls == 1 {
				panic(callbackFault(1))
			}
			calls++
			return order(x, y)
		})
	})
	if p != callbackFault(1) || !slices.Equal(arr, want) {
		panic("comparator panic rolled back completed write")
	}
	calls = 0
	p = recovered(func() {
		slices.SortedStableFunc(func(yield func(int64) bool) {
			yield(513)
			yield(256)
			panic(callbackFault(7))
		}, func(x, y int64) int { calls++; return order(x, y) })
	})
	if p != callbackFault(7) || calls != 0 {
		panic("producer panic did not precede sorting")
	}
	runLean("godiff-slices-public-stable", [][2]string{
		{"1 1 1 -1 0 513 256 512", "ok/slice:256,513,512/0/256,513;512,513|ok/slice:256,513,512/3/256,513;512,513"},
		{"1 1 1 1 0 513 256 512", "compare:1/slice:256,513,512/0/256,513;512,513|compare:1/nil/3/256,513;512,513"},
		{"1 2 1 -1 0 513 256", "ok/slice:256,513/0/256,513|producer:7/nil/2/"},
		{"1 1 0 -1 0", "ok/slice:/0/|ok/nil/0/"},
	})
	tiesExchange := slices.Clone(before)
	slices.SortStableFunc(tiesExchange, func(x, y int64) int {
		if order(x, y) <= 0 {
			return -1
		}
		return 1
	})
	if slices.Equal(want, tiesExchange) || slices.Equal(want, before) {
		panic("tie exchange or rollback mutation survived")
	}
}

func main() {
	chunkControl()
	stableControl()
	fmt.Println("seven fixed Lean rows and matching pinned Go controls passed")
}
