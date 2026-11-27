import GoLean.Compress.Flate.MatchCodes
import GoDiff.Math.Bits.Tuple

open GoLean.Compress.Flate

def main : IO Unit := runTuple fun row => do
  let [offset, length] := row | throw (IO.userError "expected offset and length")
  let word ← wordInput 32 offset
  let lengthWord ← wordInput 8 length
  return s!"{offsetCode word.toNat} {lengthCode lengthWord}"
