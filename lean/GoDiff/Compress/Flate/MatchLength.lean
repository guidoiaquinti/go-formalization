import GoLean.Compress.Flate.MatchLength
import GoDiff.Math.Bits.Tuple

open GoLean.Compress.Flate

def main : IO Unit := runTuple fun row => do
  let leftLength :: rightLength :: raw := row | throw (IO.userError "expected slice geometry")
  let values ← raw.mapM (wordInput 8)
  let left := values.take leftLength.toNat
  let right := values.drop leftLength.toNat
  return (matchLen left rightLength.toNat right).map toString |>.getD "bounds"
