import GoLean.Bytes.BufferReadByte
import GoDiff.Math.Bits.Tuple

private def wireBytes (bytes : List (BitVec 8)) : String :=
  String.intercalate "," (bytes.map fun byte => toString byte.toNat)

def main : IO Unit := runTuple fun row => do
  let receiver ← match row with
    | [0] => pure none
    | tag :: length :: offset :: lastRead :: raw => do
      let bytes ← raw.mapM (wordInput 8)
      let stored := if tag == 1 then none else some (bytes.take length.toNat)
      let tail : GoLean.Bytes.BufferTailState := ⟨stored, bytes.drop length.toNat, offset.toNat⟩
      pure (some (⟨tail, lastRead⟩ : GoLean.Bytes.BufferReadState))
    | _ => throw (IO.userError "expected receiver tag and buffer state")
  match GoLean.Bytes.bufferReadByte receiver with
  | none => return "nil-receiver"
  | some result =>
    let state := result.state
    let tail := state.tail
    let status := if result.eof then "eof" else "ok"
    let marker := if tail.stored.isNone then "nil" else "slice"
    let stored := tail.stored.getD []
    let unread := stored.drop tail.offset
    let capacity := unread.length + tail.spare.length
    let unreturned := if state.lastRead == 0 then "rejected"
      else (stored[tail.offset - 1]?.map fun byte => toString byte.toNat).getD "missing"
    return s!"{status}:{result.byte.toNat}:{marker}:{tail.offset}:{capacity}:{wireBytes unread}:" ++
      s!"{wireBytes tail.spare}:{unreturned}:{decide (state.lastRead ≤ 0)}"
