package main

import (
	"cmp"
	"fmt"
	"math"
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
	status = strconv.Itoa(slices.CompareFunc(xs, right, func(i int, j string) int {
		if i != calls || j != strconv.Itoa(calls) {
			panic("unexpected call order")
		}
		calls++
		if c.responses[i] == 2 {
			panic(panicMarker)
		}
		return c.responses[i]
	}))
	return
}

func main() {
	edges := []input{{nil, 0}, {[]int{}, 0}, {[]int{0}, 1}, {[]int{1}, 1}, {[]int{2}, 1}, {[]int{2}, 0}, {nil, 1}, {[]int{0, 1, 2}, 3}, {[]int{0, 2, 1}, 3}, {[]int{math.MinInt}, 2}, {[]int{math.MaxInt}, 1}, {[]int{0, 2}, 1}}
	check.CompareSamples("godiff-slices-compare-func", edges, func(r *rand.Rand, i int) input {
		xs := make([]int, r.Intn(65))
		for j := range xs {
			if i%4 != 0 && r.Intn(8) == 0 {
				xs[j] = []int{math.MinInt, -7, 2, 9, math.MaxInt}[r.Intn(5)]
			}
		}
		n := len(xs)
		if i%5 == 0 {
			n++
		} else if i%5 == 1 && n > 0 {
			n--
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
		overlap := min(len(xs), c.rightLength)
		stop := overlap
		for i, x := range slices.Backward(xs[:overlap]) {
			if x != 0 {
				stop = i
			}
		}
		if stop == overlap {
			return fmt.Sprintf("%d %d", cmp.Compare(len(xs), c.rightLength), overlap)
		}
		if xs[stop] == 2 {
			return fmt.Sprintf("panic %d", stop+1)
		}
		return fmt.Sprintf("%d %d", xs[stop], stop+1)
	})
}
