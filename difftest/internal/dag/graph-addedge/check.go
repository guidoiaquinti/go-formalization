package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type addition struct {
	receiver, from, to int
	rows               [2]int
}

func (c addition) want() string {
	rows := c.rows
	if c.receiver != 2 || rows[c.from] < 2 {
		return "panic"
	}
	factor := 1
	if c.to == 1 {
		factor = 3
	}
	current := ((rows[c.from] - 2) / factor) % 3
	rows[c.from] += (2 - current) * factor
	return fmt.Sprintf("%d,%d", rows[0], rows[1])
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var input bytes.Buffer
	var cases []addition
	for receiver := range 3 {
		// Only an allocated outer map can hold distinct row states.
		rowCount := 1
		if receiver == 2 {
			rowCount = 11
		}
		for first := range rowCount {
			for second := range rowCount {
				for from := range 2 {
					for to := range 2 {
						c := addition{receiver, from, to, [2]int{first, second}}
						cases = append(cases, c)
						fmt.Fprintln(&input, receiver, from, to, first, second)
					}
				}
			}
		}
	}
	dir, err := os.MkdirTemp("", "golean-addedge-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "inputs"), input.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_ADDEDGE_DIR", dir); err != nil {
		return err
	}
	lean := check.LeanRows("godiff-internal-dag-graph-addedge", len(cases), &input)
	const fixture = "difftest/internal/dag/graph-addedge/testdata/addedge_test.go"
	if err := check.TestInPackage("internal/dag", "TestGoLeanAddEdge", fixture); err != nil {
		return err
	}
	goOutput, err := os.ReadFile(filepath.Join(dir, "outputs"))
	if err != nil {
		return err
	}
	goRows := strings.Fields(string(goOutput))
	if len(goRows) != len(cases) {
		return fmt.Errorf("Go produced %d rows, want %d", len(goRows), len(cases))
	}
	for i, c := range cases {
		if want := c.want(); goRows[i] != want || lean[i] != want {
			return fmt.Errorf("case %+v: Go=%q Lean=%q want=%q", c, goRows[i], lean[i], want)
		}
	}
	fmt.Printf("internal/dag.Graph.AddEdge: %d addition cases agree.\n", len(cases))
	return nil
}
