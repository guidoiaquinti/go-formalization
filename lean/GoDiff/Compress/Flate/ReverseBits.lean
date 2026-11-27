import GoLean.Compress.Flate.ReverseBits
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let [word, count] := row | throw (IO.userError "expected word and count")
  let word ← wordInput 16 word
  let count ← wordInput 8 count
  return toString (GoLean.Compress.Flate.reverseBits word count).toNat
