namespace GoLean.Internal.Dag

/-- None represents a nil slice or map. Present association lists represent
    allocated maps, with the entries kept in the source key and value types. -/
structure Graph where
  nodes : Option (List String)
  byLabel : Option (List (String × Int))
  edges : Option (List (String × List (String × Bool)))

/-- Go initializes both maps, leaving Nodes at its slice zero value.
    https://github.com/golang/go/blob/go1.27.0/src/internal/dag/parse.go#L51-L59 -/
def newGraph : Graph :=
  { nodes := none, byLabel := some [], edges := some [] }

theorem newGraph_spec :
    newGraph.nodes = none ∧ newGraph.byLabel = some [] ∧ newGraph.edges = some [] :=
  ⟨rfl, rfl, rfl⟩

end GoLean.Internal.Dag
