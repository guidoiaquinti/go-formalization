import GoLean.Slices.IsSortedFunc
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let length :: responses := fields | throw (IO.userError "expected length and callback responses")
  let (result, calls) := GoLean.Slices.isSortedFunc (List.range length.toNat) 0 fun i _ =>
    match responses[i - 1]? with
    | some 0 => .ok 7
    | some 1 => .ok (-7)
    | some 3 => .ok 0
    | _ => .error ()
  let status := match result with
    | .error () => "panic"
    | .ok value => toString value
  return String.intercalate " " (status :: calls.map toString)
