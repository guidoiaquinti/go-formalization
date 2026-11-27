import GoLean.Strings.SingleStringReplacerWriteString
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
  let via :: count :: patternWidth :: valueWidth :: fields := fields
    | throw (IO.userError "expected writer, script length and pattern and value widths")
  let (script, raw) ← replies count.toNat fields
  let raw ← raw.mapM (wordInput 8)
  if patternWidth ≤ 0 || valueWidth < 0 || patternWidth.toNat + valueWidth.toNat > raw.length then
    throw (IO.userError "invalid pattern or value width")
  let pattern := raw.take patternWidth.toNat
  let rest := raw.drop patternWidth.toNat
  let value := rest.take valueWidth.toNat
  let input := rest.drop valueWidth.toNat
  let result := singleStringReplacerWriteString pattern value input (via == 1) script
  return s!"{result.written}:{result.error}:" ++ String.intercalate ";" (result.calls.map event)
