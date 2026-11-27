package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
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
	file, err := os.CreateTemp("", "golean-slices-choose-pivot-")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	encoder := json.NewEncoder(file)
	compare := func(label string, cases []fixture) error {
		var rows bytes.Buffer
		for _, c := range cases {
			fmt.Fprint(&rows, strconv.IntSize, " ", c.A, " ", c.B, " ", c.Mode, " ", c.FailAt, " ", c.Mask)
			for _, word := range c.Words {
				fmt.Fprint(&rows, " ", word)
			}
			fmt.Fprintln(&rows)
		}
		for i, got := range check.LeanRows("godiff-slices-choose-pivot", len(cases), &rows) {
			c := cases[i]
			c.Want = expected(c)
			if got != c.Want {
				return fmt.Errorf("%+v: Lean=%q", c, got)
			}
			if err := encoder.Encode(c); err != nil {
				return err
			}
		}
		fmt.Printf("choosePivot ordered and comparator %s: %d Lean rows agree.\n", label, len(cases))
		return nil
	}
	var edges []fixture
	for _, size := range []int{0, 1, 3, 4, 7, 8, 9, 12, 47, 48, 49, 50, 51, 52, 99, 100} {
		for _, a := range []int{0, 1, 5} {
			for pattern := range 4 {
				data := make([]int64, a+size+3)
				for i := range data {
					switch pattern {
					case 0:
						data[i] = int64(i)
					case 1:
						data[i] = int64(len(data) - i)
					case 2:
						data[i] = 0
					case 3:
						data[i] = int64((i * 7) % 11)
					}
				}
				for mode := range 3 {
					for failAt := -1; failAt < 12; failAt++ {
						for _, mask := range []int{0, 4095, 2730} {
							edges = append(edges, fixture{Words: data, A: a, B: a + size, Mode: mode, FailAt: failAt, Mask: mask})
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
	if err := os.Setenv("GOLEAN_SLICES_PIVOT_CASES", file.Name()); err != nil {
		return err
	}
	const data = "difftest/slices/choose-pivot/testdata/"
	return check.TestInPackage("slices", "TestGoLeanChoosePivot", data+"bridge_test.go", data+"pivot_test.go")
}

func sample(rng *rand.Rand, i int) fixture {
	words := make([]int64, rng.Intn(130))
	for j := range words {
		words[j] = int64(rng.Uint64())
		if i%3 == 0 {
			words[j] %= 5
		}
	}
	a := rng.Intn(len(words) + 1)
	return fixture{Words: words, A: a, B: a + rng.Intn(len(words)-a+1), Mode: rng.Intn(3),
		FailAt: rng.Intn(20) - 8, Mask: rng.Intn(4096)}
}

func expected(c fixture) string {
	return observe(c, false) + "|" + observe(c, true)
}

func observe(c fixture, callback bool) string {
	var calls []string
	render := func(status string) string {
		return status + "/" + strings.Join(calls, ";")
	}
	size := c.B - c.A
	q := size / 4
	indices := [3]int{c.A + q, c.A + 2*q, c.A + 3*q}
	if size < 8 {
		return render(fmt.Sprintf("ok:%d,1", indices[1]))
	}
	if callback && c.Mode == 0 {
		return render("nil")
	}
	swaps := 0
	failed := false
	selectMedian := func(positions [3]int) int {
		for _, left := range []int{0, 1, 0} {
			right := left + 1
			x, y := c.Words[positions[right]], c.Words[positions[left]]
			invert := x < y
			if callback {
				call := len(calls)
				calls = append(calls, fmt.Sprintf("%d,%d", x, y))
				if call == c.FailAt {
					failed = true
					return 0
				}
				if c.Mode == 2 {
					invert = c.Mask&(1<<call) != 0
				}
			}
			if invert {
				swaps++
				positions[left], positions[right] = positions[right], positions[left]
			}
		}
		return positions[1]
	}
	if size >= 50 {
		for slot, index := range indices {
			indices[slot] = selectMedian([3]int{index - 1, index, index + 1})
			if failed {
				return render(fmt.Sprintf("panic:%d", c.FailAt))
			}
		}
	}
	pivot := selectMedian(indices)
	if failed {
		return render(fmt.Sprintf("panic:%d", c.FailAt))
	}
	hint := 0
	switch swaps {
	case 0:
		hint = 1
	case 12:
		hint = 2
	}
	return render(fmt.Sprintf("ok:%d,%d", pivot, hint))
}
