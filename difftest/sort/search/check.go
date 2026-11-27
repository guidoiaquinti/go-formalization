package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct{ n, cut int }

func main() {
	edges := []input{{math.MinInt, 0}, {-1, 0}, {0, 0}, {1, 0}, {1, 1}, {math.MaxInt, 0}, {math.MaxInt, math.MaxInt - 1}, {math.MaxInt, math.MaxInt}}
	check.CompareSamples("godiff-sort-search", edges, func(r *rand.Rand, i int) input {
		n := int(r.Uint64() & uint64(math.MaxInt))
		if i%3 == 0 {
			n = r.Intn(257)
		}
		cut := int(r.Uint64() % (uint64(n) + 1))
		return input{n, cut}
	}, func(c input) string {
		return fmt.Sprintf("%d %d", c.n, c.cut)
	}, func(c input) string {
		return strconv.Itoa(sort.Search(c.n, func(i int) bool {
			if i < 0 || i >= c.n {
				log.Fatalf("Search(%d): callback index %d out of bounds", c.n, i)
			}
			return i >= c.cut
		}))
	}, func(c input) string {
		return strconv.Itoa(c.cut)
	})
}
