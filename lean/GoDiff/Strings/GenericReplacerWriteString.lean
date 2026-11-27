import GoLean.Strings.GenericReplacerWriteString
import GoDiff.Math.Bits.Tuple

open GoLean.Strings

private def replies : Nat → List Int → IO (List ReplacerWriteReply × List Int)
  | 0, fields => pure ([], fields)
  | count + 1, limit :: err :: fields => do
    let (script, raw) ← replies count fields
    return (⟨if limit < 0 then none else some limit.toNat, err.toNat⟩ :: script, raw)
  | _, _ => throw (IO.userError "expected writer script")

private def pairs : Nat → List Int →
    IO (List (List (BitVec 8) × List (BitVec 8)) × List Int)
  | 0, fields => pure ([], fields)
  | count + 1, keyLength :: valueLength :: fields => do
    if keyLength < 0 || valueLength < 0 || keyLength.toNat + valueLength.toNat > fields.length then
      throw (IO.userError "invalid key or value length")
    let key ← (fields.take keyLength.toNat).mapM (wordInput 8)
    let fields := fields.drop keyLength.toNat
    let value ← (fields.take valueLength.toNat).mapM (wordInput 8)
    let (entries, raw) ← pairs count (fields.drop valueLength.toNat)
    return ((key, value) :: entries, raw)
  | _, _ => throw (IO.userError "expected key and value lengths")

private def event (call : ReplacerWriteEvent) : String :=
  (if call.stringWriter then "S" else "B") ++
    String.join (call.bytes.map fun byte => "," ++ toString byte.toNat)

def main : IO Unit := runTuple fun fields => do
  let via :: replyCount :: pairCount :: fields := fields
    | throw (IO.userError "expected writer mode and counts")
  let (script, raw) ← replies replyCount.toNat fields
  let (entries, raw) ← pairs pairCount.toNat raw
  let input ← raw.mapM (wordInput 8)
  let result := genericReplacerWriteString entries input (via == 1) script
  return s!"{result.written}:{result.error}:" ++ String.intercalate ";" (result.calls.map event)
