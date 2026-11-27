import GoLean.Bytes.BufferAvailableBuffer
import GoDiff.Math.Bits.Tuple

private def wireBytes (bytes : List (BitVec 8)) : String :=
  String.intercalate "," (bytes.map fun byte => toString byte.toNat)

def main : IO Unit := runTuple fun row => do
  match row with
  | [0] =>
    match GoLean.Bytes.bufferAvailableBuffer none with
    | none => return "panic"
    | some _ => return "unexpected successful nil receiver"
  | tag :: length :: spareLength :: offset :: raw =>
    let bytes ← raw.mapM (wordInput 8)
    let (stored, tail) := (bytes.take length.toNat, bytes.drop length.toNat)
    let (spare, payload) := (tail.take spareLength.toNat, tail.drop spareLength.toNat)
    let state : GoLean.Bytes.BufferTailState := ⟨if tag == 1 then none else some stored, spare, offset.toNat⟩
    let result := GoLean.Bytes.bufferAvailableBuffer (some state)
    let view := (result.join).getD ⟨[], []⟩
    let unread := (state.stored.getD []).drop state.offset
    let tailAfter := if payload.length ≤ view.spare.length then payload ++ view.spare.drop payload.length
      else view.spare
    let marker := if result.join.isNone then "nil" else "slice"
    return s!"{marker}:{view.bytes.length}:{view.spare.length}:{wireBytes view.spare}:" ++
      s!"{wireBytes tailAfter}:{wireBytes unread}:{payload.length}:{wireBytes (unread ++ payload)}"
  | _ => throw (IO.userError "expected receiver tag, lengths, offset and raw bytes")
