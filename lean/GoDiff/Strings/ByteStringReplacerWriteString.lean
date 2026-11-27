import GoLean.Strings.ByteStringReplacerWriteString
import GoDiff.Math.Bits.Tuple

open GoLean.Strings

private def replies : Nat → List Int → IO (List ReplacerWriteReply × List Int)
  | 0, fields => pure ([], fields)
  | n + 1, limit :: err :: fields => do
    let (script, raw) ← replies n fields
    return (⟨if limit < 0 then none else some limit.toNat, err.toNat⟩ :: script, raw)
  | _, _ => throw (IO.userError "expected writer script")

private def event (call : ReplacerWriteEvent) : String :=
  (if call.stringWriter then "S" else "B") ++
    String.join (call.bytes.map (fun byte => "," ++ toString byte.toNat))

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
  let via :: hasTable :: count :: entryCount :: fields := fields
    | throw (IO.userError "expected writer, table marker and reply and entry counts")
  let (script, raw) ← replies count.toNat fields
  let (entries, raw) ← mappingRows entryCount.toNat raw
  let input ← raw.mapM (wordInput 8)
  let lookup := if hasTable == 1 then
    some (fun (byte : BitVec 8) => (entries.find? (fun entry => entry.1 == byte)).map Prod.snd)
    else none
  let some result := byteStringReplacerWriteString lookup input (via == 1) script | return "panic"
  return s!"{result.written}:{result.error}:" ++ String.intercalate ";" (result.calls.map event)
