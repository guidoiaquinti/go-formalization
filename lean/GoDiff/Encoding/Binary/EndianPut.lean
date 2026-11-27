import GoLean.Encoding.Binary.EndianPut
import GoDiff.Math.Bits.Tuple

open GoLean.Encoding.Binary

def main : IO Unit := runTuple fun fields => do
  let word :: offset :: size :: raw := fields | throw (IO.userError "expected word and destination view")
  if offset < 0 || size < 0 || offset.toNat + size.toNat > raw.length then
    throw (IO.userError "invalid destination view")
  let word ← wordInput 64 word
  let bytes ← raw.mapM (wordInput 8)
  let view := (bytes.drop offset.toNat).take size.toNat
  let output := [bigPutUint16 (BitVec.ofNat 16 word.toNat) view,
    bigPutUint32 (BitVec.ofNat 32 word.toNat) view, bigPutUint64 word view,
    littlePutUint16 (BitVec.ofNat 16 word.toNat) view,
    littlePutUint32 (BitVec.ofNat 32 word.toNat) view, littlePutUint64 word view]
  return String.intercalate "|" (output.map fun (success, result) =>
    let storage := bytes.take offset.toNat ++ result ++ bytes.drop (offset.toNat + size.toNat)
    (if success then "ok:" else "panic:") ++ String.intercalate "," (storage.map (toString ∘ BitVec.toNat)))
