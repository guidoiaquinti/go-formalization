package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type scenario struct {
	start     int
	decisions uint64
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var cases []scenario
	for start := range 128 {
		for decisions := range 16 {
			cases = append(cases, scenario{start, uint64(decisions)})
		}
	}
	for _, start := range []int{1023, 1024, math.MaxInt - 1, math.MaxInt} {
		cases = append(cases, scenario{start, 0}, scenario{start, math.MaxUint64})
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			decisions := rng.Uint64()
			if i%4 == 0 {
				decisions = math.MaxUint64
			}
			cases = append(cases, scenario{rng.Int(), decisions})
		}
	}

	var lean, fixture bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&lean, c.start, c.decisions)
		fmt.Fprintln(&fixture, c.start, c.decisions, c.expected())
	}
	for i, got := range check.LeanRows("godiff-container-heap-up", len(cases), &lean) {
		if want := cases[i].expected(); got != want {
			return fmt.Errorf("Lean trace for %+v was %q, want %q", cases[i], got, want)
		}
	}

	dir, err := os.MkdirTemp("", "golean-heap-up-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	file := filepath.Join(dir, "cases")
	if err := os.WriteFile(file, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_HEAP_UP_CASES", file); err != nil {
		return err
	}
	const overlay = "difftest/container/heap/up/testdata/up_test.go"
	if err := check.TestInPackage("container/heap", "TestGoLeanUp", overlay); err != nil {
		return err
	}
	fmt.Printf("container/heap.up: %d boundaries and %d random scenarios agree, seeds %v.\n",
		boundaries, len(cases)-boundaries, seeds)
	return nil
}

// Numbering nodes from one gives each ancestor by dropping another low bit.
func (c scenario) expected() string {
	trace := bytes.NewBufferString("trace")
	node := uint64(c.start) + 1
	for depth := 0; node>>depth > 1; depth++ {
		child := (node >> depth) - 1
		parent := (node >> (depth + 1)) - 1
		fmt.Fprintf(trace, " L:%d:%d", child, parent)
		if c.decisions&(1<<depth) == 0 {
			break
		}
		fmt.Fprintf(trace, " S:%d:%d", parent, child)
	}
	return trace.String()
}
