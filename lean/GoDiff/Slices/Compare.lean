import GoLean.Slices.Compare
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let n :: values := fields | throw (IO.userError "expected left length and both slices")
  unless 0 ≤ n && n.toNat ≤ values.length do
    throw (IO.userError "invalid left length")
  return toString (GoLean.Slices.sliceCompare GoLean.Slices.compareInt (values.take n.toNat) (values.drop n.toNat))
