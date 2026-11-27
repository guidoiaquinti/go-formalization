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

type deletion struct {
	receiver, from, to int
	rows               [2]int
}

func (c deletion) want() string {
	if c.receiver == 0 {
		return "panic"
	}
	if c.receiver == 1 {
		return "nil"
	}
	rows := c.rows
	if code := rows[c.from]; code >= 2 {
		factor := 1
		if c.to == 1 {
			factor = 3
		}
		rows[c.from] -= ((code - 2) / factor) % 3 * factor
	}
	return fmt.Sprintf("%d,%d", rows[0], rows[1])
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var input bytes.Buffer
	var cases []deletion
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
						c := deletion{receiver, from, to, [2]int{first, second}}
						cases = append(cases, c)
						fmt.Fprintln(&input, receiver, from, to, first, second)
					}
				}
			}
		}
	}
	dir, err := os.MkdirTemp("", "golean-deledge-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "inputs"), input.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_DELEDGE_DIR", dir); err != nil {
		return err
	}
	lean := check.LeanRows("godiff-internal-dag-graph-deledge", len(cases), &input)
	const fixture = "difftest/internal/dag/graph-deledge/testdata/deledge_test.go"
	if err := check.TestInPackage("internal/dag", "TestGoLeanDelEdge", fixture); err != nil {
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
	fmt.Printf("internal/dag.Graph.DelEdge: %d deletion cases agree.\n", len(cases))
	return nil
}
