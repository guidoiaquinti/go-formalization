import GoLean.Compress.Flate.TokenFields
import GoDiff.Math.Bits.Tuple

open GoLean.Compress.Flate

def main : IO Unit := runTuple fun row => do
  let [raw] := row | throw (IO.userError "expected one token word")
  let word ← wordInput 32 raw
  return s!"{tokenType word} {tokenLiteral word} {tokenOffset word} {tokenLength word}"
