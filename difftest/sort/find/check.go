package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct{ n, cut, endMatch int }

func main() {
	edges := []input{{math.MinInt, 0, 0}, {-1, 0, 0}, {0, 0, 0}, {1, 0, 1}, {1, 0, 0}, {1, 1, 1},
		{5, 1, 4}, {5, 2, 2}, {math.MaxInt, math.MaxInt - 1, math.MaxInt}, {math.MaxInt, math.MaxInt, math.MaxInt}}
	check.CompareSamples("godiff-sort-find", edges, func(r *rand.Rand, i int) input {
		n := int(r.Uint64() & uint64(math.MaxInt))
		if i%3 == 0 {
			n = r.Intn(257)
		}
		cut := int(r.Uint64() % (uint64(n) + 1))
		endMatch := cut
		if i%2 == 0 {
			endMatch += int(r.Uint64() % (uint64(n-cut) + 1))
		}
		return input{n, cut, endMatch}
	}, func(c input) string {
		return fmt.Sprintf("%d %d %d", c.n, c.cut, c.endMatch)
	}, func(c input) string {
		i, found := sort.Find(c.n, func(i int) int {
			if i < 0 || i >= c.n {
				log.Fatalf("Find(%d): callback index %d out of bounds", c.n, i)
			}
			if i < c.cut {
				return math.MaxInt
			}
			if i < c.endMatch {
				return 0
			}
			return math.MinInt
		})
		return fmt.Sprintf("%d %t", i, found)
	}, func(c input) string {
		return fmt.Sprintf("%d %t", c.cut, c.cut < c.endMatch)
	})
}
