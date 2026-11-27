import GoLean.Internal.Dag.TransitiveReduction
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let size :: values := fields | throw (IO.userError "expected a size and adjacency matrix")
  let n := size.toNat
  if size < 0 || values.length != n * n || !values.all (fun x => x == 0 || x == 1) then
    throw (IO.userError "invalid adjacency matrix")
  let nodes := List.range n
  let edges := nodes.flatMap fun source =>
    nodes.filterMap fun dest =>
      if values[source * n + dest]? == some 1 then some (source, dest) else none
  let result := GoLean.Internal.Dag.transitiveReduction nodes edges
  return String.intercalate " " ("edges" :: nodes.flatMap fun source =>
    nodes.map fun dest => if result.contains (source, dest) then "1" else "0")
