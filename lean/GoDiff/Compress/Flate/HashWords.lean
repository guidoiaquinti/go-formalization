import GoLean.Compress.Flate.HashWords
import GoDiff.Math.Bits.Tuple

open GoLean.Compress.Flate

def main : IO Unit := runTuple fun row => do
  let [word, bits, count] := row | throw (IO.userError "expected word, output bits and byte count")
  let word ← wordInput 64 word
  let bits ← wordInput 8 bits
  let count ← wordInput 8 count
  return s!"{(hash4u (word.truncate 32) bits).toNat} {(hashLen word bits count).toNat}"
