import GoLean.Slices.Contains
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let key :: xs := fields | throw (IO.userError "expected key and slice elements")
  return toString (GoLean.Slices.sliceContains (· == ·) xs key)
