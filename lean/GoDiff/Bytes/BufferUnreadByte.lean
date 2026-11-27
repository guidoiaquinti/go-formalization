import GoLean.Bytes.BufferUnreadByte
import GoDiff.Math.Bits.Tuple

open GoLean.Bytes

private def wire (bytes : List (BitVec 8)) : String :=
  String.intercalate "," (bytes.map fun byte => toString byte.toNat)

def main : IO Unit := runTuple fun row => do
  let receiver ← match row with
    | [0] => pure none
    | tag :: length :: offset :: lastRead :: raw => do
      let bytes ← raw.mapM (wordInput 8)
      let stored := if tag == 1 then none else some (bytes.take length.toNat)
      let tail : BufferTailState := ⟨stored, bytes.drop length.toNat, offset.toNat⟩
      pure (some (⟨tail, lastRead⟩ : BufferReadState))
    | _ => throw (IO.userError "expected receiver tag and buffer state")
  match bufferUnreadByte receiver with
  | none => return "panic"
  | some (failed, state) =>
    let tail := state.tail
    let marker := if tail.stored.isNone then "nil" else "slice"
    let status := if failed then "sentinel-error" else "ok"
    let unread := (tail.stored.getD []).drop tail.offset
    return s!"{status}:{state.lastRead}:{tail.offset}:{unread.length + tail.spare.length}:{marker}:" ++
      s!"{wire unread}:{wire tail.spare}:{wire (tail.stored.getD [])}"
