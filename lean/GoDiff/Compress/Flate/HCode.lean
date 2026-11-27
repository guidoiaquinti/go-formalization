import GoLean.Compress.Flate.HCode
import GoDiff.Math.Bits.Tuple

open GoLean.Compress.Flate

def main : IO Unit := runTuple fun row => do
  let [raw, code, length] := row | throw (IO.userError "expected original word, code and length")
  let word ← wordInput 32 raw
  let code ← wordInput 16 code
  let length ← wordInput 8 length
  return s!"{(newHCode code length).toNat} {(hCodeLen word).toNat} {(hCodeCode64 word).toNat} {hCodeZero word} {(hCodeSet word code length).toNat}"
