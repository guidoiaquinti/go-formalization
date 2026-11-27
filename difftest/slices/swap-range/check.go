package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Backing                   []int
	A, B, Count, Length, Mode int
	Want                      string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	file, err := os.CreateTemp("", "golean-slices-swap-range-")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	writer := json.NewEncoder(file)
	compare := func(label string, cases []fixture) error {
		var rows bytes.Buffer
		for _, c := range cases {
			tag := 1
			if c.Backing == nil {
				tag = 0
			}
			fmt.Fprint(&rows, tag, " ", c.A, " ", c.B, " ", c.Count, " ", c.Length, " ", c.Mode)
			for _, value := range c.Backing {
				fmt.Fprint(&rows, " ", value)
			}
			fmt.Fprintln(&rows)
		}
		for i, output := range check.LeanRows("godiff-slices-swap-range", len(cases), &rows) {
			c := cases[i]
			c.Want = expected(c)
			if output != c.Want {
				return fmt.Errorf("offsets %d/%d count %d backing %v: Lean=%q want %q", c.A, c.B, c.Count, c.Backing, output, c.Want)
			}
			if err := writer.Encode(c); err != nil {
				return err
			}
		}
		fmt.Printf("swapRange ordered and comparator %s: %d Lean rows agree.\n", label, len(cases))
		return nil
	}
	edges := []fixture{{}, {Count: 1}, {Backing: []int{}, Mode: 1},
		{Backing: []int{1, 2, 99}, Length: 2, A: math.MinInt, B: math.MaxInt, Count: math.MinInt},
		{Backing: []int{1, 2, 99}, Length: 2, A: math.MaxInt, Count: 1, Mode: 1},
		{Backing: []int{1, 2, 99}, Length: 2, B: math.MinInt, Count: 1},
		{Backing: []int{1, 2, 99}, Length: 2, A: 0, B: 0, Count: math.MaxInt},
	}
	for n := 0; n <= 8; n++ {
		backing := make([]int, n+3)
		for i := range backing {
			backing[i] = 10*i - n
		}
		for a := -1; a <= n+1; a++ {
			for b := -1; b <= n+1; b++ {
				for count := -1; count <= n+2; count++ {
					edges = append(edges, fixture{Backing: backing, A: a, B: b, Count: count, Length: n, Mode: (a + b + count) & 1})
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
	if err := os.Setenv("GOLEAN_SLICE_SWAP_RANGE_CASES", file.Name()); err != nil {
		return err
	}
	const data = "difftest/slices/swap-range/testdata/"
	return check.TestInPackage("slices", "TestGoLeanSwapRange", data+"bridge_test.go", data+"swap_range_test.go")
}

func sample(rng *rand.Rand, i int) fixture {
	if i%17 == 0 {
		return fixture{Mode: i % 2}
	}
	n := rng.Intn(65)
	backing := make([]int, n+rng.Intn(17))
	for i := range backing {
		backing[i] = rng.Intn(2001) - 1000
	}
	a, b := rng.Intn(n+3)-1, rng.Intn(n+3)-1
	count := rng.Intn(n+5) - 2
	if i%5 == 0 {
		a, b, count = 0, min(1, n), rng.Intn(max(1, n))
	}
	if i%7 == 0 {
		count = math.MaxInt
	}
	return fixture{Backing: backing, A: a, B: b, Count: count, Length: n, Mode: i % 2}
}

func expected(c fixture) string {
	limit := 0
	if min(c.A, c.B) >= 0 && max(c.A, c.B) < c.Length {
		limit = min(max(0, c.Count), c.Length-max(c.A, c.B))
	}
	permutation := make([]int, c.Length)
	for i := range permutation {
		permutation[i] = i
	}
	for i := range limit {
		x, y := c.A+i, c.B+i
		permutation[x], permutation[y] = permutation[y], permutation[x]
	}
	state := "ok"
	if c.Count > limit {
		state = "panic"
	}
	if c.Backing == nil {
		return state + ":nil|" + state + ":nil"
	}
	values := slices.Clone(c.Backing)
	for i, source := range permutation {
		values[i] = c.Backing[source]
	}
	fields := make([]string, len(values))
	for i, value := range values {
		fields[i] = strconv.Itoa(value)
	}
	row := state + ":slice:" + strings.Join(fields, ",")
	return row + "|" + row
}
