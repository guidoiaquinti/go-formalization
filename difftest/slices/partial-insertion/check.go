package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Words                    []int64
	A, B, Mode, FailAt, Mask int
	Want                     string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	file, err := os.CreateTemp("", "golean-slices-partial-insertion-")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	encoder := json.NewEncoder(file)
	compare := func(label string, cases []fixture) error {
		var rows bytes.Buffer
		for _, c := range cases {
			fmt.Fprint(&rows, c.A, " ", c.B, " ", c.Mode, " ", c.FailAt, " ", c.Mask)
			for _, word := range c.Words {
				fmt.Fprint(&rows, " ", word)
			}
			fmt.Fprintln(&rows)
		}
		for i, got := range check.LeanRows("godiff-slices-partial-insertion", len(cases), &rows) {
			c := cases[i]
			c.Want = expected(c)
			if got != c.Want {
				return fmt.Errorf("%+v: Lean=%q", c, got)
			}
			if err := encoder.Encode(c); err != nil {
				return err
			}
		}
		fmt.Printf("partial-insertions ordered and comparator %s: %d Lean rows agree.\n", label, len(cases))
		return nil
	}
	var edges []fixture
	for _, size := range []int{0, 1, 2, 3, 7, 47, 48, 49, 50, 51, 52, 65, 100} {
		for _, a := range []int{0, 1, 2, 7} {
			for pattern := range 5 {
				data := make([]int64, a+size+2)
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
					default:
						data[i] = int64(i + (i%11)*5)
					}
				}
				for mode := range 3 {
					for _, stop := range []int{-1, 0, 1, 2, 3, 7, 22, 48, 49, 50, 88} {
						for _, mask := range []int{0, 2730, 4095} {
							edges = append(edges, fixture{Words: data, A: a, B: a + size, Mode: mode, FailAt: stop, Mask: mask})
						}
					}
				}
			}
		}
	}
	cross := []int64{900, 901, 10, 11, 0}
	for i := range 47 {
		cross = append(cross, int64(i+13))
	}
	panicAfterSwap := []int64{2, 1}
	for range 48 {
		panicAfterSwap = append(panicAfterSwap, 3)
	}
	edges = append(edges, fixture{Words: cross, A: 2, B: 52, Mode: 1, FailAt: -1},
		fixture{Words: panicAfterSwap, B: 50, Mode: 2, FailAt: 1, Mask: 4095})
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
	if err := os.Setenv("GOLEAN_SLICES_PARTIAL_INSERTION_CASES", file.Name()); err != nil {
		return err
	}
	const data = "difftest/slices/partial-insertion/testdata/"
	return check.TestInPackage("slices", "TestGoLeanPartialInsertion", data+"bridge_test.go", data+"partial_insertion_test.go")
}

func sample(rng *rand.Rand, i int) fixture {
	data := make([]int64, rng.Intn(122))
	for j := range data {
		if i%3 == 0 {
			data[j] = int64(j)
		} else {
			data[j] = int64(rng.Uint64())
		}
	}
	if len(data) > 0 && i%3 == 0 {
		for range rng.Intn(8) {
			data[rng.Intn(len(data))] = int64(rng.Uint64())
		}
	}
	a := rng.Intn(len(data) + 1)
	if i%2 == 0 {
		a = 0
	}
	return fixture{Words: data, A: a, B: a + rng.Intn(len(data)-a+1),
		Mode: rng.Intn(3), FailAt: rng.Intn(95) - 8, Mask: rng.Intn(4096)}
}

func expected(c fixture) string {
	return observe(c, false) + "|" + observe(c, true)
}

func observe(c fixture, callback bool) string {
	data := slices.Clone(c.Words)
	failure := ""
	var calls []string
	inverted := func(i int) bool {
		x, y := data[i], data[i-1]
		if !callback {
			return x < y
		}
		if c.Mode == 0 {
			failure = "nil"
			return false
		}
		step := len(calls)
		calls = append(calls, fmt.Sprintf("%d,%d", x, y))
		if step == c.FailAt {
			failure = fmt.Sprintf("panic:%d", step)
			return false
		}
		if c.Mode == 2 {
			return c.Mask&(1<<(step%12)) != 0
		}
		return x < y
	}
	swap := func(i int) { data[i-1], data[i] = data[i], data[i-1] }
	shift := func(start, step int, valid func(int) bool) {
		for j := start; failure == "" && valid(j) && inverted(j); j += step {
			swap(j)
		}
	}
	index := c.A + 1
	done := false
	for repairs := 0; repairs < 5; repairs++ {
		for index < c.B && failure == "" && !inverted(index) {
			index++
		}
		if failure != "" {
			break
		}
		if index == c.B {
			done = true
			break
		}
		if c.B-c.A < 50 {
			break
		}
		swap(index)
		if index-c.A >= 2 {
			shift(index-1, -1, func(i int) bool { return i > 0 })
		}
		shift(index+1, 1, func(i int) bool { return i < c.B })
		if failure != "" {
			break
		}
	}
	status := failure
	if status == "" {
		status = fmt.Sprintf("ok:%t", done)
	}
	values := make([]string, len(data))
	for i, v := range data {
		values[i] = fmt.Sprint(v)
	}
	return status + "/" + strings.Join(values, ",") + "/" + strings.Join(calls, ";")
}
