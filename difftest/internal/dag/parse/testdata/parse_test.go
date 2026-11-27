package dag

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanParse(t *testing.T) {
	data, err := os.ReadFile(os.Getenv("GOLEAN_DAG_PARSE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Text     []byte
		Expected string
	}
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}
	for i, c := range cases {
		got := goLeanParseOutcome(t, string(c.Text))
		if got != c.Expected {
			t.Fatalf("case %d %q: Go=%q want=%q", i, c.Text, got, c.Expected)
		}
	}
}

func goLeanParseOutcome(t *testing.T, source string) string {
	t.Helper()
	graph, err := Parse(source)
	if err != nil {
		if graph != nil {
			t.Fatal("Parse returned a graph on error")
		}
		if _, ok := err.(syntaxError); ok {
			return "syntax " + goLeanParseBytes(err.Error())
		}
		return "semantic " + goLeanParseBytes(err.Error())
	}
	if graph == nil || graph.byLabel == nil || graph.edges == nil {
		t.Fatal("Parse returned an incomplete graph")
	}
	if len(graph.byLabel) != len(graph.Nodes) || len(graph.edges) != len(graph.Nodes) {
		t.Fatal("Parse returned a map with the wrong number of nodes")
	}
	rows := []string{"graph", strconv.Itoa(len(graph.Nodes))}
	for i, label := range graph.Nodes {
		rank, ok := graph.byLabel[label]
		if !ok || rank != i || graph.edges[label] == nil {
			t.Fatalf("bad rank or missing edge row for %q", label)
		}
		rows = append(rows, goLeanParseBytes(label))
	}
	for _, from := range graph.Nodes {
		for to, value := range graph.edges[from] {
			if _, exists := graph.byLabel[to]; !exists || !value {
				t.Fatalf("unexpected edge %q to %q=%t", from, to, value)
			}
		}
		for _, to := range graph.Nodes {
			if graph.edges[from][to] {
				rows = append(rows, "1")
			} else {
				rows = append(rows, "0")
			}
		}
	}
	return strings.Join(rows, " ")
}

func goLeanParseBytes(s string) string {
	var row strings.Builder
	fmt.Fprint(&row, len(s))
	for i := range len(s) {
		fmt.Fprint(&row, " ", s[i])
	}
	return row.String()
}
