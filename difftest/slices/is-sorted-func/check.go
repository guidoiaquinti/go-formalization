package main

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
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
	trace []int
}

var panicMarker = new(int)

func (r *recorder) compare(i, j int) int {
	if j != i-1 || i <= 0 || i >= r.length {
		panic("unexpected comparison operands")
	}
	r.trace = append(r.trace, i)
	switch r.responses[i-1] {
	case 0:
		return math.MaxInt
	case 1:
		return math.MinInt
	case 3:
		return 0
	default:
		panic(panicMarker)
	}
}

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
		row = observation(status, r.trace)
	}()
	values := make([]int, c.length)
	for i := range values {
		values[i] = i
	}
	status = strconv.FormatBool(slices.IsSortedFunc(values, r.compare))
	return
}

func main() {
	edges := []input{{0, nil}, {1, nil}, {4, []int{0, 3, 0}}, {4, []int{2, 1, 0}}, {4, []int{1, 2, 0}}, {2, []int{1}}}
	check.CompareSamples("godiff-slices-is-sorted-func", edges, func(r *rand.Rand, i int) input {
		n := r.Intn(33)
		responses := make([]int, max(0, n-1))
		for j := range responses {
			if i%4 != 0 && r.Intn(5) == 0 {
				responses[j] = 1 + r.Intn(3)
			}
		}
		return input{n, responses}
	}, func(c input) string {
		return observation(strconv.Itoa(c.length), c.responses)
	}, actual, func(c input) string {
		stop, status := 1, "true"
		for i, v := range slices.Backward(c.responses) {
			if v == 1 || v == 2 {
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
