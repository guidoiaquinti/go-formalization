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
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Kind, Shape, A, B int
	Replies           []int
	Want              string
}

func expected(c fixture) string {
	if c.Kind == 1 {
		if c.Shape != 0 {
			return "nil:"
		}
		j := c.Replies[0]
		result := strconv.FormatBool(j < c.A)
		if c.B == 2 {
			result = "callback"
		}
		return fmt.Sprintf("%s:%d,%d", result, j, c.A)
	}
	if c.Shape == 5 {
		return "invalid:"
	}
	if c.Shape >= 6 {
		return "unsupported:"
	}
	stop := -1
	for i, r := range c.Replies {
		if r != 0 {
			stop = i
		}
	}
	result := "true"
	if stop >= 0 {
		result = "false"
		if c.Replies[stop] == 2 {
			result = "callback"
		}
	}
	var calls []string
	for i := c.A - 1; i > max(stop+1, 0); i-- {
		calls = append(calls, fmt.Sprintf("%d,%d", i, i-1))
	}
	if stop >= 0 {
		calls = append(calls, fmt.Sprintf("%d,%d", stop+1, stop))
	}
	return result + ":" + strings.Join(calls, ";")
}

func main() {
	f, err := os.CreateTemp("", "golean-sort-comparison-")
	must(err)
	defer os.Remove(f.Name())
	buffer := bufio.NewWriter(f)
	rows := json.NewEncoder(buffer)
	for kind := range 2 {
		compare(fmt.Sprintf("kind=%d fixed", kind), boundaries(kind), rows)
		for _, seed := range []int64{17, 2718, 65537} {
			rng := rand.New(rand.NewSource(seed))
			for block := range 4 {
				cases := make([]fixture, 5000)
				for i := range cases {
					if kind == 1 {
						cases[i] = fixture{kind, rng.Intn(2), int(rng.Uint64()), rng.Intn(3), []int{int(rng.Uint64())}, ""}
					} else {
						n := rng.Intn(65)
						replies := make([]int, max(n-1, 0))
						for j := range replies {
							if i%3 != 0 && rng.Intn(13) == 0 {
								replies[j] = 1 + rng.Intn(2)
							}
						}
						cases[i] = fixture{kind, rng.Intn(8), n, 0, replies, ""}
					}
				}
				compare(fmt.Sprintf("kind=%d seed=%d block=%d", kind, seed, block), cases, rows)
			}
		}
	}
	must(buffer.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_SORT_COMPARISON_CASES", f.Name()))
	must(check.TestInPackage("sort", "TestGoLeanComparisonEntrypoints", "difftest/sort/comparison-entrypoints/testdata/comparison_test.go"))
}

func boundaries(kind int) []fixture {
	var cases []fixture
	if kind == 1 {
		for _, p := range [][2]int{{0, 0}, {1, 3}, {3, 1}, {math.MinInt, math.MaxInt}, {math.MaxInt, math.MinInt}} {
			for shape := range 2 {
				for _, reply := range []int{0, 2} {
					cases = append(cases, fixture{kind, shape, p[0], reply, []int{p[1]}, ""})
				}
			}
		}
		return cases
	}
	for shape := range 8 {
		for n, power := range []int{1, 1, 3, 9, 27, 81, 243} {
			for code := range power {
				replies := make([]int, max(n-1, 0))
				cur := code
				for i := range replies {
					replies[i] = cur % 3
					cur /= 3
				}
				cases = append(cases, fixture{kind, shape, n, 0, replies, ""})
			}
		}
	}
	return cases
}

func compare(label string, cases []fixture, rows *json.Encoder) {
	var buffer bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&buffer, "%d %d %d %d", c.Kind, c.Shape, c.A, c.B)
		for _, r := range c.Replies {
			fmt.Fprintf(&buffer, " %d", r)
		}
		fmt.Fprintln(&buffer)
	}
	for i, got := range check.LeanRows("godiff-sort-comparison-entrypoints", len(cases), &buffer) {
		c := cases[i]
		c.Want = expected(c)
		if got != c.Want {
			log.Fatalf("%+v: Lean %s", c, got)
		}
		must(rows.Encode(c))
	}
	fmt.Printf("sort comparison entrypoints %s: %d matched.\n", label, len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
