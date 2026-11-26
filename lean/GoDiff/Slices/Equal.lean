import GoLean.Slices.Equal
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let n :: values := fields | throw (IO.userError "expected left length and both slices")
  unless 0 ≤ n && n.toNat ≤ values.length do
    throw (IO.userError "invalid left length")
  return toString (GoLean.Slices.sliceEqual (· == ·) (values.take n.toNat) (values.drop n.toNat))
