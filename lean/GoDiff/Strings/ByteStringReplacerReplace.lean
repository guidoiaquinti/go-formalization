import GoLean.Strings.ByteStringReplacerReplace
import GoDiff.Math.Bits.Tuple

open GoLean.Strings

private def mappingRows : Nat → List Int → IO (List (BitVec 8 × List (BitVec 8)) × List Int)
  | 0, fields => pure ([], fields)
  | count + 1, key :: width :: fields => do
    if width < 0 || width.toNat > fields.length then
      throw (IO.userError "invalid replacement width")
    let byte ← wordInput 8 key
    let value ← (fields.take width.toNat).mapM (wordInput 8)
    let (entries, input) ← mappingRows count (fields.drop width.toNat)
    return ((byte, value) :: entries, input)
  | _, _ => throw (IO.userError "expected mapping key and width")

def main : IO Unit := runTuple fun fields => do
  let marker :: count :: fields := fields | throw (IO.userError "expected mapping and input")
  let (entries, raw) ← mappingRows count.toNat fields
  let input ← raw.mapM (wordInput 8)
  let mapping : ByteStringReplacement :=
    ⟨entries.map Prod.fst, fun key => (entries.find? (fun entry => entry.1 == key)).map Prod.snd⟩
  let some bytes := byteStringReplacerReplace (if marker == 1 then some mapping else none) input
    | return "panic"
  return "bytes:" ++ String.intercalate "," (bytes.map (fun byte => toString byte.toNat))
