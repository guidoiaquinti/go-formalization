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
	Variant, A, M, B int
	Keys             []int
	Want             string
}

func main() {
	f, err := os.CreateTemp("", "golean-sort-sym-merge-")
	must(err)
	defer os.Remove(f.Name())
	buffer := bufio.NewWriter(f)
	rows := json.NewEncoder(buffer)
	compare("fixed", boundaries(), rows)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 0, 10000)
			for i := range 5000 {
				left, right := 1+rng.Intn(32), 1+rng.Intn(32)
				if i%113 == 0 {
					left, right = 1+rng.Intn(80), 1+rng.Intn(80)
				}
				if i%7 == 0 {
					left = 1
				}
				if i%11 == 0 {
					right = 1
				}
				a := rng.Intn(6)
				m, b := a+left, a+left+right
				keys := make([]int, b+rng.Intn(6))
				for j := range keys {
					keys[j] = rng.Intn(17) - 8
					if i%5 == 0 {
						keys[j] = int(rng.Uint64())
					}
					if i%13 == 0 {
						keys[j] = 1
					}
				}
				slices.Sort(keys[a:m])
				slices.Sort(keys[m:b])
				cases = add(cases, a, m, b, keys)
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases, rows)
		}
	}
	must(buffer.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_SORT_SYM_MERGE_CASES", f.Name()))
	must(check.TestInPackage("sort", "TestGoLeanSymMerge", "difftest/sort/sym-merge/testdata/sym_merge_test.go"))
}
func add(cases []fixture, a, m, b int, keys []int) []fixture {
	return append(cases, fixture{0, a, m, b, keys, ""}, fixture{1, a, m, b, keys, ""})
}
func boundaries() []fixture {
	var cases []fixture
	for width, power := range []int{1, 3, 9, 27, 81, 243, 729} {
		for code := range power {
			keys := make([]int, width)
			cursor := code
			for i := range keys {
				keys[i] = cursor%3 - 1
				cursor /= 3
			}
			for a := 0; a < width; a++ {
				for m := a + 1; m < width; m++ {
					for b := m + 1; b <= width; b++ {
						if slices.IsSorted(keys[a:m]) && slices.IsSorted(keys[m:b]) {
							cases = add(cases, a, m, b, keys)
						}
					}
				}
			}
		}
	}
	return cases
}
func compare(label string, cases []fixture, rows *json.Encoder) {
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d %d %d %d", c.Variant, c.A, c.M, c.B)
		for _, key := range c.Keys {
			fmt.Fprintf(&wire, " %d", key)
		}
		fmt.Fprintln(&wire)
	}
	for i, got := range check.LeanRows("godiff-sort-sym-merge", len(cases), &wire) {
		c := cases[i]
		c.Want = got
		must(rows.Encode(c))
	}
	fmt.Printf("sort symMerge %s: %d Lean rows.\n", label, len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
