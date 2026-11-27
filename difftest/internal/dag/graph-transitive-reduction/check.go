package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type reductionCase struct {
	size  int
	codes []int
}

func (c reductionCase) covers() []int {
	covers := slices.Clone(c.codes)
	for from := range c.size {
		for to := range c.size {
			for via := range c.size {
				if c.codes[from*c.size+via] == 1 && c.codes[via*c.size+to] == 1 {
					covers[from*c.size+to] = 0
					break
				}
			}
		}
	}
	return covers
}

// The mutable matrix is independent of both Go's nested maps and the Lean list
// representation. Each lookup observes earlier deletions.
func (c reductionCase) sequentialResult() []int {
	current := slices.Clone(c.codes)
	for from := range c.size {
		for via := range c.size {
			if current[from*c.size+via] == 0 {
				continue
			}
			for to := range c.size {
				if current[via*c.size+to] == 1 {
					current[from*c.size+to] = 0
				}
			}
		}
	}
	return current
}

func (c reductionCase) insertionOrderIsTopological() bool {
	for from := range c.size {
		for to := range from {
			if c.codes[from*c.size+to] == 1 {
				return false
			}
		}
	}
	return true
}

func encodeEdges(edges []int) string {
	var row strings.Builder
	row.WriteString("edges")
	for _, edge := range edges {
		fmt.Fprintf(&row, " %d", edge)
	}
	return row.String()
}

func main() {
	counterexample := flag.Bool("cover-counterexample", false, "reproduce a mismatch with exact cover reduction")
	flag.Parse()
	if err := run(*counterexample); err != nil {
		log.Fatal(err)
	}
}

func generateCases(seeds []int64) ([]reductionCase, int) {
	var cases []reductionCase
	for size := range 5 {
		order := make([]int, size)
		for i := range order {
			order[i] = i
		}
		forEachOrder(order, 0, func(order []int) {
			for mask := range 1 << (size * (size - 1) / 2) {
				cases = append(cases, closedDAG(order, func() int {
					bit := mask & 1
					mask >>= 1
					return bit
				}))
			}
		})
	}
	boundaries := len(cases)
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			order := rng.Perm(rng.Intn(10))
			cases = append(cases, closedDAG(order, func() int {
				if rng.Intn(9) <= i%8 {
					return 1
				}
				return 0
			}))
		}
	}
	return cases, boundaries
}

func run(counterexample bool) error {
	seeds := []int64{17, 2718, 65537}
	var cases []reductionCase
	var boundaries int
	if counterexample {
		// Insertion order is 0,1,2,3. Topological order is 1,0,2,3.
		cases = []reductionCase{{4, []int{0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0}}}
		boundaries = 1
	} else {
		cases, boundaries = generateCases(seeds)
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&input, c.size)
		for _, code := range c.codes {
			fmt.Fprint(&input, " ", code)
		}
		fmt.Fprintln(&input)
	}
	dir, err := os.MkdirTemp("", "golean-reduction-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "inputs"), input.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_REDUCTION_DIR", dir); err != nil {
		return err
	}
	lean := check.LeanRows("godiff-internal-dag-graph-transitive-reduction", len(cases), &input)
	const fixture = "difftest/internal/dag/graph-transitive-reduction/testdata/reduction_test.go"
	if err := check.TestInPackage("internal/dag", "TestGoLeanReduction", fixture); err != nil {
		return err
	}
	output, err := os.ReadFile(filepath.Join(dir, "outputs"))
	if err != nil {
		return err
	}
	goRows := strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")
	if len(goRows) != len(cases) {
		return fmt.Errorf("Go produced %d rows, want %d", len(goRows), len(cases))
	}
	var topological int
	for i, c := range cases {
		sequential := c.sequentialResult()
		if want := encodeEdges(sequential); goRows[i] != want || lean[i] != want {
			return fmt.Errorf("case %d %+v, seeds %v: Go=%q Lean=%q sequential=%q", i, c, seeds, goRows[i], lean[i], want)
		}
		covers := c.covers()
		for pos, got := range sequential {
			if got > c.codes[pos] || got < covers[pos] {
				return fmt.Errorf("case %d %+v: changed non-deletable edge %d: Go=%q Lean=%q", i, c, pos, goRows[i], lean[i])
			}
		}
		if counterexample || c.insertionOrderIsTopological() {
			topological++
			if want := encodeEdges(covers); goRows[i] != want {
				return fmt.Errorf("case %d %+v: Go=%q Lean=%q exact cover=%q", i, c, goRows[i], lean[i], want)
			}
		}
	}
	fmt.Printf("internal/dag.Graph.TransitiveReduction: %d boundary and %d sampled cases agree, seeds %v. Exact cover also agrees on %d cases with topological insertion order.\n",
		boundaries, len(cases)-boundaries, seeds, topological)
	return nil
}

// Build a DAG whose hidden topological order can differ from insertion order,
// then close its incidence relation before calling the source method.
func closedDAG(order []int, nextEdge func() int) reductionCase {
	n := len(order)
	codes := make([]int, n*n)
	for i, from := range order {
		for _, to := range order[i+1:] {
			codes[from*n+to] = nextEdge()
		}
	}
	for via := range n {
		for from := range n {
			for to := range n {
				codes[from*n+to] |= codes[from*n+via] & codes[via*n+to]
			}
		}
	}
	return reductionCase{n, codes}
}

func forEachOrder(order []int, i int, visit func([]int)) {
	if i == len(order) {
		visit(order)
		return
	}
	for j := i; j < len(order); j++ {
		order[i], order[j] = order[j], order[i]
		forEachOrder(order, i+1, visit)
		order[i], order[j] = order[j], order[i]
	}
}
