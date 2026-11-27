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
	Words                                 []int64
	Kind, A, B, Pivot, Mode, FailAt, Mask int
	Want                                  string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	file, err := os.CreateTemp("", "golean-slices-partition-")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	encoder := json.NewEncoder(file)
	compare := func(label string, cases []fixture) error {
		var rows bytes.Buffer
		for _, c := range cases {
			fmt.Fprint(&rows, c.Kind, " ", c.A, " ", c.B, " ", c.Pivot, " ", c.Mode, " ", c.FailAt, " ", c.Mask)
			for _, word := range c.Words {
				fmt.Fprint(&rows, " ", word)
			}
			fmt.Fprintln(&rows)
		}
		for i, got := range check.LeanRows("godiff-slices-partition", len(cases), &rows) {
			c := cases[i]
			c.Want = expected(c)
			if got != c.Want {
				return fmt.Errorf("%+v: Lean=%q", c, got)
			}
			if err := encoder.Encode(c); err != nil {
				return err
			}
		}
		fmt.Printf("partitions ordered and comparator %s: %d Lean rows agree.\n", label, len(cases))
		return nil
	}
	var edges []fixture
	for n := 1; n <= 8; n++ {
		for pattern := range 5 {
			data := make([]int64, n+4)
			for i := range data {
				switch pattern {
				case 0:
					data[i] = int64(i)
				case 1:
					data[i] = int64(-i)
				case 2:
					data[i] = int64(i % 3)
				case 3:
					data[i] = int64(-1 - i%4)
				default:
					data[i] = 7
				}
			}
			for kind := range 2 {
				for _, a := range []int{0, 2} {
					for pivot := a; pivot < a+n; pivot++ {
						for mode := range 3 {
							for _, stop := range []int{-1, 0, 1, 2, 3, 8} {
								edges = append(edges, fixture{Words: data, Kind: kind, A: a, B: a + n, Pivot: pivot, Mode: mode, FailAt: stop, Mask: 2730})
							}
						}
					}
				}
			}
		}
	}
	edges = append(edges,
		fixture{Words: []int64{8, 3, 9}, A: 1, B: 2, Pivot: 1, Mode: 1, FailAt: -1},
		fixture{Words: []int64{7, 3, 5}, B: 3, Pivot: 2, Mode: 1, FailAt: 0},
		fixture{Words: []int64{8, 3, 3, 3, 9}, Kind: 1, A: 1, B: 4, Pivot: 2, Mode: 1, FailAt: -1},
		fixture{Words: []int64{7, 5, 3}, Kind: 1, B: 3, Pivot: 2, Mode: 1, FailAt: 0},
	)
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
			if err := compare(fmt.Sprintf("kind=%d seed=%d", kind, seed), cases); err != nil {
				return err
			}
		}
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_SLICES_PARTITION_CASES", file.Name()); err != nil {
		return err
	}
	const data = "difftest/slices/partition/testdata/"
	return check.TestInPackage("slices", "TestGoLeanPartition", data+"bridge_test.go", data+"partition_test.go")
}

func sample(rng *rand.Rand, kind int) fixture {
	data := make([]int64, 1+rng.Intn(30))
	for i := range data {
		data[i] = int64(rng.Uint64())
		if i%4 != 0 {
			data[i] %= 5
		}
	}
	a := rng.Intn(len(data))
	b := a + 1 + rng.Intn(len(data)-a)
	pivot := a + rng.Intn(b-a)
	if kind == 1 && rng.Intn(2) == 0 {
		minimum := slices.Min(data[a:b])
		data[pivot] = minimum
	}
	return fixture{Words: data, Kind: kind, A: a, B: b, Pivot: pivot,
		Mode: rng.Intn(3), FailAt: rng.Intn(20) - 5, Mask: rng.Intn(4096)}
}

func expected(c fixture) string {
	return observe(c, false) + "|" + observe(c, true)
}

func observe(c fixture, callback bool) string {
	data := slices.Clone(c.Words)
	status := "ok"
	var calls []string
	less := func(x, y int64) bool {
		if !callback {
			return x < y
		}
		if c.Mode == 0 {
			status = "nil"
			return false
		}
		step := len(calls)
		calls = append(calls, fmt.Sprintf("%d,%d", x, y))
		if step == c.FailAt {
			status = fmt.Sprintf("panic:%d", step)
			return false
		}
		if c.Mode == 2 {
			return c.Mask&(1<<(step%12)) != 0
		}
		return x < y
	}
	equal := c.Kind == 1
	data[c.A], data[c.Pivot] = data[c.Pivot], data[c.A]
	pivot := data[c.A]
	leftGroup := func(value int64) bool {
		if equal {
			return !less(pivot, value)
		}
		return less(value, pivot)
	}
	low, high := c.A+1, c.B
	already := true
	for low < high && status == "ok" {
		switch {
		case leftGroup(data[low]):
			low++
		default:
			for low < high && status == "ok" {
				if leftGroup(data[high-1]) {
					if status == "ok" {
						data[low], data[high-1] = data[high-1], data[low]
						already = false
						low++
						high--
					}
					break
				}
				high--
			}
		}
	}
	if status == "ok" {
		if equal {
			status = fmt.Sprintf("ok:%d", low)
		} else {
			stop := high - 1
			data[c.A], data[stop] = data[stop], data[c.A]
			status = fmt.Sprintf("ok:%d,%t", stop, already)
		}
	}
	values := make([]string, len(data))
	for i, v := range data {
		values[i] = fmt.Sprint(v)
	}
	return status + "/" + strings.Join(values, ",") + "/" + strings.Join(calls, ";")
}
