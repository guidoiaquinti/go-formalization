package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	r, length int
	cells     []int
}

func row(prefix string, cells []int) string {
	var result strings.Builder
	result.WriteString(prefix)
	for _, v := range cells {
		fmt.Fprintf(&result, " %d", v)
	}
	return result.String()
}

func expected(c input) string {
	// A positive mathematical shift exceeding MaxInt wraps to a negative Go int.
	// A right count exceeding length also produces a negative left count.
	if c.r < c.length-math.MaxInt || c.r > c.length {
		return row("panic", c.cells)
	}
	out := append([]int(nil), c.cells...)
	if c.r < 0 {
		prefix := c.length - c.r
		if prefix <= len(c.cells) {
			for i := range c.cells[:prefix] {
				out[i] = c.cells[prefix-1-i]
			}
		}
		return row("panic", out)
	}
	for i, value := range c.cells[:c.length] {
		out[(i+c.r)%c.length] = value
	}
	return row("ok", out)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []input{
		{-1, 0, nil}, {0, 0, nil}, {1, 0, nil}, {math.MinInt, 0, nil},
		{-3, 3, []int{1, 2, 3, 4, 5}}, {-2, 3, []int{1, 2, 3, 4, 5}},
		{-1, 3, []int{1, 2, 3, 4, 5}}, {0, 3, []int{1, 2, 3, 4, 5}},
		{1, 3, []int{1, 2, 3, 4, 5}}, {2, 3, []int{1, 2, 3, 4, 5}},
		{3, 3, []int{1, 2, 3, 4, 5}}, {4, 3, []int{1, 2, 3, 4, 5}},
		{math.MinInt, 3, []int{1, 2, 3, 4, 5}}, {math.MaxInt, 3, []int{1, 2, 3, 4, 5}},
		{3 - math.MaxInt, 3, []int{1, 2, 3, 4, 5}}, {2 - math.MaxInt, 3, []int{1, 2, 3, 4, 5}},
		{-2, 0, []int{1, 2, 3}}, {1, 4, []int{1, 2, 1, 2}},
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			capacity := rng.Intn(65)
			length := rng.Intn(capacity + 1)
			cells := make([]int, capacity)
			for i := range cells {
				cells[i] = rng.Intn(101) - 50
			}
			r := rng.Intn(capacity+3) - (capacity - length + 1)
			if i%17 == 0 {
				r = int(rng.Uint64())
			}
			cases = append(cases, input{r, length, cells})
		}
	}
	var wire, leanWire bytes.Buffer
	for _, c := range cases {
		encoded := row(fmt.Sprintf("%d %d", c.r, c.length), c.cells)
		fmt.Fprintln(&wire, encoded)
		fmt.Fprintln(&leanWire, strconv.IntSize, encoded)
	}
	dir, err := os.MkdirTemp("", "golean-rotate-right-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "inputs"), wire.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_ROTATE_RIGHT_DIR", dir); err != nil {
		return err
	}
	lean := check.LeanRows("godiff-slices-rotate-right", len(cases), &leanWire)
	const fixture = "difftest/slices/rotate-right/testdata/rotate_right_test.go"
	if err := check.TestInPackage("slices", "TestGoLeanRotateRight", fixture, "difftest/slices/rotate-right/testdata/bridge_test.go"); err != nil {
		return err
	}
	out, err := os.ReadFile(filepath.Join(dir, "outputs"))
	if err != nil {
		return err
	}
	actual := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(actual) != len(cases) {
		return fmt.Errorf("Go produced %d rows, want %d", len(actual), len(cases))
	}
	for i, c := range cases {
		want := expected(c)
		if actual[i] != want || lean[i] != want {
			return fmt.Errorf("case %+v: Go=%q Lean=%q want=%q", c, actual[i], lean[i], want)
		}
	}
	fmt.Printf("slices.rotateRight: %d boundary and 60000 sampled cases agree; seeds %v\n", boundaries, seeds)
	return nil
}
