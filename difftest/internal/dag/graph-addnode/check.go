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

type input struct{ size, mode, query, existingRow int }

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var cases []input
	for size := range 4 {
		for mode := range 5 {
			for query := range size + 1 {
				for existingRow := range 2 {
					cases = append(cases, input{size, mode, query, existingRow})
				}
			}
		}
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for range 20_000 {
			size := rng.Intn(33)
			cases = append(cases, input{size, rng.Intn(5), rng.Intn(size + 1), rng.Intn(2)})
		}
	}

	var stdin bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&stdin, c.size, c.mode, c.query, c.existingRow)
	}
	dir, err := os.MkdirTemp("", "golean-addnode-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "inputs"), stdin.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_ADDNODE_DIR", dir); err != nil {
		return err
	}

	lean := check.LeanRows("godiff-internal-dag-graph-addnode", len(cases), &stdin)
	const fixture = "difftest/internal/dag/graph-addnode/testdata/addnode_test.go"
	if err := check.TestInPackage("internal/dag", "TestGoLeanAddNode", fixture); err != nil {
		return err
	}
	data, err := os.ReadFile(filepath.Join(dir, "outputs"))
	if err != nil {
		return err
	}
	rows := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	if len(rows) != len(cases) {
		return fmt.Errorf("Go returned %d rows, want %d", len(rows), len(cases))
	}
	for i, c := range cases {
		if expected := c.want(); rows[i] != expected || lean[i] != expected {
			return fmt.Errorf("input %+v: Go=%q Lean=%q want=%q", c, rows[i], lean[i], expected)
		}
	}
	fmt.Printf("internal/dag.Graph.addNode: %d boundary and %d sampled cases agree, seeds %v.\n",
		boundaries, len(cases)-boundaries, seeds)
	return nil
}

func (c input) want() string {
	if c.mode == 4 {
		return "panic nil"
	}
	rankPresent := c.mode != 1 && c.mode != 3
	edgesPresent := c.mode != 2 && c.mode != 3
	insert := rankPresent && c.query == c.size
	status := "panic"
	if insert && edgesPresent {
		status = "true"
	} else if rankPresent && !insert {
		status = "false"
	}
	size := c.size
	if insert {
		size++
	}
	names := make([]string, size)
	for i := range names {
		names[i] = strconv.Itoa(i)
	}
	nodes := "-"
	if size > 0 {
		nodes = strings.Join(names, ",")
	}
	ranks := "-"
	if rankPresent {
		out := append([]string{strconv.Itoa(size)}, names...)
		if !insert {
			out = append(out, "-")
		}
		ranks = strings.Join(out, ",")
	}
	edges := "-"
	if edgesPresent {
		count := c.size + c.existingRow
		if insert && c.existingRow == 0 {
			count++
		}
		out := []string{strconv.Itoa(count)}
		for i := range c.size + 1 {
			value := "1"
			if i == c.size && c.existingRow == 0 {
				value = "-"
			}
			if insert && i == c.query {
				value = "0"
			}
			out = append(out, value)
		}
		edges = strings.Join(out, ",")
	}
	return fmt.Sprintf("%s nodes=%s rank=%s edges=%s", status, nodes, ranks, edges)
}
