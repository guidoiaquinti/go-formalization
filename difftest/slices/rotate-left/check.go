package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
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
	if c.r < 0 || c.r > len(c.cells) {
		return row("panic", c.cells)
	}
	if c.r > c.length {
		out := append([]int(nil), c.cells...)
		for i := range c.cells[:c.r] {
			out[i] = c.cells[c.r-1-i]
		}
		return row("panic", out)
	}
	out := append([]int(nil), c.cells...)
	for i := range out[:c.length] {
		out[i] = c.cells[(i+c.r)%c.length]
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
		{-1, 0, nil}, {0, 0, nil}, {1, 0, nil},
		{-1, 1, []int{1}}, {0, 1, []int{1}}, {1, 1, []int{1}}, {2, 1, []int{1}},
		{0, 5, []int{1, 2, 3, 4, 5}}, {1, 5, []int{1, 2, 3, 4, 5}},
		{2, 5, []int{1, 2, 3, 4, 5}}, {4, 5, []int{1, 2, 3, 4, 5}}, {5, 5, []int{1, 2, 3, 4, 5}},
		{1, 4, []int{1, 2, 1, 2}}, {2, 4, []int{1, 2, 3, 4}},
		{2, 3, []int{1, 2, 3, 4, 5}}, {4, 3, []int{1, 2, 3, 4, 5}},
		{5, 3, []int{1, 2, 3, 4, 5}}, {6, 3, []int{1, 2, 3, 4, 5}},
		{2, 0, []int{1, 2, 3}},
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for range 20_000 {
			capacity := rng.Intn(65)
			length := rng.Intn(capacity + 1)
			cells := make([]int, capacity)
			for i := range cells {
				cells[i] = rng.Intn(101) - 50
			}
			cases = append(cases, input{rng.Intn(capacity+3) - 1, length, cells})
		}
	}
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&wire, row(fmt.Sprintf("%d %d", c.r, c.length), c.cells))
	}
	dir, err := os.MkdirTemp("", "golean-rotate-left-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "inputs"), wire.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_ROTATE_LEFT_DIR", dir); err != nil {
		return err
	}
	lean := check.LeanRows("godiff-slices-rotate-left", len(cases), &wire)
	const fixture = "difftest/slices/rotate-left/testdata/rotate_left_test.go"
	if err := check.TestInPackage("slices", "TestGoLeanRotateLeft", fixture, "difftest/slices/rotate-left/testdata/bridge_test.go"); err != nil {
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
	fmt.Printf("slices.rotateLeft: %d boundary and 60000 sampled cases agree; seeds %v\n", boundaries, seeds)
	return nil
}
