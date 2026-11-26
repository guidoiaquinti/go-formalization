import GoLean.Slices.IsSorted
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let 0 :: xs := fields | throw (IO.userError "expected marker and integer elements")
  return toString (GoLean.Slices.isSorted (fun a b => decide (a < b)) 0 xs)
