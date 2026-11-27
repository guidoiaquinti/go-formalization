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

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Kind, Variant, A, B, C int
	Values                 []int
	Want                   string
}

func main() {
	file, err := os.CreateTemp("", "golean-sort-local-sorts-")
	must(err)
	defer os.Remove(file.Name())
	buffer := bufio.NewWriter(file)
	rows := json.NewEncoder(buffer)
	for kind := range 3 {
		compare(fmt.Sprintf("helper=%d fixed", kind), boundaries(kind), rows)
		for _, seed := range []int64{17, 2718, 65537} {
			rng := rand.New(rand.NewSource(seed))
			for block := range 4 {
				cases := make([]fixture, 0, 10000)
				for i := range 5000 {
					width := rng.Intn(25)
					if i%101 == 0 {
						width = 49 + rng.Intn(17)
					}
					values := make([]int, width)
					for j := range values {
						values[j] = rng.Intn(17) - 8
						if i%13 == 0 {
							values[j] = int(rng.Uint64())
						}
						if i%11 == 0 {
							values[j] = width - j
						}
					}
					a, b, c := rng.Intn(width+1), rng.Intn(width+1), 0
					a, b = min(a, b), max(a, b)
					if i%3 == 0 {
						a, b = 0, width
					}
					if kind == 1 {
						c = a
						b -= a
						a = rng.Intn(b + 1)
					}
					cases = add(cases, kind, a, b, c, values)
				}
				compare(fmt.Sprintf("helper=%d seed=%d block=%d", kind, seed, block), cases, rows)
			}
		}
	}
	must(buffer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_SORT_LOCAL_SORT_CASES", file.Name()))
	must(check.TestInPackage("sort", "TestGoLeanLocalSorts", "difftest/sort/local-sorts/testdata/local_sorts_test.go"))
}

func add(cases []fixture, kind, a, b, c int, values []int) []fixture {
	return append(cases, fixture{kind, 0, a, b, c, values, ""}, fixture{kind, 1, a, b, c, values, ""})
}

func boundaries(kind int) []fixture {
	var cases []fixture
	powers := []int{1, 3, 9, 27, 81, 243}
	for width := 0; width <= 5; width++ {
		for code := 0; code < powers[width]; code++ {
			values := make([]int, width)
			digits := code
			for i := range values {
				values[i] = digits%3 - 1
				digits /= 3
			}
			for a := 0; a <= width; a++ {
				for b := a; b <= width; b++ {
					if kind == 1 {
						size := b - a
						for root := 0; root <= size; root++ {
							cases = add(cases, kind, root, size, a, values)
						}
					} else {
						cases = add(cases, kind, a, b, 0, values)
					}
				}
			}
		}
	}
	values := []int{math.MaxInt, math.MinInt, 0, math.MinInt, math.MaxInt}
	if kind == 1 {
		return add(cases, kind, 0, 4, 1, values)
	}
	return add(cases, kind, 1, 5, 0, values)
}

func compare(label string, cases []fixture, rows *json.Encoder) {
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d %d %d %d %d", c.Kind, c.Variant, c.A, c.B, c.C)
		for _, v := range c.Values {
			fmt.Fprintf(&wire, " %d", v)
		}
		fmt.Fprintln(&wire)
	}
	for i, result := range check.LeanRows("godiff-sort-local-sorts", len(cases), &wire) {
		c := cases[i]
		c.Want = result
		must(rows.Encode(c))
	}
	fmt.Printf("sort local sorts %s: %d Lean rows.\n", label, len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
