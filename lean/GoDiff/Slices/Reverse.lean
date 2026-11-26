import GoLean.Slices.Reverse
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let 0 :: values := fields | throw (IO.userError "expected marker and slice elements")
  return String.intercalate " " ("ok" :: (GoLean.Slices.sliceReverse values).map toString)
