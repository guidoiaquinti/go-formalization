import GoLean.Internal.Dag.NewGraph

private def field (entries : Option (List α)) : String :=
  match entries with
  | none => "nil 0"
  | some xs => s!"present {xs.length}"

def main : IO Unit := do
  let graph := GoLean.Internal.Dag.newGraph
  IO.println s!"{field graph.nodes} {field graph.byLabel} {field graph.edges}"
