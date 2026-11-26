import GoLean.Sort.IntsAreSorted
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let 0 :: xs := fields | throw (IO.userError "expected marker and slice elements")
  return toString (GoLean.Sort.intsAreSorted xs)
