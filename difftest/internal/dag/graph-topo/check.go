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

type topoCase struct {
	size  int
	codes []int
}

func (c topoCase) want() string {
	// An explicit stack of enter and exit events computes the deterministic
	// reverse finishing order without recursive calls.
	type event struct {
		node int
		exit bool
	}
	var finished []int
	marked := make([]bool, c.size)
	for root := range c.size {
		stack := []event{{node: root}}
		for len(stack) != 0 {
			task := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if task.exit {
				marked[task.node] = true
				finished = append(finished, task.node)
				continue
			}
			if marked[task.node] {
				continue
			}
			stack = append(stack, event{node: task.node, exit: true})
			for dest := c.size - 1; dest >= 0; dest-- {
				if c.codes[task.node*c.size+dest] != 0 {
					stack = append(stack, event{node: dest})
				}
			}
		}
	}
	if len(finished) != c.size {
		panic("oracle did not finish each node once")
	}
	rank := make([]int, c.size)
	seen := make([]bool, c.size)
	var want strings.Builder
	want.WriteString("ok")
	for i, node := range finished {
		if seen[node] {
			panic("oracle repeated a node")
		}
		seen[node] = true
		rank[node] = c.size - 1 - i
		fmt.Fprintf(&want, " %d", finished[c.size-1-i])
	}
	for from := range c.size {
		for to := range c.size {
			if c.codes[from*c.size+to] != 0 && rank[from] >= rank[to] {
				panic("oracle returned a non-topological order")
			}
		}
	}
	return want.String()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	seeds := []int64{17, 2718, 65537}
	cases := []topoCase{{}}
	// Exhaust acyclic matrices through three nodes, allowing absent,
	// false-valued and true-valued entries. Exclude cyclic candidates.
	for size := 1; size <= 3; size++ {
		count := 1
		for range size * (size - 1) {
			count *= 3
		}
		for mask := range count {
			codes := make([]int, size*size)
			degree := make([]int, size)
			for from := range size {
				for to := range size {
					if from != to {
						codes[from*size+to], mask = mask%3, mask/3
						if codes[from*size+to] != 0 {
							degree[to]++
						}
					}
				}
			}
			var ready []int
			for node, d := range degree {
				if d == 0 {
					ready = append(ready, node)
				}
			}
			for i := 0; i < len(ready); i++ {
				from := ready[i]
				for to := range size {
					if codes[from*size+to] != 0 {
						degree[to]--
						if degree[to] == 0 {
							ready = append(ready, to)
						}
					}
				}
			}
			if len(ready) == size {
				cases = append(cases, topoCase{size, codes})
			}
		}
	}
	boundaries := len(cases)
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for range 20_000 {
			size := rng.Intn(9)
			codes := make([]int, size*size)
			height := rng.Perm(size)
			for from := range size {
				for to := range size {
					if height[from] < height[to] {
						codes[from*size+to] = rng.Intn(3)
					}
				}
			}
			cases = append(cases, topoCase{size, codes})
		}
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&input, c.size)
		for _, code := range c.codes {
			fmt.Fprint(&input, " ", code)
		}
		fmt.Fprintln(&input)
	}
	dir, err := os.MkdirTemp("", "golean-topo-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "inputs"), input.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_TOPO_DIR", dir); err != nil {
		return err
	}
	lean := check.LeanRows("godiff-internal-dag-graph-topo", len(cases), &input)
	const fixture = "difftest/internal/dag/graph-topo/testdata/topo_test.go"
	if err := check.TestInPackage("internal/dag", "TestGoLeanTopo", fixture); err != nil {
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
	for i, c := range cases {
		if want := c.want(); goRows[i] != want || lean[i] != want {
			return fmt.Errorf("case %d %+v, seeds %v: Go=%q Lean=%q want=%q", i, c, seeds, goRows[i], lean[i], want)
		}
	}
	fmt.Printf("internal/dag.Graph.Topo: %d boundary and %d sampled cases agree, seeds %v.\n", boundaries, len(cases)-boundaries, seeds)
	return nil
}
