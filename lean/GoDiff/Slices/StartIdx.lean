import GoLean.Slices.StartIdx
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

def main : IO Unit := runTuple fun fields => do
  let [a, len, b, needleSize] := fields |
    throw (IO.userError "expected haystack start and length, needle start and length")
  let addresses := (List.range len.toNat).map fun offset => a.toNat + offset
  match sliceStartIdx addresses (if needleSize = 0 then none else some b.toNat) with
  | .ok i => return s!"index {i}"
  | .error .emptyNeedle => return "bounds panic"
  | .error .needleNotFound => return "panic needle not found"
