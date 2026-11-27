package main

import (
	"fmt"
	"iter"
	"log"
	"reflect"
	"runtime"
	"slices"
)

type namedSlice []int
type entry struct{ i, v int }
type callbackFault string

type observation struct {
	trace []entry
	cells namedSlice
	panic string
}

func observe(makeSeq func(namedSlice) iter.Seq2[int, int], initial namedSlice, action string) (out observation) {
	out.cells = slices.Clone(initial)
	seq := makeSeq(out.cells)
	defer func() {
		switch fault := recover().(type) {
		case nil:
		case runtime.Error:
			if action != "nil" {
				panic(fault)
			}
			out.panic = "nil callback"
		case callbackFault:
			out.panic = string(fault)
		default:
			panic(fault)
		}
	}()
	if action == "nil" {
		seq(nil)
		return
	}
	seq(func(i, v int) bool {
		out.trace = append(out.trace, entry{i, v})
		if i > 0 {
			out.cells[i-1] += 10
		}
		if i == 1 && action == "panic" {
			panic(callbackFault("callback payload"))
		}
		return i != len(out.cells)-1 || action != "stop"
	})
	return
}

func main() {
	fixtures := []struct {
		input  namedSlice
		action string
		want   observation
	}{
		{nil, "nil", observation{}},
		{namedSlice{1}, "nil", observation{cells: namedSlice{1}, panic: "nil callback"}},
		{namedSlice{1, 2, 3}, "full", observation{[]entry{{2, 3}, {1, 12}, {0, 11}}, namedSlice{11, 12, 3}, ""}},
		{namedSlice{1, 2, 3}, "stop", observation{[]entry{{2, 3}}, namedSlice{1, 12, 3}, ""}},
		{namedSlice{1, 2, 3}, "panic", observation{[]entry{{2, 3}, {1, 12}}, namedSlice{11, 12, 3}, "callback payload"}},
	}
	snapshot := func(s namedSlice) iter.Seq2[int, int] { return slices.Backward(slices.Clone(s)) }
	ignoreStop := func(s namedSlice) iter.Seq2[int, int] {
		return func(yield func(int, int) bool) {
			for i := len(s) - 1; i >= 0; i-- {
				yield(i, s[i])
			}
		}
	}
	swallowPanic := func(s namedSlice) iter.Seq2[int, int] {
		return func(yield func(int, int) bool) {
			defer func() { _ = recover() }()
			slices.Backward(s)(yield)
		}
	}
	for _, check := range []struct {
		name string
		run  func(namedSlice) iter.Seq2[int, int]
		want int
	}{
		{"Go", slices.Backward[namedSlice, int], 0},
		{"Forward", slices.All[namedSlice, int], 3},
		{"Snapshot", snapshot, 2},
		{"Ignore stop", ignoreStop, 1},
		{"Swallow panic", swallowPanic, 2},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := observe(check.run, f.input, f.action); !reflect.DeepEqual(got, f.want) {
				bad++
				if check.want == 0 {
					log.Fatalf("%s on %v: got %+v, want %+v", f.action, f.input, got, f.want)
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", check.name, bad, len(fixtures))
		if bad != check.want {
			log.Fatalf("%s: wanted %d disagreements", check.name, check.want)
		}
	}
}
