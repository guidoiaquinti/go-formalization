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

const slots = 64

type interval struct{ start, length int }

type input struct {
	a, b     interval
	separate bool
}

func expected(c input) string {
	if c.b.length == 0 {
		return "bounds panic"
	}
	offset := c.b.start - c.a.start
	if c.separate || offset < 0 || offset >= c.a.length {
		return "panic needle not found"
	}
	return fmt.Sprintf("index %d", offset)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var intervals []interval
	for start := range 9 {
		for length := range 9 - start {
			intervals = append(intervals, interval{start, length})
		}
	}
	var cases []input
	for _, a := range intervals {
		for _, b := range intervals {
			cases = append(cases, input{a, b, false}, input{a, b, true})
		}
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		span := func() interval {
			start := rng.Intn(slots + 1)
			return interval{start, rng.Intn(slots - start + 1)}
		}
		for i := range 20_000 {
			a, b := span(), span()
			if a.length != 0 && i%2 == 0 {
				b.start = a.start + rng.Intn(a.length)
				b.length = 1 + rng.Intn(a.start+a.length-b.start)
			}
			cases = append(cases, input{a, b, i%4 == 3})
		}
	}
	var goInput, leanInput bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&goInput, c.a.start, c.a.length, c.b.start, c.b.length, c.separate)
		bStart := c.b.start
		if c.separate {
			bStart += slots
		}
		fmt.Fprintln(&leanInput, c.a.start, c.a.length, bStart, c.b.length)
	}
	dir, err := os.MkdirTemp("", "golean-slices-start-idx-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "inputs"), goInput.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_START_IDX_DIR", dir); err != nil {
		return err
	}
	lean := check.LeanRows("godiff-slices-start-idx", len(cases), &leanInput)
	const fixture = "difftest/slices/start-idx/testdata/"
	if err := check.TestInPackage("slices", "TestGoLeanStartIdx", fixture+"start_idx_test.go", fixture+"bridge_test.go"); err != nil {
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
	fmt.Printf("slices.startIdx: %d boundary and 60000 sampled cases agree; seeds %v\n", boundaries, seeds)
	return nil
}
