import GoLean.Bytes.BufferNext
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
  match GoLean.Bytes.bufferNext receiver count with
  | none => return "nil-receiver"
  | some result =>
    let state := result.state
    let tail := state.tail
    let view := result.view.getD ⟨[], []⟩
    let status := if result.boundsPanic then "bounds" else "ok"
    let marker := if result.view.isNone then "nil" else "slice"
    let stored := tail.stored.getD []
    let unread := stored.drop tail.offset
    let mutated := view.bytes.map fun byte => byte ^^^ 1
    let backing := stored.take (tail.offset - view.bytes.length) ++ mutated ++ unread ++ tail.spare
    let unreturned := if state.lastRead == 0 then "rejected"
      else (mutated.getLast?.map fun byte => toString byte.toNat).getD "missing"
    return s!"{status}:{marker}:{view.bytes.length + view.spare.length}:{wireBytes view.bytes}:" ++
      s!"{tail.offset}:{wireBytes unread}:{wireBytes backing}:{unreturned}:{decide (state.lastRead ≤ 0)}"
