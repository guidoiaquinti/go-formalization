package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Variant, A, B int
	Values        []int
	CrossesStart  bool
	Want          string
}

func main() {
	f, err := os.CreateTemp("", "golean-sort-partial-")
	must(err)
	defer os.Remove(f.Name())
	w := bufio.NewWriter(f)
	rows := json.NewEncoder(w)
	compare("boundaries", boundaries(), rows)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 0, 10000)
			for i := range 5000 {
				a, size := rng.Intn(8), rng.Intn(97)
				if i%3 == 0 {
					size = 49 + rng.Intn(3)
				}
				b := a + size
				values := make([]int, b+rng.Intn(8))
				for j := range values {
					values[j] = rng.Intn(31) - 15
				}
				if i%2 == 0 {
					slices.Sort(values[a:b])
					for range rng.Intn(9) {
						if size >= 2 {
							j := a + rng.Intn(size-1)
							values[j], values[j+1] = values[j+1], values[j]
						}
					}
				}
				cases = add(cases, a, b, values, false)
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases, rows)
		}
	}
	must(w.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_SORT_PARTIAL_CASES", f.Name()))
	must(check.TestInPackage("sort", "TestGoLeanPartialInsertion", "difftest/sort/partial-insertion-sort/testdata/partial_insertion_test.go"))
}

func add(cases []fixture, a, b int, values []int, cross bool) []fixture {
	return append(cases, fixture{0, a, b, values, cross, ""}, fixture{1, a, b, values, cross, ""})
}

func boundaries() []fixture {
	var cases []fixture
	powers := []int{1, 3, 9, 27, 81, 243}
	for n, power := range powers {
		for code := range power {
			values := make([]int, n)
			cur := code
			for i := range values {
				values[i] = cur%3 - 1
				cur /= 3
			}
			for a := 0; a <= n; a++ {
				for b := a; b <= n; b++ {
					cases = add(cases, a, b, values, false)
				}
			}
		}
	}
	for _, size := range []int{49, 50, 51, 80} {
		for a := range 5 {
			for inversions := range 8 {
				values := make([]int, a+size+2)
				for j := range values {
					values[j] = j
				}
				for j := range inversions {
					index := a + 2*j
					values[index], values[index+1] = values[index+1], values[index]
				}
				cases = add(cases, a, a+size, values, false)
			}
		}
	}
	cross := make([]int, 54)
	for i := range cross {
		cross[i] = i
	}
	copy(cross, []int{100, 101, 0, 1, -1})
	return add(cases, 2, 52, cross, true)
}

func compare(label string, cases []fixture, rows *json.Encoder) {
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d %d %d", c.Variant, c.A, c.B)
		for _, v := range c.Values {
			fmt.Fprintf(&wire, " %d", v)
		}
		fmt.Fprintln(&wire)
	}
	for i, got := range check.LeanRows("godiff-sort-partial-insertion", len(cases), &wire) {
		c := cases[i]
		c.Want = got
		must(rows.Encode(c))
	}
	fmt.Printf("sort partial insertion %s: %d Lean rows.\n", label, len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
