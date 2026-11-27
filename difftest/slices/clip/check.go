package main

import (
	"fmt"
	"log"
	"math/rand"
	"slices"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type namedSlice []int
type input struct {
	values        []int
	spare, offset int
}

func row(header string, values []int) string {
	var out strings.Builder
	out.WriteString(header)
	for _, v := range values {
		fmt.Fprintf(&out, " %d", v)
	}
	return out.String()
}

func main() {
	edges := []input{{}, {[]int{}, 0, 0}, {[]int{}, 1, 1}, {[]int{1}, 0, 0}, {[]int{1, 2}, 7, 3}}
	check.CompareSamples("godiff-slices-clip", edges, func(r *rand.Rand, _ int) input {
		values := make([]int, r.Intn(33))
		for i := range values {
			values[i] = int(r.Uint64())
		}
		return input{values, r.Intn(9), r.Intn(5)}
	}, func(c input) string {
		nilFlag := 0
		if c.values == nil {
			nilFlag = 1
		}
		return row(fmt.Sprintf("%d %d", nilFlag, c.spare), c.values)
	}, func(c input) string {
		var view namedSlice
		if c.values != nil {
			backing := make(namedSlice, c.offset+len(c.values)+c.spare)
			view = backing[c.offset : c.offset+len(c.values)]
			copy(view, c.values)
		}
		clipped := slices.Clip(view)
		if len(clipped) > 0 && (len(view) == 0 || &clipped[0] != &view[0]) {
			log.Fatal("Clip changed the backing address of a nonempty slice")
		}
		return row(fmt.Sprintf("%t %d", clipped == nil, cap(clipped)), clipped)
	}, func(c input) string {
		return row(fmt.Sprintf("%t %d", c.values == nil, len(c.values)), c.values)
	})
}
