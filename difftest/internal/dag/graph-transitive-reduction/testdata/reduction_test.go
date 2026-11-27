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

func TestGoLeanReduction(t *testing.T) {
	dir := os.Getenv("GOLEAN_REDUCTION_DIR")
	if dir == "" {
		t.Fatal("missing GOLEAN_REDUCTION_DIR")
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
		for from := range n {
			source := strconv.Itoa(from)
			for to := range n {
				status := fields[1+from*n+to]
				switch status {
				case "1":
					graph.AddEdge(source, strconv.Itoa(to))
				case "0":
				default:
					t.Fatalf("invalid bit: %q", status)
				}
			}
		}

		graph.TransitiveReduction()
		if graph.edges == nil || len(graph.edges) != n {
			t.Fatalf("wrong outer map shape at %v", fields)
		}
		if len(graph.Nodes) != n || len(graph.byLabel) != n {
			t.Fatalf("node list or rank changed at %v", fields)
		}
		for i, label := range graph.Nodes {
			if label != strconv.Itoa(i) || graph.byLabel[label] != i {
				t.Fatalf("changed node or rank at index %d: %v", i, fields)
			}
		}
		result := []string{"edges"}
		for from := range n {
			row := graph.edges[strconv.Itoa(from)]
			if row == nil {
				t.Fatalf("nil destination row %d at %v", from, fields)
			}
			count := 0
			for to := range n {
				value, present := row[strconv.Itoa(to)]
				switch {
				case !present:
					result = append(result, "0")
				case value:
					result = append(result, "1")
					count++
				default:
					t.Fatalf("false value after reduction at %v", fields)
				}
			}
			if count != len(row) {
				t.Fatalf("unexpected destination key at %v", fields)
			}
		}
		if _, err := fmt.Fprintln(out, strings.Join(result, " ")); err != nil {
			t.Fatal(err)
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
}
