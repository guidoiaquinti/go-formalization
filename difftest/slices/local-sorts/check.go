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
	Words                                []int64
	Kind, A, B, Root, Mode, FailAt, Mask int
	Want                                 string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	file, err := os.CreateTemp("", "golean-slices-local-sorts-")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	encoder := json.NewEncoder(file)
	compare := func(label string, cases []fixture) error {
		var rows bytes.Buffer
		for _, c := range cases {
			fmt.Fprint(&rows, c.Kind, " ", c.A, " ", c.B, " ", c.Root, " ", c.Mode, " ", c.FailAt, " ", c.Mask)
			for _, word := range c.Words {
				fmt.Fprint(&rows, " ", word)
			}
			fmt.Fprintln(&rows)
		}
		for i, got := range check.LeanRows("godiff-slices-local-sorts", len(cases), &rows) {
			c := cases[i]
			c.Want = expected(c)
			if got != c.Want {
				return fmt.Errorf("%+v: Lean=%q", c, got)
			}
			if err := encoder.Encode(c); err != nil {
				return err
			}
		}
		fmt.Printf("local-sort ordered and comparator %s: %d Lean rows agree.\n", label, len(cases))
		return nil
	}
	var edges []fixture
	for n := 0; n <= 8; n++ {
		for pattern := range 4 {
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
				}
			}
			for kind := range 3 {
				for _, a := range []int{0, 2} {
					for root := 0; root <= min(n, 3); root++ {
						for mode := range 3 {
							for _, stop := range []int{-1, 0, 1, 2, 3, 8} {
								edges = append(edges, fixture{Words: data, Kind: kind, A: a, B: a + n, Root: root, Mode: mode, FailAt: stop, Mask: 2730})
							}
						}
					}
				}
			}
		}
	}
	edges = append(edges,
		fixture{Words: []int64{3, 2, 1}, B: 3, Mode: 2, Mask: 4095, FailAt: 1},
		fixture{Words: []int64{0, 1, 2}, Kind: 1, B: 3, Mode: 1, FailAt: 0},
		fixture{Words: []int64{8, 3, 1, 9}, Kind: 2, A: 1, B: 3, Mode: 1, FailAt: -1},
	)
	if err := compare("boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for kind := range 3 {
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
	if err := os.Setenv("GOLEAN_SLICES_LOCAL_SORT_CASES", file.Name()); err != nil {
		return err
	}
	const data = "difftest/slices/local-sorts/testdata/"
	return check.TestInPackage("slices", "TestGoLeanLocalSorts", data+"bridge_test.go", data+"local_sorts_test.go")
}

func sample(rng *rand.Rand, kind int) fixture {
	data := make([]int64, rng.Intn(22))
	for i := range data {
		data[i] = int64(rng.Uint64())
	}
	a := rng.Intn(len(data) + 1)
	width := rng.Intn(len(data) - a + 1)
	return fixture{Words: data, Kind: kind, A: a, B: a + width, Root: rng.Intn(width + 1),
		Mode: rng.Intn(3), FailAt: rng.Intn(18) - 5, Mask: rng.Intn(4096)}
}

func expected(c fixture) string {
	return observe(c, false) + "|" + observe(c, true)
}

func observe(c fixture, callback bool) string {
	data := slices.Clone(c.Words)
	status := "ok"
	var calls []string
	less := func(i, j int) bool {
		if status != "ok" {
			return false
		}
		x, y := data[i], data[j]
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
	swap := func(i, j int) {
		if status == "ok" {
			data[i], data[j] = data[j], data[i]
		}
	}
	sift := func(root, size int) {
		for child := 2*root + 1; child < size; child = 2*root + 1 {
			best := child
			if child+1 < size && less(c.A+child, c.A+child+1) {
				best = child + 1
			}
			if !less(c.A+root, c.A+best) {
				return
			}
			swap(c.A+root, c.A+best)
			root = best
		}
	}
	switch {
	case !callback && c.Kind != 1:
		for i := c.A; i < c.B; i++ {
			minimum := i
			for j := i + 1; j < c.B; j++ {
				if data[j] < data[minimum] {
					minimum = j
				}
			}
			swap(i, minimum)
		}
	case c.Kind == 0:
		for i := c.A + 1; i < c.B && status == "ok"; i++ {
			for j := i; j > c.A && less(j, j-1); j-- {
				swap(j, j-1)
			}
		}
	case c.Kind == 1:
		sift(c.Root, c.B-c.A)
	default:
		n := c.B - c.A
		for root := (n - 1) / 2; root >= 0 && status == "ok"; root-- {
			sift(root, n)
		}
		for stop := n - 1; stop >= 0 && status == "ok"; stop-- {
			swap(c.A, c.A+stop)
			sift(0, stop)
		}
	}
	return render(status, data, calls)
}

func render(status string, data []int64, calls []string) string {
	values := make([]string, len(data))
	for i, v := range data {
		values[i] = fmt.Sprint(v)
	}
	return status + "/" + strings.Join(values, ",") + "/" + strings.Join(calls, ";")
}
