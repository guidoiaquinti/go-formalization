import GoLean.Slices.ContainsFunc
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let 0 :: xs := fields | throw (IO.userError "expected marker and callback responses")
  let (result, calls) := GoLean.Slices.sliceContainsFunc (fun x =>
    if x = 0 then .ok false else if x = 1 then .ok true else .error ()) xs
  let status := match result with
    | .ok value => toString value
    | .error () => "panic"
  return s!"{status} {calls}"
