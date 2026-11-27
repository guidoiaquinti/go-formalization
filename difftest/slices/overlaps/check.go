package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const slots = 64

type interval struct{ start, length int }

type input struct {
	a, b     interval
	size     int
	separate bool
}

func expected(c input) bool {
	if c.separate || c.size == 0 {
		return false
	}
	for i := range c.a.length {
		offset := c.a.start + i - c.b.start
		if offset >= 0 && offset < c.b.length {
			return true
		}
	}
	return false
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	sizes := []int{0, 1, 3, 8}
	var intervals []interval
	for start := range 9 {
		for length := range 9 - start {
			intervals = append(intervals, interval{start, length})
		}
	}
	var cases []input
	for _, a := range intervals {
		for _, b := range intervals {
			for _, size := range sizes {
				cases = append(cases, input{a, b, size, false}, input{a, b, size, true})
			}
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
		for range 20_000 {
			cases = append(cases, input{span(), span(), sizes[rng.Intn(len(sizes))], rng.Intn(2) == 0})
		}
	}
	var goInput, leanInput bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&goInput, c.a.start, c.a.length, c.b.start, c.b.length, c.size, c.separate)
		bStart := c.b.start
		if c.separate {
			bStart += slots
		}
		fmt.Fprintln(&leanInput, c.a.start*c.size, c.a.length, bStart*c.size, c.b.length, c.size)
	}
	dir, err := os.MkdirTemp("", "golean-slices-overlaps-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "inputs"), goInput.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_OVERLAPS_DIR", dir); err != nil {
		return err
	}
	lean := check.LeanRows("godiff-slices-overlaps", len(cases), &leanInput)
	const fixture = "difftest/slices/overlaps/testdata/"
	if err := check.TestInPackage("slices", "TestGoLeanOverlaps", fixture+"overlaps_test.go", fixture+"bridge_test.go"); err != nil {
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
		want := strconv.FormatBool(expected(c))
		if actual[i] != want || lean[i] != want {
			return fmt.Errorf("case %+v: Go=%q Lean=%q want=%q", c, actual[i], lean[i], want)
		}
	}
	fmt.Printf("slices.overlaps: %d boundary and 60000 sampled cases agree; seeds %v\n", boundaries, seeds)
	return nil
}
