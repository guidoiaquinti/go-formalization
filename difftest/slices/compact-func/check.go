package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"slices"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type namedSlice []int

type testCase struct {
	values        []int
	mode, fault   int
	offset, spare int
}

type pair struct{ current, previous int }

func format(header string, backing []int, calls []pair) string {
	var row strings.Builder
	row.WriteString(header)
	for _, value := range backing {
		fmt.Fprintf(&row, " %d", value)
	}
	row.WriteString(" trace")
	for _, call := range calls {
		fmt.Fprintf(&row, " %d %d", call.current, call.previous)
	}
	return row.String()
}

func actual(c testCase) (row string) {
	var backing, view namedSlice
	if c.values != nil {
		backing = make(namedSlice, c.offset+len(c.values)+c.spare)
		for i := range backing {
			backing[i] = i + 101
		}
		view = backing[c.offset : c.offset+len(c.values)]
		copy(view, c.values)
	}
	checkGuards := func() {
		for i, value := range backing {
			if (i < c.offset || i >= c.offset+len(view)) && value != i+101 {
				log.Fatalf("CompactFunc wrote outside the original view at %d", i)
			}
		}
	}

	var calls []pair
	type panicValue struct{ token int }
	payload := &panicValue{c.fault + 1}
	defer func() {
		p := recover()
		if p == nil {
			return
		}
		checkGuards()
		if p == payload {
			row = format(fmt.Sprint("panic ", payload.token), view, calls)
		} else if _, ok := p.(runtime.Error); ok && c.mode == 3 {
			row = format("panic 0", view, calls)
		} else {
			panic(p)
		}
	}()

	var compare func(int, int) bool
	if c.mode != 3 {
		compare = func(current, previous int) bool {
			i := len(calls)
			calls = append(calls, pair{current, previous})
			if i == c.fault {
				panic(payload)
			}
			switch c.mode {
			case 0:
				return current == previous
			case 1:
				return current == previous+1
			default:
				return i%3 == 0
			}
		}
	}
	result := slices.CompactFunc(view, compare)
	checkGuards()
	if (result == nil) != (view == nil) || cap(result) != cap(view) || len(result) > len(view) {
		log.Fatal("CompactFunc changed nilness or capacity, or expanded the slice")
	}
	if cap(view) > 0 && &result[:cap(result)][0] != &view[:cap(view)][0] {
		log.Fatal("CompactFunc changed the base address")
	}
	if !slices.Equal(result, view[:len(result)]) {
		log.Fatal("CompactFunc did not return the compacted prefix")
	}
	return format(fmt.Sprint("ok ", len(result)), view, calls)
}

func expected(c testCase) string {
	backing := slices.Clone(c.values)
	if len(backing) < 2 {
		return format(fmt.Sprint("ok ", len(backing)), backing, nil)
	}
	if c.mode == 3 {
		return format("panic 0", backing, nil)
	}
	var calls []pair
	kept := []int{c.values[0]}
	for i, current := range c.values[1:] {
		previous := c.values[i]
		calls = append(calls, pair{current, previous})
		if i == c.fault {
			copy(backing, kept)
			return format(fmt.Sprint("panic ", i+1), backing, calls)
		}
		drop := (c.mode == 0 && current == previous) ||
			(c.mode == 1 && current-previous == 1) || (c.mode == 2 && i%3 == 0)
		if !drop {
			kept = append(kept, current)
		}
	}
	copy(backing, kept)
	clear(backing[len(kept):])
	return format(fmt.Sprint("ok ", len(kept)), backing, calls)
}

func main() {
	edges := []testCase{
		{nil, 3, -1, 0, 0},
		{[]int{}, 3, -1, 2, 4},
		{[]int{7}, 3, -1, 1, 4},
		{[]int{7, 8}, 3, -1, 1, 4},
		{[]int{1, 1, 2, 2, 3}, 0, -1, 1, 4},
		{[]int{1, 2, 3, 2, 3, 9}, 1, -1, 1, 1},
		{[]int{1, 2, 3, 4, 5, 6}, 2, -1, 1, 3},
		{[]int{1, 2, 3, 4}, 0, 0, 1, 4},
		{[]int{1, 2, 3, 4}, 0, 2, 0, 4},
		{[]int{1, 1, 4, 9, 9}, 0, 2, 2, 4},
		{[]int{1, 1, 4, 9, 9}, 0, 3, 2, 4},
		{[]int{1, 2, 3, 4, 5}, 2, 3, 2, 4},
	}
	check.CompareSamples("godiff-slices-compact-func", edges, func(r *rand.Rand, i int) testCase {
		values := make([]int, r.Intn(33))
		for j := range values {
			values[j] = r.Intn(11) - 5
		}
		fault := -1
		if i%3 == 0 && len(values) > 1 {
			fault = r.Intn(len(values) - 1)
		}
		mode := r.Intn(3)
		if i%19 == 0 {
			mode = 3
		}
		return testCase{values, mode, fault, r.Intn(5), r.Intn(9)}
	}, func(c testCase) string {
		var header strings.Builder
		fmt.Fprintf(&header, "%d %d", c.mode, c.fault)
		for _, v := range c.values {
			fmt.Fprintf(&header, " %d", v)
		}
		return header.String()
	}, actual, expected)
}
