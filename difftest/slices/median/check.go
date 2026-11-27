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
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Words                                []int64
	A, B, C, Center, Mode, FailAt, Count int
	Replies                              [3]int
	HasCounter                           bool
	Want                                 string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	file, err := os.CreateTemp("", "golean-slices-median-")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	encoder := json.NewEncoder(file)
	compare := func(label string, cases []fixture) error {
		var rows bytes.Buffer
		for _, c := range cases {
			present := 0
			if c.HasCounter {
				present = 1
			}
			fmt.Fprint(&rows, strconv.IntSize, " ", c.A, " ", c.B, " ", c.C, " ", c.Center, " ", present, " ", c.Count, " ", c.Mode, " ", c.FailAt, " ", c.Replies[0], " ", c.Replies[1], " ", c.Replies[2])
			for _, word := range c.Words {
				fmt.Fprint(&rows, " ", word)
			}
			fmt.Fprintln(&rows)
		}
		for i, got := range check.LeanRows("godiff-slices-median", len(cases), &rows) {
			c := cases[i]
			c.Want = expected(c)
			if got != c.Want {
				return fmt.Errorf("%+v: Lean=%q", c, got)
			}
			if err := encoder.Encode(c); err != nil {
				return err
			}
		}
		fmt.Printf("median ordered and comparator %s: %d Lean rows agree.\n", label, len(cases))
		return nil
	}
	edges := []fixture{{Words: []int64{30, 20, 10}, A: 0, B: 1, C: 2, Center: 1,
		Mode: 1, FailAt: 2, HasCounter: true, Count: 5}}
	for x := int64(-1); x <= 1; x++ {
		for y := int64(-1); y <= 1; y++ {
			for z := int64(-1); z <= 1; z++ {
				for a := range 3 {
					for b := range 3 {
						for c := range 3 {
							for mode := range 3 {
								for failAt := -1; failAt <= 2; failAt++ {
									for _, present := range []bool{false, true} {
										edges = append(edges, fixture{Words: []int64{x, y, z}, A: a, B: b, C: c, Center: 1,
											Mode: mode, FailAt: failAt, HasCounter: present, Count: math.MaxInt,
											Replies: [3]int{math.MinInt, math.MaxInt, -1}})
									}
								}
							}
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
	if err := os.Setenv("GOLEAN_SLICES_MEDIAN_CASES", file.Name()); err != nil {
		return err
	}
	const data = "difftest/slices/median/testdata/"
	return check.TestInPackage("slices", "TestGoLeanMedian", data+"bridge_test.go", data+"median_test.go")
}

func sample(rng *rand.Rand, i int) fixture {
	words := make([]int64, 3+rng.Intn(12))
	for j := range words {
		words[j] = int64(rng.Uint64())
	}
	return fixture{Words: words, A: rng.Intn(len(words)), B: rng.Intn(len(words)), C: rng.Intn(len(words)),
		Center: 1 + rng.Intn(len(words)-2), Mode: rng.Intn(3), FailAt: rng.Intn(7) - 4,
		HasCounter: i%4 != 0, Count: int(rng.Uint64()),
		Replies: [3]int{int(rng.Uint64()), int(rng.Uint64()), int(rng.Uint64())}}
}

func expected(c fixture) string {
	return strings.Join([]string{
		observe(c, [3]int{c.A, c.B, c.C}, false),
		observe(c, [3]int{c.A, c.B, c.C}, true),
		observe(c, [3]int{c.Center - 1, c.Center, c.Center + 1}, false),
		observe(c, [3]int{c.Center - 1, c.Center, c.Center + 1}, true),
	}, "|")
}

func observe(c fixture, order [3]int, callback bool) string {
	count := c.Count
	var calls []string
	render := func(result string) string {
		counter := "nil"
		if c.HasCounter {
			counter = strconv.Itoa(count)
		}
		return result + "/" + counter + "/" + strings.Join(calls, ";")
	}
	if callback && c.Mode == 0 {
		return render("nil")
	}
	for call, left := range []int{0, 1, 0} {
		right := left + 1
		x, y := c.Words[order[right]], c.Words[order[left]]
		swap := x < y
		if callback {
			calls = append(calls, fmt.Sprintf("%d,%d", x, y))
			if c.FailAt == call {
				return render(fmt.Sprintf("panic:%d", call))
			}
			if c.Mode == 2 {
				swap = c.Replies[call] < 0
			}
		}
		if swap {
			if !c.HasCounter {
				return render("nil")
			}
			count++
			order[left], order[right] = order[right], order[left]
		}
	}
	return render(fmt.Sprintf("ok:%d", order[1]))
}
