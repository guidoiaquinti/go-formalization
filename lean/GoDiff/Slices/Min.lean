import GoLean.Slices.Min
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let 0 :: xs := fields | throw (IO.userError "expected marker and elements")
  match GoLean.Slices.sliceMin min xs with
  | .ok value => return toString value
  | .error message => return message
