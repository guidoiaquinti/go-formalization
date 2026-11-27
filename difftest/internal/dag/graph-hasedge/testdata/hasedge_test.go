package dag

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
)

func TestGoLeanHasEdge(t *testing.T) {
	dir := os.Getenv("GOLEAN_HASEDGE_DIR")
	if dir == "" {
		t.Fatal("missing GOLEAN_HASEDGE_DIR")
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
		if _, err := fmt.Fprintln(out, goLeanHasEdgeOutcome(graph, keys[from], keys[to])); err != nil {
			t.Fatal(err)
		}
	}
}

func goLeanHasEdgeOutcome(graph *Graph, from, to string) (result string) {
	defer func() {
		if recovered := recover(); recovered != nil {
			if _, ok := recovered.(runtime.Error); !ok {
				panic(recovered)
			}
			result = "panic"
		}
	}()
	return strconv.FormatBool(graph.HasEdge(from, to))
}
