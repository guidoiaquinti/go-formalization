package dag

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanEdges(t *testing.T) {
	dir := os.Getenv("GOLEAN_EDGES_DIR")
	if dir == "" {
		t.Fatal("missing GOLEAN_EDGES_DIR")
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
		var graph *Graph
		switch fields[0] {
		case "0":
		case "1":
			graph = &Graph{}
		case "2":
			graph = newGraph()
		case "3":
			graph = newGraph()
			graph.edges["source"] = nil
		case "4":
			graph = newGraph()
			row := make(map[string]bool)
			graph.edges["source"] = row
			for i, field := range fields[1:] {
				code, err := strconv.Atoi(field)
				if err != nil {
					t.Fatal(err)
				}
				key := strconv.Itoa(i)
				graph.byLabel[key] = code / 3
				if status := code % 3; status != 0 {
					row[key] = status == 2
				}
			}
		default:
			t.Fatalf("unknown graph state %q", fields[0])
		}
		if _, err := fmt.Fprintln(out, goLeanEdgesOutcome(graph)); err != nil {
			t.Fatal(err)
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
}

func goLeanEdgesOutcome(graph *Graph) (result string) {
	defer func() {
		if recovered := recover(); recovered != nil {
			if _, ok := recovered.(runtime.Error); !ok {
				panic(recovered)
			}
			result = "panic"
		}
	}()
	keys := graph.Edges("source")
	if keys == nil {
		return "nil"
	}
	return strings.Join(append([]string{"ok"}, keys...), " ")
}
