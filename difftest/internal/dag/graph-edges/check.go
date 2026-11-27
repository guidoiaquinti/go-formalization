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

type edgeCase struct {
	state int
	codes []int
}

func (c edgeCase) want() string {
	if c.state == 0 {
		return "panic"
	}
	if c.state < 4 {
		return "ok"
	}
	byRank := make([]string, len(c.codes))
	for key, code := range c.codes {
		if code%3 != 0 {
			byRank[code/3] = " " + strconv.Itoa(key)
		}
	}
	return "ok" + strings.Join(byRank, "")
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	seeds := []int64{17, 2718, 65537}
	cases := []edgeCase{{state: 0}, {state: 1}, {state: 2}, {state: 3}, {state: 4}}
	for _, ranks := range [][3]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}} {
		for mask := range 27 {
			cases = append(cases, edgeCase{4, []int{
				3*ranks[0] + mask%3, 3*ranks[1] + mask/3%3, 3*ranks[2] + mask/9,
			}})
		}
	}
	for _, n := range []int{15, 16, 17, 33} {
		codes := make([]int, n)
		for i := range codes {
			codes[i] = 3*(n-1-i) + 1 + i%2
		}
		cases = append(cases, edgeCase{4, codes})
	}
	boundaries := len(cases)
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for range 20_000 {
			codes := rng.Perm(1 + rng.Intn(32))
			for key, rank := range codes {
				codes[key] = 3*rank + rng.Intn(3)
			}
			cases = append(cases, edgeCase{4, codes})
		}
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&input, c.state)
		for _, code := range c.codes {
			fmt.Fprint(&input, " ", code)
		}
		fmt.Fprintln(&input)
	}
	dir, err := os.MkdirTemp("", "golean-edges-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "inputs"), input.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_EDGES_DIR", dir); err != nil {
		return err
	}
	lean := check.LeanRows("godiff-internal-dag-graph-edges", len(cases), &input)
	const fixture = "difftest/internal/dag/graph-edges/testdata/edges_test.go"
	if err := check.TestInPackage("internal/dag", "TestGoLeanEdges", fixture); err != nil {
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
	fmt.Printf("internal/dag.Graph.Edges: %d boundary and %d sampled cases agree, seeds %v.\n", boundaries, len(cases)-boundaries, seeds)
	return nil
}
