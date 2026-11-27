import GoLean.Slices.Overlaps
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

def main : IO Unit := runTuple fun fields => do
  let [aStart, aLength, bStart, bLength, size] := fields |
    throw (IO.userError "expected both byte starts, both element lengths and the element size")
  return toString (sliceOverlaps aStart.toNat aLength.toNat bStart.toNat bLength.toNat size.toNat)
