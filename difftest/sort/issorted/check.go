package main

import (
	"fmt"
	"math/rand"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	length    int
	responses []int
}

type recorder struct {
	input
	trace    []int
	lenCalls int
}

var panicMarker = new(int)

func (r *recorder) Len() int {
	r.lenCalls++
	return r.length
}

func (r *recorder) Less(i, j int) bool {
	if j != i-1 || i <= 0 || i >= r.length {
		panic("unexpected comparison indices")
	}
	r.trace = append(r.trace, i)
	switch r.responses[i-1] {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(panicMarker)
	}
}

func (*recorder) Swap(i, j int) { panic("unexpected Swap") }

func observation(status string, trace []int) string {
	var row strings.Builder
	row.WriteString(status)
	for _, i := range trace {
		fmt.Fprintf(&row, " %d", i)
	}
	return row.String()
}

func actual(c input) (row string) {
	r := recorder{input: c}
	status := ""
	defer func() {
		if p := recover(); p != nil {
			if p == panicMarker {
				status = "panic"
			} else {
				status = fmt.Sprintf("unexpected panic %v", p)
			}
		}
		if r.lenCalls != 1 {
			status = fmt.Sprintf("Len called %d times", r.lenCalls)
		}
		row = observation(status, r.trace)
	}()
	status = strconv.FormatBool(sort.IsSorted(&r))
	return
}

func main() {
	edges := []input{{0, nil}, {1, nil}, {4, []int{0, 0, 0}}, {4, []int{2, 1, 0}}, {4, []int{1, 2, 0}}, {2, []int{1}}}
	check.CompareSamples("godiff-sort-issorted", edges, func(r *rand.Rand, i int) input {
		n := r.Intn(33)
		responses := make([]int, max(0, n-1))
		for j := range responses {
			if i%4 != 0 && r.Intn(5) == 0 {
				responses[j] = 1 + r.Intn(2)
			}
		}
		return input{n, responses}
	}, func(c input) string {
		return observation(strconv.Itoa(c.length), c.responses)
	}, actual, func(c input) string {
		stop, status := 1, "true"
		for i, v := range slices.Backward(c.responses) {
			if v != 0 {
				stop = i + 1
				status = "false"
				if v == 2 {
					status = "panic"
				}
				break
			}
		}
		var trace []int
		for i := c.length - 1; i >= stop; i-- {
			trace = append(trace, i)
		}
		return observation(status, trace)
	})
}
