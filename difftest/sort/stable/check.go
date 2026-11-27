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
	Kind, Tag, N int
	Keys         []int
	Want         string
}

func main() {
	f, err := os.CreateTemp("", "golean-sort-stable-")
	must(err)
	defer os.Remove(f.Name())
	buffer := bufio.NewWriter(f)
	rows := json.NewEncoder(buffer)
	for kind := range 4 {
		compare(fmt.Sprintf("kind=%d fixed", kind), boundaries(kind), rows)
		for _, seed := range []int64{17, 2718, 65537} {
			rng := rand.New(rand.NewSource(seed))
			for block := range 4 {
				cases := make([]fixture, 5000)
				for i := range cases {
					length := rng.Intn(57)
					if i%127 == 0 {
						length = 57 + rng.Intn(144)
					}
					n, tag := length, 0
					if kind < 3 && i%7 == 0 {
						n = rng.Intn(length+3) - 2
					}
					if kind == 2 && i%15 == 0 {
						tag = 1
					}
					if kind == 3 && i%13 == 0 {
						tag = 1 + rng.Intn(4)
					}
					values := keys(rng, length, i%5)
					if kind == 3 && tag == 4 {
						values = nil
					}
					cases[i] = fixture{kind, tag, n, values, ""}
				}
				compare(fmt.Sprintf("kind=%d seed=%d block=%d", kind, seed, block), cases, rows)
			}
		}
	}
	must(buffer.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_SORT_STABLE_CASES", f.Name()))
	must(check.TestInPackage("sort", "TestGoLeanStable", "difftest/sort/stable/testdata/stable_test.go"))
}

func keys(rng *rand.Rand, width, pattern int) []int {
	keys := make([]int, width)
	for i := range keys {
		switch pattern {
		case 0:
			keys[i] = rng.Intn(7) - 3
		case 1:
			keys[i] = i / 3
		case 2:
			keys[i] = (width - i) / 3
		case 3:
			keys[i] = 1
		case 4:
			keys[i] = int(rng.Uint64())
		}
	}
	return keys
}

func boundaries(kind int) []fixture {
	var cases []fixture
	rng := rand.New(rand.NewSource(1))
	for _, n := range []int{0, 1, 2, 19, 20, 21, 38, 39, 40, 41, 42, 59, 60, 61, 79, 80, 81, 159, 160, 161} {
		for pattern := range 5 {
			values := keys(rng, n, pattern)
			cases = append(cases, fixture{kind, 0, n, values, ""})
			if kind < 3 {
				cases = append(cases, fixture{kind, 0, n / 2, values, ""})
			}
			if kind == 2 {
				cases = append(cases, fixture{kind, 1, n, values, ""})
			}
			if kind == 3 {
				for tag := 1; tag <= 4; tag++ {
					if tag == 4 {
						cases = append(cases, fixture{kind, tag, 0, nil, ""})
					} else {
						cases = append(cases, fixture{kind, tag, n, values, ""})
					}
				}
			}
		}
	}
	if kind < 3 {
		cases = append(cases, fixture{kind, 0, math.MinInt, []int{4, 1}, ""})
	}
	if kind == 3 {
		cases = append(cases, fixture{kind, 4, 0, nil, ""})
	}
	return cases
}

func compare(label string, cases []fixture, rows *json.Encoder) {
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d %d %d", c.Kind, c.Tag, c.N)
		for _, v := range c.Keys {
			fmt.Fprintf(&wire, " %d", v)
		}
		fmt.Fprintln(&wire)
	}
	for i, got := range check.LeanRows("godiff-sort-stable", len(cases), &wire) {
		c := cases[i]
		c.Want = got
		must(rows.Encode(c))
	}
	fmt.Printf("sort stable %s: %d Lean rows.\n", label, len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
