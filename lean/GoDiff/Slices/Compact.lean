import GoLean.Slices.Compact
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let 0 :: xs := fields | throw (IO.userError "expected marker and elements")
  let (kept, backing) := GoLean.Slices.sliceCompact (· == ·) (0 : Int) xs
  return String.intercalate " " (toString kept.length :: backing.map toString)
