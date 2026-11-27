package dag

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"testing"
)

func TestGoLeanAddEdge(t *testing.T) {
	dir := os.Getenv("GOLEAN_ADDEDGE_DIR")
	if dir == "" {
		t.Fatal("missing GOLEAN_ADDEDGE_DIR")
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

	keys := [2]string{"", "\x00"}
	scanner := bufio.NewReader(in)
	for {
		var receiver, from, to int
		var codes [2]int
		_, err := fmt.Fscan(scanner, &receiver, &from, &to, &codes[0], &codes[1])
		if err == io.EOF {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		var graph *Graph
		switch receiver {
		case 1:
			graph = &Graph{}
		case 2:
			graph = &Graph{edges: make(map[string]map[string]bool)}
			for i, code := range codes {
				if code == 1 {
					graph.edges[keys[i]] = nil
				} else if code >= 2 {
					row := make(map[string]bool)
					digits := code - 2
					for _, key := range keys {
						if value := digits % 3; value != 0 {
							row[key] = value == 2
						}
						digits /= 3
					}
					graph.edges[keys[i]] = row
				}
			}
		}
		if graph != nil {
			graph.Nodes = []string{"sentinel"}
			graph.byLabel = map[string]int{"sentinel": 1}
		}
		result := goLeanAddEdgeOutcome(t, graph, keys, from, to)
		if graph != nil && (!slices.Equal(graph.Nodes, []string{"sentinel"}) ||
			len(graph.byLabel) != 1 || graph.byLabel["sentinel"] != 1) {
			t.Fatal("AddEdge modified unrelated graph fields")
		}
		if _, err := fmt.Fprintln(out, result); err != nil {
			t.Fatal(err)
		}
	}
}

func goLeanAddEdgeOutcome(t *testing.T, graph *Graph, keys [2]string, from, to int) (result string) {
	t.Helper()
	defer func() {
		if recovered := recover(); recovered != nil {
			if _, ok := recovered.(runtime.Error); !ok {
				panic(recovered)
			}
			result = "panic"
		}
	}()
	graph.AddEdge(keys[from], keys[to])
	if graph.edges == nil {
		t.Error("successful AddEdge left a nil edge map")
		return "nil"
	}
	var codes [2]int
	rowCount := 0
	for i, key := range keys {
		row, present := graph.edges[key]
		if present {
			rowCount++
		}
		switch {
		case !present:
		case row == nil:
			codes[i] = 1
		default:
			codes[i] = 2
			factor, entryCount := 1, 0
			for _, destination := range keys {
				if value, ok := row[destination]; ok {
					entryCount++
					codes[i] += factor
					if value {
						codes[i] += factor
					}
				}
				factor *= 3
			}
			if entryCount != len(row) {
				t.Errorf("AddEdge created an unexpected destination in row %q: %v", key, row)
			}
		}
	}
	if rowCount != len(graph.edges) {
		t.Errorf("AddEdge created an unexpected source: %v", graph.edges)
	}
	return fmt.Sprintf("%d,%d", codes[0], codes[1])
}
