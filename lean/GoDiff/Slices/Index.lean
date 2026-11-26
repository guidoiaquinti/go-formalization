import GoLean.Slices.Index
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let key :: xs := fields | throw (IO.userError "expected key and slice elements")
  return match GoLean.Slices.sliceIndex (· == ·) xs key with
  | none => "-1"
  | some i => toString i
