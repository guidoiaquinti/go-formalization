import GoLean.Bytes.BufferTruncate
import GoDiff.Math.Bits.Tuple

private def wireBytes (bytes : List (BitVec 8)) : String :=
  String.intercalate "," (bytes.map fun byte => toString byte.toNat)

def main : IO Unit := runTuple fun row => do
  let (receiver, count) ← match row with
    | [0, count] => pure (none, count)
    | tag :: length :: offset :: lastRead :: count :: raw => do
      let bytes ← raw.mapM (wordInput 8)
      let stored := if tag == 1 then none else some (bytes.take length.toNat)
      let tail : GoLean.Bytes.BufferTailState := ⟨stored, bytes.drop length.toNat, offset.toNat⟩
      pure (some (⟨tail, lastRead⟩ : GoLean.Bytes.BufferReadState), count)
    | _ => throw (IO.userError "expected receiver tag, buffer state and count")
  match GoLean.Bytes.bufferTruncate receiver count with
  | none => return "nil-receiver"
  | some (valid, state) =>
    let tail := state.tail
    let status := if valid then "ok" else "bounds"
    let marker := if tail.stored.isNone then "nil" else "slice"
    let stored := tail.stored.getD []
    let unread := stored.drop tail.offset
    let capacity := unread.length + tail.spare.length
    return s!"{status}:{marker}:{tail.offset}:{capacity}:{wireBytes unread}:" ++
      s!"{wireBytes tail.spare}:{state.lastRead == 0}:{decide (state.lastRead ≤ 0)}"
