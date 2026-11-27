import GoLean.Internal.Dag.Edges
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let state :: codes := fields
    | throw (IO.userError "expected graph state and rank/edge codes")
  let entries := codes.zipIdx
  let rank (key : String) : Int :=
    match key.toNat?.bind (fun i => codes[i]?) with
    | some code => code / 3
    | none => 0
  let row := entries.filterMap fun (code, key) =>
    if code % 3 == 0 then none else some (toString key, code % 3 == 2)
  let graph := match state with
    | 0 => none
    | 1 => some none
    | 2 => some (some [])
    | 3 => some (some [("source", none)])
    | _ => some (some [("source", some row)])
  return match GoLean.Internal.Dag.edges graph rank "source" with
    | none => "panic"
    | some keys => String.intercalate " " ("ok" :: keys)
