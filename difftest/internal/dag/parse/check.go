package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type testCase struct {
	Text, Expected string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []testCase{
		{"", "graph 0"},
		{"# comment\n", "graph 0"},
		{"NONE < a;", "graph 1 1 97 0"},
		{"NONE < \xff,é,\x00;", "graph 3 1 255 2 195 169 1 0 0 0 0 0 0 0 0 0 0"},
		{"NONE < a,b; a,b < c;", "graph 3 1 97 1 98 1 99 0 0 0 0 0 0 1 1 0"},
		{"NONE < a; a < b; b < c; a !< c;", "semantic " + encode("graph edge assertion failed: a !< c")},
		{"NONE < a; NONE < a;", "semantic " + encode("multiple definitions for a")},
		{"unknown < a;", "semantic " + encode("use of unknown before its definition")},
		{"NONE < NONE;", "semantic " + encode("NONE cannot be a predecessor")},
		{"NONE < a; a < a;", "semantic " + encode("multiple definitions for a")},
		{"NONE < a,b; a < b; b < a;", "semantic " + encode("multiple definitions for b\nmultiple definitions for a\ngraph cycle: b < a < b\ngraph cycle: a < b < a")},
	}
	// Check a syntax error against its exact source diagnostic.
	cases = append(cases, testCase{"!", "syntax " + encode("parsing graph: line 1: syntax error: unexpected token ! near ")})
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			cases = append(cases, sample(rng, i))
		}
	}

	var stdin bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&stdin, strconv.IntSize)
		for j := range len(c.Text) {
			fmt.Fprint(&stdin, " ", c.Text[j])
		}
		fmt.Fprintln(&stdin)
	}
	for i, got := range check.LeanRows("godiff-internal-dag-parse", len(cases), &stdin) {
		if c := cases[i]; got != c.Expected {
			return fmt.Errorf("case %d %q: Lean=%q want=%q", i, c.Text, got, c.Expected)
		}
	}
	dir, err := os.MkdirTemp("", "golean-parse-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	wire := make([]struct {
		Text     []byte
		Expected string
	}, len(cases))
	for i, c := range cases {
		wire[i].Text, wire[i].Expected = []byte(c.Text), c.Expected
	}
	data, err := json.Marshal(wire)
	if err != nil {
		return err
	}
	file := filepath.Join(dir, "cases.json")
	if err := os.WriteFile(file, data, 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_DAG_PARSE_CASES", file); err != nil {
		return err
	}
	const fixture = "difftest/internal/dag/parse/testdata/parse_test.go"
	if err := check.TestInPackage("internal/dag", "TestGoLeanParse", fixture); err != nil {
		return err
	}
	fmt.Printf("internal/dag.Parse: %d boundaries and %d sampled graphs agree, seeds %v.\n",
		boundaries, len(cases)-boundaries, seeds)
	return nil
}

func sample(rng *rand.Rand, index int) testCase {
	n := 1 + rng.Intn(8)
	var source strings.Builder
	nodes := make([]string, n)
	direct := make([][]int, n)
	for i := range nodes {
		nodes[i] = fmt.Sprintf("n%d", rng.Uint32())
		nodes[i] += fmt.Sprintf("_%d", i)
		var predecessors []string
		for j := range i {
			if rng.Intn(3) == 0 {
				direct[i] = append(direct[i], j)
				predecessors = append(predecessors, nodes[j])
			}
		}
		if len(predecessors) == 0 {
			predecessors = []string{"NONE"}
		}
		fmt.Fprintf(&source, "%s < %s;\n", strings.Join(predecessors, ","), nodes[i])
	}
	// Search direct edges independently from each node to compute closure.
	closure := make([][]bool, n)
	for i := range nodes {
		closure[i] = make([]bool, n)
		var visit func(int)
		visit = func(v int) {
			for _, to := range direct[v] {
				if !closure[i][to] {
					closure[i][to] = true
					visit(to)
				}
			}
		}
		visit(i)
	}
	if n > 1 && index%2 != 0 {
		from := 1 + rng.Intn(n-1)
		to := rng.Intn(from)
		fmt.Fprintf(&source, "%s !< %s;", nodes[to], nodes[from])
		if closure[from][to] {
			message := fmt.Sprintf("graph edge assertion failed: %s !< %s", nodes[to], nodes[from])
			return testCase{source.String(), "semantic " + encode(message)}
		}
	}
	rows := []string{"graph", strconv.Itoa(n)}
	for _, node := range nodes {
		rows = append(rows, encode(node))
	}
	for _, row := range closure {
		for _, edge := range row {
			if edge {
				rows = append(rows, "1")
			} else {
				rows = append(rows, "0")
			}
		}
	}
	return testCase{source.String(), strings.Join(rows, " ")}
}

func encode(s string) string {
	var row strings.Builder
	fmt.Fprint(&row, len(s))
	for i := range len(s) {
		fmt.Fprint(&row, " ", s[i])
	}
	return row.String()
}
