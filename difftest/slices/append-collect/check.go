package main

import (
	"fmt"
	"iter"
	"math/rand"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type namedSlice []int
type producerFault int

type input struct {
	cells, emitted      []int
	length, mode, panic int
}

func main() {
	edges := []input{
		{nil, nil, 0, 0, 0},
		{[]int{}, nil, 0, 0, 0},
		{[]int{1, 2, 3}, nil, 2, 0, 0},
		{nil, []int{4, 5}, 0, 0, 0},
		{[]int{8, 99, 98}, []int{4, 5}, 1, 0, 0},
		{[]int{8, 99, 98}, []int{4, 5, 6}, 1, 0, 0},
		{[]int{11, 90, 91}, []int{4, 5, 6}, 1, 1, -73},
		{[]int{9, 8}, []int{7}, 0, 1, 23},
		{[]int{9, 8}, nil, 0, 1, 0},
		{[]int{9, 8}, nil, 1, 2, 0},
		{nil, nil, 0, 2, 0},
		{[]int{9, 8}, []int{6}, 0, 0, 0},
	}
	check.CompareSamples("godiff-slices-append-collect", edges, sample, encode, actual, expected)
}

func sample(rng *rand.Rand, i int) input {
	var cells []int
	if i%7 != 0 {
		cells = make([]int, rng.Intn(33))
		for j := range cells {
			cells[j] = rng.Intn(201) - 100
		}
	}
	mode := 0
	if i%5 == 1 {
		mode = 1
	} else if i%5 == 2 {
		mode = 2
	}
	var emitted []int
	if mode != 2 {
		emitted = make([]int, rng.Intn(49))
		for j := range emitted {
			emitted[j] = rng.Intn(201) - 100
		}
	}
	return input{cells, emitted, rng.Intn(len(cells) + 1), mode, rng.Intn(1001) - 500}
}

func encode(c input) string {
	tag := 1
	if c.cells == nil {
		tag = 0
	}
	var row strings.Builder
	fmt.Fprintf(&row, "%d %d %d %d %d", tag, c.length, len(c.cells), c.mode, c.panic)
	for _, block := range [][]int{c.cells, c.emitted} {
		for _, v := range block {
			fmt.Fprint(&row, " ", v)
		}
	}
	return row.String()
}

func actual(c input) string {
	return invoke(c, false) + "|" + invoke(c, true)
}

func invoke(c input, collect bool) (observation string) {
	var cells []int
	if !collect {
		cells = slices.Clone(c.cells)
	}
	acked := 0
	defer func() {
		switch fault := recover().(type) {
		case nil:
		case producerFault:
			observation = wire(fmt.Sprintf("producer:%d", fault), acked, false, cells)
		case runtime.Error:
			if c.mode != 2 || fault.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(fault)
			}
			observation = wire("nil-seq", acked, false, cells)
		default:
			panic(fault)
		}
	}()
	var sequence iter.Seq[int]
	if c.mode != 2 {
		sequence = func(yield func(int) bool) {
			for _, value := range c.emitted {
				if !yield(value) {
					return
				}
				acked++
			}
			if c.mode == 1 {
				panic(producerFault(c.panic))
			}
		}
	}
	var output []int
	if collect {
		output = slices.Collect(sequence)
	} else {
		output = slices.AppendSeq(namedSlice(cells[:c.length:len(cells)]), sequence)
	}
	reused := len(cells) != 0 && cap(output) != 0 && &cells[0] == &output[:1][0]
	return wire(sliceStatus(output), acked, reused, cells)
}

func expected(c input) string {
	return expectedOne(c, false) + "|" + expectedOne(c, true)
}

func expectedOne(c input, collect bool) string {
	if collect {
		c.cells, c.length = nil, 0
	}
	cells := slices.Clone(c.cells)
	if c.mode == 2 {
		return wire("nil-seq", 0, false, cells)
	}
	writes := copy(cells[c.length:], c.emitted)
	if c.mode == 1 {
		return wire(fmt.Sprintf("producer:%d", c.panic), len(c.emitted), false, cells)
	}
	var output []int
	if c.cells != nil || len(c.emitted) != 0 {
		output = make([]int, c.length+len(c.emitted))
		copy(output, c.cells[:c.length])
		copy(output[c.length:], c.emitted)
	}
	reused := len(cells) != 0 && writes == len(c.emitted)
	return wire(sliceStatus(output), len(c.emitted), reused, cells)
}

func sliceStatus(output []int) string {
	if output == nil {
		return "nil"
	}
	return "slice:" + cellsWire(output)
}

func wire(status string, acknowledged int, reused bool, cells []int) string {
	return fmt.Sprintf("%s/%d/%t/%s", status, acknowledged, reused, cellsWire(cells))
}

func cellsWire(cells []int) string {
	fields := make([]string, len(cells))
	for i, value := range cells {
		fields[i] = strconv.Itoa(value)
	}
	return strings.Join(fields, ",")
}
