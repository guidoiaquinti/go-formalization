import GoLean.Slices.EqualFunc
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let rightLength :: responses := fields | throw (IO.userError "expected right length and responses")
  let (result, calls) := GoLean.Slices.sliceEqualFunc (fun x _ =>
    if x = 0 then .ok true else if x = 1 then .ok false else .error ())
    responses (List.range rightLength.toNat)
  let status := match result with
    | .ok value => toString value
    | .error () => "panic"
  return s!"{status} {calls}"
