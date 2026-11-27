package dag

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanAddNode(t *testing.T) {
	dir := os.Getenv("GOLEAN_ADDNODE_DIR")
	if dir == "" {
		t.Fatal("GOLEAN_ADDNODE_DIR is missing")
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
	reader, writer := bufio.NewReader(in), bufio.NewWriter(out)
	for {
		var size, mode, query, existingRow int
		if _, err := fmt.Fscan(reader, &size, &mode, &query, &existingRow); err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		graph := goLeanAddNodeGraph(size, mode, existingRow)
		if _, err := fmt.Fprintln(writer, goLeanAddNodeCall(graph, strconv.Itoa(query)), goLeanAddNodeState(graph, size)); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}
}

func goLeanAddNodeGraph(size, mode, existingRow int) *Graph {
	if mode == 4 {
		return nil
	}
	graph := &Graph{byLabel: make(map[string]int), edges: make(map[string]map[string]bool)}
	for i := range size {
		node := strconv.Itoa(i)
		graph.Nodes = append(graph.Nodes, node)
		graph.byLabel[node] = i
	}
	for i := range size + existingRow {
		graph.edges[strconv.Itoa(i)] = map[string]bool{"sentinel": false}
	}
	if mode == 1 || mode == 3 {
		graph.byLabel = nil
	}
	if mode == 2 || mode == 3 {
		graph.edges = nil
	}
	return graph
}

func goLeanAddNodeCall(graph *Graph, query string) (result string) {
	defer func() {
		if recovered := recover(); recovered != nil {
			if _, ok := recovered.(runtime.Error); ok {
				result = "panic"
			} else {
				result = fmt.Sprintf("unexpected panic %v", recovered)
			}
		}
	}()
	return strconv.FormatBool(graph.addNode(query))
}

func goLeanAddNodeState(graph *Graph, size int) string {
	if graph == nil {
		return "nil"
	}
	nodes, ranks, edges := "-", "-", "-"
	if graph.Nodes != nil {
		nodes = strings.Join(graph.Nodes, ",")
	}
	if graph.byLabel != nil {
		out := []string{strconv.Itoa(len(graph.byLabel))}
		for i := range size + 1 {
			value := "-"
			if rank, ok := graph.byLabel[strconv.Itoa(i)]; ok {
				value = strconv.Itoa(rank)
			}
			out = append(out, value)
		}
		ranks = strings.Join(out, ",")
	}
	if graph.edges != nil {
		out := []string{strconv.Itoa(len(graph.edges))}
		for i := range size + 1 {
			row, present := graph.edges[strconv.Itoa(i)]
			value := "?"
			if !present {
				value = "-"
			}
			if v, ok := row["sentinel"]; len(row) == 1 && ok && !v {
				value = "1"
			}
			if row != nil && len(row) == 0 {
				value = "0"
			}
			out = append(out, value)
		}
		edges = strings.Join(out, ",")
	}
	return fmt.Sprintf("nodes=%s rank=%s edges=%s", nodes, ranks, edges)
}
