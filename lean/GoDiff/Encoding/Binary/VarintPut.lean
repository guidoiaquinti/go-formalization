import GoLean.Encoding.Binary.VarintPut
import GoDiff.Math.Bits.Tuple

open GoLean.Encoding.Binary

def main : IO Unit := runTuple fun fields => do
  let word :: offset :: size :: raw := fields | throw (IO.userError "expected word and destination view")
  if offset < 0 || size < 0 || offset.toNat + size.toNat > raw.length then
    throw (IO.userError "invalid destination view")
  let word ← wordInput 64 word
  let bytes ← raw.mapM (wordInput 8)
  let backing := bytes.drop offset.toNat
  let results := [putUvarint word size.toNat backing, putVarint word size.toNat backing]
  return String.intercalate "|" (results.map fun (reported, after) =>
    let count := (reported.map toString).getD "bounds"
    let raw := bytes.take offset.toNat ++ after
    s!"{count}:{String.intercalate "," (raw.map (toString ∘ BitVec.toNat))}")
