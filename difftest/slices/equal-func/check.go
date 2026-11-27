package main

import (
	"fmt"
	"math/rand"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

var panicMarker = new(int)

type input struct {
	responses   []int
	rightLength int
}

func actual(c input) (row string) {
	calls := 0
	status := ""
	defer func() {
		if p := recover(); p != nil {
			if p == panicMarker {
				status = "panic"
			} else {
				status = fmt.Sprintf("unexpected panic %v", p)
			}
		}
		row = fmt.Sprintf("%s %d", status, calls)
	}()
	var xs []int
	if c.responses != nil {
		xs = make([]int, len(c.responses))
	}
	for i := range xs {
		xs[i] = i
	}
	right := make([]string, c.rightLength)
	for j := range right {
		right[j] = strconv.Itoa(j)
	}
	status = strconv.FormatBool(slices.EqualFunc(xs, right, func(i int, j string) bool {
		if i != calls || j != strconv.Itoa(calls) {
			panic("unexpected call order")
		}
		calls++
		switch c.responses[i] {
		case 0:
			return true
		case 1:
			return false
		default:
			panic(panicMarker)
		}
	}))
	return
}

func main() {
	edges := []input{{nil, 0}, {[]int{}, 0}, {[]int{0}, 1}, {[]int{1}, 1}, {[]int{2}, 1}, {[]int{2}, 0}, {nil, 1}, {[]int{0, 1, 2}, 3}, {[]int{0, 2, 1}, 3}}
	check.CompareSamples("godiff-slices-equal-func", edges, func(r *rand.Rand, i int) input {
		xs := make([]int, r.Intn(65))
		for j := range xs {
			if i%4 != 0 && r.Intn(8) == 0 {
				xs[j] = 1 + r.Intn(2)
			}
		}
		n := len(xs)
		if i%5 == 0 {
			n++
		}
		return input{xs, n}
	}, func(c input) string {
		var row strings.Builder
		fmt.Fprint(&row, c.rightLength)
		for _, x := range c.responses {
			fmt.Fprintf(&row, " %d", x)
		}
		return row.String()
	}, actual, func(c input) string {
		xs := c.responses
		if len(xs) != c.rightLength {
			return "false 0"
		}
		stop := len(xs)
		for i, x := range slices.Backward(xs) {
			if x != 0 {
				stop = i
			}
		}
		if stop == len(xs) {
			return fmt.Sprintf("true %d", stop)
		}
		if xs[stop] == 2 {
			return fmt.Sprintf("panic %d", stop+1)
		}
		return fmt.Sprintf("false %d", stop+1)
	})
}
