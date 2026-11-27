package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Backing                    []int
	A, Middle, B, Length, Mode int
	Want                       string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	file, err := os.CreateTemp("", "golean-slices-rotate-range-")
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
			fmt.Fprint(&rows, tag, " ", c.A, " ", c.Middle, " ", c.B, " ", c.Length, " ", c.Mode)
			for _, value := range c.Backing {
				fmt.Fprint(&rows, " ", value)
			}
			fmt.Fprintln(&rows)
		}
		for i, output := range check.LeanRows("godiff-slices-rotate-range", len(cases), &rows) {
			c := cases[i]
			c.Want = expected(c)
			if output != c.Want {
				return fmt.Errorf("range [%d,%d,%d), backing %v: Lean=%q want %q", c.A, c.Middle, c.B, c.Backing, output, c.Want)
			}
			if err := writer.Encode(c); err != nil {
				return err
			}
		}
		fmt.Printf("rotate ordered and comparator %s: %d Lean rows agree.\n", label, len(cases))
		return nil
	}
	edges := []fixture{{}, {Backing: []int{}, Mode: 1}}
	for n := 0; n <= 10; n++ {
		backing := make([]int, n+3)
		for i := range backing {
			backing[i] = 10*i - n
		}
		for a := 0; a <= n; a++ {
			for mode := range 2 {
				edges = append(edges, fixture{Backing: backing, A: a, Middle: a, B: a, Length: n, Mode: mode})
			}
			for middle := a + 1; middle < n; middle++ {
				for b := middle + 1; b <= n; b++ {
					for mode := range 2 {
						edges = append(edges, fixture{Backing: backing, A: a, Middle: middle, B: b, Length: n, Mode: mode})
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
	if err := os.Setenv("GOLEAN_SLICE_ROTATE_RANGE_CASES", file.Name()); err != nil {
		return err
	}
	const data = "difftest/slices/rotate-range/testdata/"
	return check.TestInPackage("slices", "TestGoLeanRotateRange", data+"bridge_test.go", data+"rotate_range_test.go")
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
	a := rng.Intn(n + 1)
	middle, b := a, a
	if n-a >= 2 && i%4 != 0 {
		middle = a + 1 + rng.Intn(n-a-1)
		b = middle + 1 + rng.Intn(n-middle)
		if i%5 == 0 {
			a, middle, b = 0, min(middle, n-1), n
		}
	}
	return fixture{Backing: backing, A: a, Middle: middle, B: b, Length: n, Mode: i % 2}
}

func expected(c fixture) string {
	if c.Backing == nil {
		return "nil|nil"
	}
	values := slices.Clone(c.Backing)
	position := c.A
	position += copy(values[position:c.B], c.Backing[c.Middle:c.B])
	copy(values[position:c.B], c.Backing[c.A:c.Middle])
	fields := make([]string, len(values))
	for i, value := range values {
		fields[i] = strconv.Itoa(value)
	}
	row := "slice:" + strings.Join(fields, ",")
	return row + "|" + row
}
