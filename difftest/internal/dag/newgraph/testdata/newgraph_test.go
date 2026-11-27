package dag

import "testing"

func TestGoLeanNewGraph(t *testing.T) {
	graph := newGraph()
	if graph == nil {
		t.Fatal("newGraph returned a nil pointer")
	}
	if graph.Nodes != nil || len(graph.Nodes) != 0 {
		t.Errorf("Nodes = %#v, want nil and empty", graph.Nodes)
	}
	if graph.byLabel == nil || len(graph.byLabel) != 0 {
		t.Errorf("byLabel = %#v, want allocated and empty", graph.byLabel)
	}
	if graph.edges == nil || len(graph.edges) != 0 {
		t.Errorf("edges = %#v, want allocated and empty", graph.edges)
	}
}
