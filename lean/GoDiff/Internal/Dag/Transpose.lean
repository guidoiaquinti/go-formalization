import GoLean.Internal.Dag.Transpose
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let size :: codes := fields | throw (IO.userError "expected size followed by matrix")
  let nodes := List.range size.toNat
  let edges := nodes.flatMap fun src => nodes.filterMap fun dest =>
    let code := codes[src * size.toNat + dest]?.getD 0
    if code == 0 then none else some (src, dest, code == 2)
  let rows := GoLean.Internal.Dag.transpose nodes edges
  let bits := rows.flatMap fun (_, row) =>
    nodes.map fun dest => if row.lookup dest == some true then "1" else "0"
  return String.intercalate " " ("ok" :: bits)
