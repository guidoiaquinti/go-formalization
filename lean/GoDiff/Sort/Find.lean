import GoLean.Sort.Find
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [n, cut, endMatch] := fields | throw (IO.userError "expected length and two thresholds")
  let (i, found) := GoLean.Sort.sortFind n fun j =>
    if Int.ofNat j < cut then 1 else if Int.ofNat j < endMatch then 0 else -1
  return s!"{i} {found}"
