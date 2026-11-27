import GoLean.Strings.ByteReplacerWriteString
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

def main : IO Unit := runTuple fun fields => do
  let via :: hasTable :: count :: fields := fields
    | throw (IO.userError "expected writer, table marker and script length")
  let (script, raw) ← replies count.toNat fields
  let raw ← raw.mapM (wordInput 8)
  if hasTable == 1 && raw.length < 256 then
    throw (IO.userError "expected complete mapping")
  let table := (raw.take 256).toArray
  let lookup := if hasTable == 1 then some (fun (byte : BitVec 8) => table[byte.toNat]!) else none
  let input := if hasTable == 1 then raw.drop 256 else raw
  let some result := byteReplacerWriteString lookup input (via == 1) script | return "panic"
  return s!"{result.written}:{result.error}:" ++ String.intercalate ";" (result.calls.map event)
