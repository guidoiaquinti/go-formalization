package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/bits"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Backing            []int
	A, B, Length, Mode int
	Want               string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	file, err := os.CreateTemp("", "golean-slices-break-patterns-")
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
			fmt.Fprint(&rows, tag, " ", strconv.IntSize, " ", c.A, " ", c.B, " ", c.Length, " ", c.Mode)
			for _, value := range c.Backing {
				fmt.Fprint(&rows, " ", value)
			}
			fmt.Fprintln(&rows)
		}
		for i, output := range check.LeanRows("godiff-slices-break-patterns", len(cases), &rows) {
			c := cases[i]
			c.Want = expected(c)
			if output != c.Want {
				return fmt.Errorf("offsets %d/%d length %d backing %v: Lean=%q want %q", c.A, c.B, c.Length, c.Backing, output, c.Want)
			}
			if err := writer.Encode(c); err != nil {
				return err
			}
		}
		fmt.Printf("breakPatterns ordered and comparator %s: %d Lean rows agree.\n", label, len(cases))
		return nil
	}
	edges := []fixture{
		{}, {A: 0, B: 8},
		{Backing: []int{}, A: 0, B: 8, Mode: 1},
		{A: 0, B: math.MaxInt},
		{Backing: []int{10, 20, 30, 40, 99}, Length: 4, A: 0, B: 8},
		{A: math.MinInt, B: math.MinInt + 8, Mode: 1},
		{A: math.MaxInt - 8, B: math.MaxInt},
		{A: math.MinInt, B: -1, Mode: 1},
	}
	for n := 0; n <= 28; n++ {
		backing := make([]int, n+3)
		for i := range backing {
			backing[i] = 10*i - n
		}
		for a := -3; a <= n+2; a++ {
			for b := -3; b <= n+4; b++ {
				edges = append(edges, fixture{Backing: backing, A: a, B: b, Length: n, Mode: (a + b) & 1})
			}
		}
	}
	for _, n := range []int{31, 32, 33, 63, 64, 65, 127, 128, 129, 255, 256, 257} {
		backing := make([]int, n+6)
		for i := range backing {
			backing[i] = i
		}
		edges = append(edges,
			fixture{Backing: backing, Length: n + 4, A: 2, B: n + 2},
			fixture{Backing: backing, Length: n + 4, A: 2, B: n + 2, Mode: 1})
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
	if err := os.Setenv("GOLEAN_SLICES_BREAK_PATTERNS_CASES", file.Name()); err != nil {
		return err
	}
	const data = "difftest/slices/break-patterns/testdata/"
	return check.TestInPackage("slices", "TestGoLeanBreakPatterns", data+"bridge_test.go", data+"break_patterns_test.go")
}

func sample(rng *rand.Rand, i int) fixture {
	n := rng.Intn(81)
	c := fixture{Length: n, Mode: i & 1}
	c.Backing = make([]int, n+rng.Intn(5))
	for j := range c.Backing {
		c.Backing[j] = int(rng.Uint64())
	}
	if i%3 == 0 {
		c.A = rng.Intn(n+7) - 3
		c.B = rng.Intn(n+15) - 3
	} else {
		c.A = rng.Intn(n + 1)
		c.B = c.A + rng.Intn(n-c.A+1)
	}
	return c
}

func expected(c fixture) string {
	backing := slices.Clone(c.Backing)
	status := "ok"
	size := c.B - c.A
	if size >= 8 {
		state := uint64(size)
		period := uint(1) << bits.Len(uint(size))
		middle := c.A + (size/4)*2
		for _, left := range []int{middle - 1, middle, middle + 1} {
			state ^= state << 13
			state ^= state >> 7
			state ^= state << 17
			right := c.A + int((uint(state)%period)%uint(size))
			if left < 0 || left >= c.Length || right < 0 || right >= c.Length {
				status = "panic"
				break
			}
			backing[left], backing[right] = backing[right], backing[left]
		}
	}
	cells := "nil"
	if backing != nil {
		values := make([]string, len(backing))
		for i, value := range backing {
			values[i] = strconv.Itoa(value)
		}
		cells = "slice:" + strings.Join(values, ",")
	}
	result := status + ":" + cells
	return result + "|" + result
}
