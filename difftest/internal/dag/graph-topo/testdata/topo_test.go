package dag

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanTopo(t *testing.T) {
	dir := os.Getenv("GOLEAN_TOPO_DIR")
	if dir == "" {
		t.Fatal("missing GOLEAN_TOPO_DIR")
	}
	in, err := os.Open(filepath.Join(dir, "inputs"))
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()
	out, err := os.Create(filepath.Join(dir, "outputs"))
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		n, err := strconv.Atoi(fields[0])
		if err != nil || n < 0 || len(fields) != 1+n*n {
			t.Fatalf("invalid graph matrix: %q", fields)
		}
		graph := newGraph()
		for i := range n {
			graph.addNode(strconv.Itoa(i))
		}
		// A sparse old map can omit or use nil rows when no keys are present.
		if n == 0 {
			graph.edges = nil
		}
		for from := range n {
			source := strconv.Itoa(from)
			for to := range n {
				status := fields[1+from*n+to]
				if status != "0" {
					graph.edges[source][strconv.Itoa(to)] = status == "2"
				}
			}
			if len(graph.edges[source]) == 0 && from%2 == 0 {
				graph.edges[source] = nil
			}
		}

		nodes := graph.Topo()
		if nodes == nil {
			t.Fatalf("Topo returned nil for %v", fields)
		}
		if _, err := fmt.Fprintln(out, strings.Join(append([]string{"ok"}, nodes...), " ")); err != nil {
			t.Fatal(err)
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
}
