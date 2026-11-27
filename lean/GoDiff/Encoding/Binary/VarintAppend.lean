import GoLean.Encoding.Binary.VarintAppend
import GoDiff.Math.Bits.Tuple

open GoLean.Encoding.Binary

private def render (bytes : List (BitVec 8)) : String :=
  String.intercalate "," (bytes.map (toString ∘ BitVec.toNat))

def main : IO Unit := runTuple fun fields => do
  let word :: offset :: size :: raw := fields | throw (IO.userError "expected word and source view")
  if offset < 0 || size < 0 || offset.toNat + size.toNat > raw.length then
    throw (IO.userError "invalid source view")
  let word ← wordInput 64 word
  let bytes ← raw.mapM (wordInput 8)
  let backing := bytes.drop offset.toNat
  let results := [appendUvarint word size.toNat backing, appendVarint word size.toNat backing]
  return String.intercalate "|" (results.map fun result =>
    s!"{result.reused}:{render result.output}:{render (bytes.take offset.toNat ++ result.originalBacking)}")
