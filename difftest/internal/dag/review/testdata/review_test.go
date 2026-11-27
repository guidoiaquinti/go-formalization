package dag

import (
	"maps"
	"slices"
	"testing"
)

func checkReview(t *testing.T, real, mutant bool) {
	t.Helper()
	if !real || mutant {
		t.Fatalf("Go accepted=%t, mutant accepted=%t, want true and false", real, mutant)
	}
	t.Log("Go passed and the mutant failed")
}

func reviewPanic(f func()) (reason any) {
	defer func() { reason = recover() }()
	f()
	return
}

func TestGoLeanDAGReview(t *testing.T) {
	t.Run("newGraph", func(t *testing.T) {
		spec := func(g *Graph) bool {
			return g != nil && g.Nodes == nil && g.byLabel != nil && len(g.byLabel) == 0 &&
				g.edges != nil && len(g.edges) == 0
		}
		checkReview(t, spec(newGraph()), spec(&Graph{}))
	})
	t.Run("addNode", func(t *testing.T) {
		g := newGraph()
		first, second, duplicate := g.addNode("a"), g.addNode("b"), g.addNode("a")
		spec := func(rank int) bool {
			return first && second && !duplicate && slices.Equal(g.Nodes, []string{"a", "b"}) &&
				rank == 1 && g.edges["b"] != nil && len(g.edges["b"]) == 0
		}
		partial := &Graph{byLabel: map[string]int{}}
		if reviewPanic(func() { partial.addNode("b") }) == nil ||
			!slices.Equal(partial.Nodes, []string{"b"}) || len(partial.byLabel) != 1 || partial.byLabel["b"] != 0 {
			t.Fatal("addNode discarded the state written before its nil-map panic")
		}
		checkReview(t, spec(g.byLabel["b"]), spec(len(g.Nodes)))
	})
	t.Run("AddEdge", func(t *testing.T) {
		g := &Graph{edges: map[string]map[string]bool{"a": {"b": false, "c": false}}}
		g.AddEdge("a", "b")
		want := map[string]bool{"b": true, "c": false}
		checkReview(t, maps.Equal(g.edges["a"], want) && reviewPanic(func() { g.AddEdge("missing", "b") }) != nil,
			maps.Equal(map[string]bool{"b": false, "c": false}, want))
	})
	t.Run("DelEdge", func(t *testing.T) {
		g := &Graph{edges: map[string]map[string]bool{"a": {"b": true, "c": true}}}
		g.DelEdge("a", "b")
		want := map[string]bool{"c": true}
		checkReview(t, maps.Equal(g.edges["a"], want), maps.Equal(map[string]bool{"b": false, "c": true}, want))
	})
	t.Run("HasEdge", func(t *testing.T) {
		g := &Graph{edges: map[string]map[string]bool{"a": {"b": false, "c": true}}}
		_, mutant := g.edges["a"]["b"]
		checkReview(t, !g.HasEdge("a", "b") && g.HasEdge("a", "c"), !mutant)
	})
	t.Run("Edges", func(t *testing.T) {
		g := &Graph{byLabel: map[string]int{"z": 0, "a": 1}, edges: map[string]map[string]bool{"root": {"z": false, "a": true}}}
		want := []string{"z", "a"}
		checkReview(t, slices.Equal(g.Edges("root"), want), slices.Equal([]string{"a"}, want))
	})
	t.Run("Topo", func(t *testing.T) {
		g := &Graph{Nodes: []string{"a", "b", "c"}, byLabel: map[string]int{"a": 0, "b": 1, "c": 2},
			edges: map[string]map[string]bool{"a": {"b": false, "c": true}}}
		want := []string{"a", "c", "b"}
		checkReview(t, slices.Equal(g.Topo(), want), slices.Equal([]string{"b", "a", "c"}, want))
	})
	t.Run("Transpose", func(t *testing.T) {
		g := &Graph{Nodes: []string{"a", "b"}, edges: map[string]map[string]bool{"a": {"b": false}}}
		g.Transpose()
		want := map[string]bool{"a": true}
		checkReview(t, len(g.edges) == 2 && g.edges["a"] != nil && len(g.edges["a"]) == 0 &&
			maps.Equal(g.edges["b"], want), maps.Equal(map[string]bool{"a": false}, want))
	})
	t.Run("TransitiveReduction", func(t *testing.T) {
		g := &Graph{Nodes: []string{"a", "b", "c"}, edges: map[string]map[string]bool{
			"a": {"b": true, "c": true}, "b": {"c": true}, "c": {},
		}}
		g.TransitiveReduction()
		want := map[string]bool{"b": true}
		checkReview(t, maps.Equal(g.edges["a"], want) && maps.Equal(g.edges["b"], map[string]bool{"c": true}),
			maps.Equal(map[string]bool{"b": true, "c": true}, want))
	})
	t.Run("Parse", func(t *testing.T) {
		g, err := Parse("NONE < a; a < b; b < c;")
		order := []string{"a", "b", "c"}
		rejected, syntax := Parse("!")
		_, typed := syntax.(syntaxError)
		if rejected != nil || !typed || syntax.Error() != "parsing graph: line 1: syntax error: unexpected token ! near " {
			t.Fatalf("Parse syntax path returned graph=%v, error=%v", rejected, syntax)
		}
		empty, errEmpty := Parse("")
		single, errSingle := Parse("NONE < single;")
		if errEmpty != nil || empty == nil || len(empty.Nodes) != 0 ||
			errSingle != nil || single == nil || !slices.Equal(single.Nodes, []string{"single"}) || len(single.Edges("single")) != 0 {
			t.Fatal("Parse returned the wrong empty or singleton graph")
		}
		spec := func(result *Graph) bool {
			return result != nil && slices.Equal(result.Nodes, order) && result.HasEdge("c", "a")
		}
		mutant := &Graph{Nodes: order, edges: map[string]map[string]bool{"c": {"b": true}, "b": {"a": true}}}
		checkReview(t, err == nil && spec(g), spec(mutant))
	})
	t.Run("parseRules", func(t *testing.T) {
		out, err := parseRules("a < b c")
		_, typed := err.(syntaxError)
		spec := func(rules []rule) bool {
			return typed && err.Error() == "parsing graph: line 1: syntax error: missing < near c" &&
				len(rules) == 1 && slices.Equal(rules[0].less, []string{"a"}) &&
				rules[0].op == "<" && slices.Equal(rules[0].def, []string{"b"})
		}
		checkReview(t, spec(out), spec(nil))
	})
	t.Run("nextList", func(t *testing.T) {
		p := &rulesParser{lineno: 1, text: "a,b<c"}
		names, tok := p.nextList()
		spec := func(list []string) bool {
			return slices.Equal(list, []string{"a", "b"}) && tok == "<" && p.text == "c" && p.lastWord == "b"
		}
		checkReview(t, spec(names), spec([]string{"b", "a"}))
	})
	t.Run("nextToken", func(t *testing.T) {
		p := &rulesParser{lineno: 1, lastWord: "old", text: "# skip\n !oops"}
		got := reviewPanic(func() { p.nextToken() })
		want := syntaxError("parsing graph: line 2: syntax error: unexpected token ! near old")
		spec := func(line int) bool {
			return got == want && line == 2 && p.text == "!oops" && p.lastWord == "old"
		}
		checkReview(t, spec(p.lineno), spec(1))
	})
	t.Run("syntaxError", func(t *testing.T) {
		p := &rulesParser{lineno: -5, lastWord: "\x00\xff"}
		got := reviewPanic(func() { p.syntaxError("bad\x00") })
		want := syntaxError("parsing graph: line -5: syntax error: bad\x00 near \x00\xff")
		checkReview(t, got == want, got == string(want))
	})
	t.Run("Error", func(t *testing.T) {
		const raw = "\x00\xff"
		want := func(got string) bool { return got == raw }
		checkReview(t, want(syntaxError(raw).Error()), want("\x00\ufffd"))
	})
}
