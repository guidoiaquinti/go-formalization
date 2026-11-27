package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Words                           []int64
	A, B, Limit, Mode, FailAt, Mask int
	Want                            string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	file, err := os.CreateTemp("", "golean-slices-pdqsort-")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	encoder := json.NewEncoder(file)
	compare := func(label string, cases []fixture) error {
		var rows bytes.Buffer
		for _, c := range cases {
			fmt.Fprint(&rows, strconv.IntSize, " ", c.A, " ", c.B, " ", c.Limit, " ", c.Mode, " ", c.FailAt, " ", c.Mask)
			for _, word := range c.Words {
				fmt.Fprint(&rows, " ", word)
			}
			fmt.Fprintln(&rows)
		}
		for i, got := range check.LeanRows("godiff-slices-pdqsort", len(cases), &rows) {
			c := cases[i]
			c.Want = got
			if err := encoder.Encode(c); err != nil {
				return err
			}
		}
		fmt.Printf("pdqsorts ordered and comparator %s: %d Lean observations produced.\n", label, len(cases))
		return nil
	}
	var edges []fixture
	for _, n := range []int{0, 1, 2, 11, 12, 13, 14, 24, 25, 48, 49, 50, 51, 64, 101, 129} {
		for _, a := range []int{0, 3} {
			for pattern := range 6 {
				data := layout(a, n, pattern)
				for _, limit := range []int{math.MinInt, -1, 0, 1, 3, math.MaxInt} {
					for mode := range 3 {
						for _, stop := range []int{-1, 0, 1, 2, 3, 10, 11, 12, 13, 30, 150} {
							edges = append(edges, fixture{Words: data, A: a, B: a + n, Limit: limit,
								Mode: mode, FailAt: stop, Mask: 2730})
						}
					}
				}
			}
		}
	}
	if err := compare("boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			cases[i] = sample(rng, i)
		}
		if err := compare(fmt.Sprintf("seed=%d", seed), cases); err != nil {
			return err
		}
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_SLICES_PDQSORT_CASES", file.Name()); err != nil {
		return err
	}
	const data = "difftest/slices/pdqsort/testdata/"
	return check.TestInPackage("slices", "TestGoLeanPdqsort", data+"bridge_test.go", data+"pdqsort_test.go")
}

func layout(a, n, pattern int) []int64 {
	data := make([]int64, a+n+2)
	for i := range data {
		switch pattern {
		case 0:
			data[i] = int64(i)
		case 1:
			data[i] = int64(-i)
		case 2:
			data[i] = int64(i % 4)
		case 3:
			data[i] = 12
		case 4:
			data[i] = int64(min(i, a+n-i))
		default:
			data[i] = int64(i + (i%11)*5)
		}
	}
	for i := range a {
		data[i] = math.MinInt64
	}
	return data
}

func sample(rng *rand.Rand, i int) fixture {
	a := rng.Intn(6)
	n := rng.Intn(65)
	if i%4 == 0 {
		n = rng.Intn(20)
	}
	data := layout(a, n, rng.Intn(6))
	for j := range n {
		if i%3 != 0 && rng.Intn(5) == 0 {
			data[a+j] = int64(rng.Uint64())
		}
	}
	return fixture{Words: data, A: a, B: a + n, Limit: rng.Intn(6) - 1,
		Mode: rng.Intn(3), FailAt: rng.Intn(100) - 18, Mask: rng.Intn(4096)}
}
