package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"slices"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Kind, Variant, A, B, Pivot int
	Values                     []int
	Want                       string
}

func main() {
	file, err := os.CreateTemp("", "golean-sort-partition-")
	must(err)
	defer os.Remove(file.Name())
	buffer := bufio.NewWriter(file)
	rows := json.NewEncoder(buffer)
	for kind := range 2 {
		compare(fmt.Sprintf("helper=%d boundaries", kind), boundaries(kind), rows)
		for _, seed := range []int64{17, 2718, 65537} {
			rng := rand.New(rand.NewSource(seed))
			for block := range 4 {
				cases := make([]fixture, 0, 10000)
				for i := range 5000 {
					size := 1 + rng.Intn(33)
					if i%101 == 0 {
						size += 64
					}
					a := rng.Intn(6)
					b := a + size
					pivot := a + rng.Intn(size)
					values := make([]int, b+rng.Intn(6))
					for j := range values {
						values[j] = rng.Intn(17) - 8
						if i%13 == 0 {
							values[j] = int(rng.Uint64())
						}
						if i%19 == 0 {
							values[j] = j
						}
						if i%29 == 0 {
							values[j] = len(values) - j
						}
					}
					if kind == 1 {
						boundary := values[pivot]
						for j := a; j < b; j++ {
							if values[j] < boundary || (i%2 == 0 && j%2 == 0) {
								values[j] = boundary
							}
						}
					}
					cases = add(cases, kind, a, b, pivot, values)
				}
				compare(fmt.Sprintf("helper=%d seed=%d block=%d", kind, seed, block), cases, rows)
			}
		}
	}
	must(buffer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_SORT_PARTITION_CASES", file.Name()))
	must(check.TestInPackage("sort", "TestGoLeanPartition", "difftest/sort/partition/testdata/partition_test.go"))
}

func add(cases []fixture, kind, a, b, pivot int, values []int) []fixture {
	return append(cases, fixture{kind, 0, a, b, pivot, values, ""}, fixture{kind, 1, a, b, pivot, values, ""})
}

func boundaries(kind int) []fixture {
	var cases []fixture
	powers := []int{1, 3, 9, 27, 81, 243}
	for width := 1; width <= 5; width++ {
		for code := 0; code < powers[width]; code++ {
			values := make([]int, width)
			digits := code
			for i := range values {
				values[i] = digits%3 - 1
				digits /= 3
			}
			for a := 0; a < width; a++ {
				for b := a + 1; b <= width; b++ {
					least := slices.Min(values[a:b])
					for pivot := a; pivot < b; pivot++ {
						if kind == 0 || values[pivot] == least {
							cases = add(cases, kind, a, b, pivot, values)
						}
					}
				}
			}
		}
	}
	values := []int{math.MaxInt, math.MinInt, 0, math.MinInt, math.MaxInt}
	return add(cases, kind, 1, 5, 3, values)
}

func compare(label string, cases []fixture, rows *json.Encoder) {
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d %d %d %d %d", c.Kind, c.Variant, c.A, c.B, c.Pivot)
		for _, v := range c.Values {
			fmt.Fprintf(&wire, " %d", v)
		}
		fmt.Fprintln(&wire)
	}
	for i, result := range check.LeanRows("godiff-sort-partition", len(cases), &wire) {
		c := cases[i]
		c.Want = result
		must(rows.Encode(c))
	}
	fmt.Printf("sort partition %s: %d Lean rows.\n", label, len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
