import GoLean.Slices.CompareFunc
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let rightLength :: responses := fields | throw (IO.userError "expected right length and responses")
  let (result, calls) := GoLean.Slices.sliceCompareFunc (fun x _ =>
    if x = 2 then .error () else .ok x)
    responses (List.range rightLength.toNat)
  let status := match result with
    | .ok value => toString value
    | .error () => "panic"
  return s!"{status} {calls}"
