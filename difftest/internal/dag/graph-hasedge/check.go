package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type lookup struct {
	receiver, from, to int
	rows               [2]int
}

func (c lookup) want() string {
	if c.receiver == 0 {
		return "panic"
	}
	code := c.rows[c.from]
	if c.receiver == 1 || code < 2 {
		return "false"
	}
	digit := code - 2
	if c.to == 1 {
		digit /= 3
	}
	return strconv.FormatBool(digit%3 == 2)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var input bytes.Buffer
	var cases []lookup
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
						c := lookup{receiver, from, to, [2]int{first, second}}
						cases = append(cases, c)
						fmt.Fprintln(&input, receiver, from, to, first, second)
					}
				}
			}
		}
	}
	dir, err := os.MkdirTemp("", "golean-hasedge-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "inputs"), input.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_HASEDGE_DIR", dir); err != nil {
		return err
	}
	lean := check.LeanRows("godiff-internal-dag-graph-hasedge", len(cases), &input)
	const fixture = "difftest/internal/dag/graph-hasedge/testdata/hasedge_test.go"
	if err := check.TestInPackage("internal/dag", "TestGoLeanHasEdge", fixture); err != nil {
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
	fmt.Printf("internal/dag.Graph.HasEdge: %d lookup cases agree.\n", len(cases))
	return nil
}
