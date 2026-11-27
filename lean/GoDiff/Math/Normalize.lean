import GoLean.Math.Normalize
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [raw] := fields | throw (IO.userError "expected a binary64 word")
  let word := (← wordInput 64 raw).toNat
  let (normal, adjustment) := GoLean.Math.normalize word
  return s!"{normal} {adjustment} {GoLean.Math.ilogbFinite word}"
