package main

import (
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
	Words                             []int64
	Kind, A, M, B, Mode, FailAt, Mask int
	Want                              string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	file, err := os.CreateTemp("", "golean-slices-stable-")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	encoder := json.NewEncoder(file)
	compare := func(label string, cases []fixture) error {
		var rows bytes.Buffer
		for _, c := range cases {
			fmt.Fprint(&rows, c.Kind, " ", c.A, " ", c.M, " ", c.B, " ", c.Mode, " ", c.FailAt, " ", c.Mask)
			for _, word := range c.Words {
				fmt.Fprint(&rows, " ", word)
			}
			fmt.Fprintln(&rows)
		}
		for i, got := range check.LeanRows("godiff-slices-stable", len(cases), &rows) {
			c := cases[i]
			c.Want = got
			if err := encoder.Encode(c); err != nil {
				return err
			}
		}
		fmt.Printf("stables ordered and comparator %s: %d Lean observations produced.\n", label, len(cases))
		return nil
	}
	var edges []fixture
	for _, n := range []int{0, 1, 2, 3, 4, 5, 10, 19, 20, 21, 39, 40, 41, 59, 60, 61, 79, 80, 81, 128} {
		for pattern := range 4 {
			for mode := range 3 {
				for _, stop := range []int{-1, 0, 1, 2, 3, 19, 20, 40, 80} {
					raw := words(n+3, pattern)
					edges = append(edges, fixture{Words: raw, Kind: 1, B: n, Mode: mode, FailAt: stop, Mask: 2730})
					for _, a := range []int{0, 2} {
						for _, m := range []int{a + 1, a + n/2, a + n - 1} {
							b := a + n
							if a < m && m < b {
								data := slices.Clone(raw)
								initialRuns(data, a, m, b)
								edges = append(edges, fixture{Words: data, A: a, M: m, B: b, Mode: mode, FailAt: stop, Mask: 2730})
							}
						}
					}
				}
			}
		}
	}
	edges = append(edges,
		fixture{Words: []int64{9, 4, 1, 3, 5, 8}, A: 1, M: 2, B: 5, Mode: 2, Mask: 2730, FailAt: -1},
		fixture{Words: []int64{512, 513, 514}, M: 2, B: 3, Mode: 1, FailAt: -1},
		fixture{Words: []int64{4, 1, 3, 2, 0}, Kind: 1, B: 4, Mode: 2, FailAt: -1},
		fixture{Words: []int64{3, 2, 1, 0}, Kind: 1, B: 3, Mode: 2, Mask: 4095, FailAt: 1})
	if err := compare("boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for kind := range 2 {
			cases := make([]fixture, 20_000)
			for i := range cases {
				cases[i] = sample(rng, kind)
			}
			if err := compare(fmt.Sprintf("%d seed=%d", kind, seed), cases); err != nil {
				return err
			}
		}
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_SLICES_STABLE_CASES", file.Name()); err != nil {
		return err
	}
	const data = "difftest/slices/stable/testdata/"
	return check.TestInPackage("slices", "TestGoLeanStable", data+"bridge_test.go", data+"stable_test.go")
}

func words(n, pattern int) []int64 {
	data := make([]int64, n)
	for i := range data {
		key := i / 2
		switch pattern {
		case 1:
			key = 20 - i/2
		case 2:
			key = i%5 - 10
		case 3:
			key = 2
		}
		data[i] = int64(key)*256 + int64(i)
	}
	return data
}

func initialRuns(data []int64, a, m, b int) {
	for _, run := range [][]int64{data[a:m], data[m:b]} {
		for i := range run {
			for j := i + 1; j < len(run); j++ {
				if run[j] < run[i] {
					run[i], run[j] = run[j], run[i]
				}
			}
		}
	}
}

func sample(rng *rand.Rand, kind int) fixture {
	n := rng.Intn(85)
	data := words(n+4, rng.Intn(4))
	for i := range data {
		if rng.Intn(4) == 0 {
			data[i] = int64(rng.Intn(7)-3)*256 + int64(i)
		}
	}
	c := fixture{Words: data, Kind: kind, B: n, Mode: rng.Intn(3),
		FailAt: rng.Intn(90) - 18, Mask: rng.Intn(4096)}
	if kind == 0 {
		c.A = rng.Intn(3)
		c.B += c.A + 2
		c.M = c.A + 1 + rng.Intn(c.B-c.A-1)
		initialRuns(data, c.A, c.M, c.B)
	}
	return c
}
