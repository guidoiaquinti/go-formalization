import GoLean.Encoding.Binary.EndianAppend
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
  let length := size.toNat
  let output := [bigAppendUint16 (BitVec.ofNat 16 word.toNat) length backing,
    bigAppendUint32 (BitVec.ofNat 32 word.toNat) length backing, bigAppendUint64 word length backing,
    littleAppendUint16 (BitVec.ofNat 16 word.toNat) length backing,
    littleAppendUint32 (BitVec.ofNat 32 word.toNat) length backing, littleAppendUint64 word length backing]
  return String.intercalate "|" (output.map fun result =>
    s!"{result.reused}:{render result.output}:{render (bytes.take offset.toNat ++ result.originalBacking)}")
