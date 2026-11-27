import GoLean.Internal.Dag.NewGraph

namespace GoLean.Internal.Dag

/-- None in the first field denotes panic. The second field retains graph
    writes completed before that panic.
    https://github.com/golang/go/blob/go1.27.0/src/internal/dag/parse.go#L51-L68 -/
def addNode (receiver : Option Graph) (label : String) : Option Bool × Option Graph :=
  match receiver with
  | none => (none, none)
  | some graph =>
    match graph.byLabel with
    | none => (none, some graph)
    | some labels =>
      if (labels.lookup label).isSome then (some false, some graph)
      else
        let nodes := graph.nodes.getD []
        let written := { graph with
          nodes := some (nodes ++ [label])
          byLabel := some ((label, (nodes.length : Int)) :: labels) }
        match written.edges with
        | none => (none, some written)
        | some edges =>
          (some true, some { written with
            edges := some ((label, []) :: edges.filter (fun entry => entry.1 != label)) })

theorem addNode_spec (graph : Graph) (label : String) :
    addNode (some graph) label =
      if (graph.byLabel.bind (·.lookup label)).isSome then (some false, some graph)
      else match graph.byLabel with
        | none => (none, some graph)
        | some labels =>
          let afterLabels := { graph with
            nodes := some (graph.nodes.getD [] ++ [label])
            byLabel := some ((label, ((graph.nodes.getD []).length : Int)) :: labels) }
          match graph.edges with
          | none => (none, some afterLabels)
          | some rows => (some true, some { afterLabels with
            edges := some ((label, []) :: rows.filter (fun entry => entry.1 != label)) }) := by
  cases h : graph.byLabel with
  | none => simp [addNode, h]
  | some labels => simp [addNode, h]

theorem addNode_nil (label : String) : addNode none label = (none, none) := rfl

end GoLean.Internal.Dag
