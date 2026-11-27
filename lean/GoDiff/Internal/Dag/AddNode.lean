import GoLean.Internal.Dag.AddNode
import GoDiff.Math.Bits.Tuple

open GoLean.Internal.Dag

private def render (outcome : Option Bool × Option Graph) (size : Nat) : String :=
  let (result, state) := outcome
  let status := (result.map toString).getD "panic"
  match state with
  | none => status ++ " nil"
  | some graph =>
    let nodes := (graph.nodes.map fun rows => String.intercalate "," rows).getD "-"
    let keys := (List.range (size + 1)).map toString
    let labels := (graph.byLabel.map fun rows =>
      String.intercalate "," (toString rows.length :: keys.map fun key =>
        ((rows.lookup key).map toString).getD "-")).getD "-"
    let edges := (graph.edges.map fun rows =>
      String.intercalate "," (toString rows.length :: keys.map fun key =>
        match rows.lookup key with
        | some [] => "0"
        | some [("sentinel", false)] => "1"
        | none => "-"
        | _ => "?")).getD "-"
    status ++ " nodes=" ++ nodes ++ " rank=" ++ labels ++ " edges=" ++ edges

def main : IO Unit := runTuple fun fields => do
  let [size, mode, query, existingRow] := fields | throw (IO.userError "expected size, mode, query and collision flag")
  let nodes := (List.range size.toNat).map toString
  let label := toString query
  let labels := (List.range size.toNat).map fun (i : Nat) => (toString i, (i : Int))
  let edges := ((List.range (size.toNat + existingRow.toNat)).map fun i => (toString i, [("sentinel", false)]))
  let graph : Graph := {
    nodes := if size == 0 then none else some nodes
    byLabel := if mode == 1 || mode == 3 then none else some labels
    edges := if mode == 2 || mode == 3 then none else some edges }
  return render (addNode (if mode == 4 then none else some graph) label) size.toNat
