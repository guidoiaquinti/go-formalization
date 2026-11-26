import GoLean.Sort.IsSorted
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let length :: responses := fields | throw (IO.userError "expected length and callback responses")
  let (result, calls) := GoLean.Sort.isSorted length.toNat fun i _ =>
    match responses[i - 1]? with
    | some 0 => .ok false
    | some 1 => .ok true
    | _ => .error ()
  let status := match result with
    | .error () => "panic"
    | .ok value => toString value
  return String.intercalate " " (status :: calls.map toString)
