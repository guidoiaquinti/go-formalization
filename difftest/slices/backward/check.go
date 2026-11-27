package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"slices"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type namedSlice []int

type input struct {
	values           []int
	kind, stop, edit int
}

type entry struct{ index, value int }
type callbackPanic int

func row(status string, cells []int, trace []entry) string {
	var out strings.Builder
	out.WriteString(status)
	for _, v := range cells {
		fmt.Fprintf(&out, " %d", v)
	}
	out.WriteString(" trace")
	for _, pair := range trace {
		fmt.Fprintf(&out, " %d %d", pair.index, pair.value)
	}
	return out.String()
}

func actual(c input) (result string) {
	cells := namedSlice(slices.Clone(c.values))
	seq := slices.Backward(cells)
	var trace []entry
	defer func() {
		switch p := recover().(type) {
		case nil:
		case callbackPanic:
			result = row(fmt.Sprint("panic ", p), cells, trace)
		case runtime.Error:
			if c.kind != 3 {
				panic(p)
			}
			result = row("nil-panic", cells, trace)
		default:
			panic(p)
		}
	}()
	if c.kind == 3 {
		seq(nil)
	} else {
		seq(func(i, value int) bool {
			trace = append(trace, entry{i, value})
			if i > 0 {
				cells[i-1] += c.edit
			}
			if i == c.stop && c.kind == 2 {
				panic(callbackPanic(i + 1))
			}
			return i != c.stop
		})
	}
	return row("ok", cells, trace)
}

func expected(c input) string {
	if c.kind == 3 && len(c.values) != 0 {
		return row("nil-panic", c.values, nil)
	}
	first := len(c.values) - 1
	last := 0
	status := "ok"
	if c.kind != 3 && 0 <= c.stop && c.stop <= first {
		last = c.stop
		if c.kind == 2 {
			status = fmt.Sprint("panic ", last+1)
		}
	}
	final := slices.Clone(c.values)
	for i := max(0, last-1); i < first; i++ {
		final[i] += c.edit
	}
	var trace []entry
	for i := first; i >= last; i-- {
		trace = append(trace, entry{i, final[i]})
	}
	return row(status, final, trace)
}

func main() {
	edges := []input{
		{nil, 3, -1, 0}, {[]int{}, 0, -1, 4}, {[]int{1}, 3, -1, 0},
		{[]int{1, 2, 3}, 0, -1, 10}, {[]int{1, 2, 3}, 1, 0, 10}, {[]int{1, 2, 3}, 1, 1, 10},
		{[]int{1, 2, 3}, 2, 0, 10}, {[]int{1, 2, 3}, 2, 1, -10}, {[]int{1, 2, 3}, 2, 2, 5},
		{[]int{4, 4, 4, 4}, 0, -1, 0}, {nil, 2, 0, 0},
	}
	check.CompareSamples("godiff-slices-backward", edges, func(r *rand.Rand, i int) input {
		values := make([]int, r.Intn(65))
		for j := range values {
			values[j] = r.Intn(201) - 100
		}
		kind := i % 4
		stop := -1
		if kind == 1 || kind == 2 {
			stop = r.Intn(len(values) + 2)
		}
		return input{values, kind, stop, r.Intn(31) - 15}
	}, func(c input) string {
		var out strings.Builder
		fmt.Fprintf(&out, "%d %d %d", c.kind, c.stop, c.edit)
		for _, v := range c.values {
			fmt.Fprintf(&out, " %d", v)
		}
		return out.String()
	}, actual, expected)
}
