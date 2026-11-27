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

func actual(responses []int) (row string) {
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
	if responses != nil {
		xs = make([]int, len(responses))
	}
	for i := range xs {
		xs[i] = i
	}
	status = strconv.FormatBool(slices.ContainsFunc(xs, func(i int) bool {
		if i != calls {
			panic("unexpected call order")
		}
		calls++
		switch responses[i] {
		case 0:
			return false
		case 1:
			return true
		default:
			panic(panicMarker)
		}
	}))
	return
}

func main() {
	edges := [][]int{nil, {}, {0}, {1}, {2}, {0, 1, 2, 1}, {0, 2, 1}, {0, 0, 1, 1}}
	check.CompareSamples("godiff-slices-contains-func", edges, func(r *rand.Rand, i int) []int {
		xs := make([]int, r.Intn(65))
		for j := range xs {
			if i%4 != 0 && r.Intn(8) == 0 {
				xs[j] = 1 + r.Intn(2)
			}
		}
		return xs
	}, func(xs []int) string {
		var row strings.Builder
		row.WriteByte('0')
		for _, x := range xs {
			fmt.Fprintf(&row, " %d", x)
		}
		return row.String()
	}, actual, func(xs []int) string {
		stop := len(xs)
		for i, x := range slices.Backward(xs) {
			if x != 0 {
				stop = i
			}
		}
		if stop == len(xs) {
			return fmt.Sprintf("false %d", stop)
		}
		if xs[stop] == 2 {
			return fmt.Sprintf("panic %d", stop+1)
		}
		return fmt.Sprintf("true %d", stop+1)
	})
}
