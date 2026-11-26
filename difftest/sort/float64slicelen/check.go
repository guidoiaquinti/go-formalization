package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type bounds struct {
	size, lo, hi int
	nilSlice     bool
}

func main() {
	edges := []bounds{{nilSlice: true}, {}, {size: 256}, {256, 0, 256, false}, {256, 255, 256, false}, {256, 256, 256, false}, {256, 127, 129, false}}
	check.CompareSamples("godiff-float64slicelen", edges, func(r *rand.Rand, _ int) bounds {
		size := r.Intn(257)
		lo := r.Intn(size + 1)
		return bounds{size, lo, lo + r.Intn(size-lo+1), false}
	}, func(b bounds) string {
		return fmt.Sprintf("%d %d %d", b.size, b.lo, b.hi)
	}, func(b bounds) string {
		var xs sort.Float64Slice
		if !b.nilSlice {
			xs = make(sort.Float64Slice, b.size)
			for i := range xs {
				xs[i] = math.Float64frombits(0x7ff0000000000000 + uint64(i))
			}
		}
		return strconv.Itoa(xs[b.lo:b.hi].Len())
	}, func(b bounds) string {
		return strconv.Itoa(b.hi - b.lo)
	})
}
