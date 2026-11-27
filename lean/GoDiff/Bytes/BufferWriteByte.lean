import GoLean.Bytes.BufferWriteByte
import GoDiff.Math.Bits.Tuple

open GoLean.Bytes

private def wire (bytes : List (BitVec 8)) : String :=
  String.intercalate "," (bytes.map fun byte => toString byte.toNat)

def main : IO Unit := runTuple fun row => do
  let (receiver, byte) ← match row with
    | [0, byte] => pure (none, ← wordInput 8 byte)
    | tag :: byte :: length :: offset :: lastRead :: raw => do
      let bytes ← raw.mapM (wordInput 8)
      let stored := if tag == 1 then none else some (bytes.take length.toNat)
      let tail : BufferTailState := ⟨stored, bytes.drop length.toNat, offset.toNat⟩
      pure (some (⟨tail, lastRead⟩ : BufferReadState), ← wordInput 8 byte)
    | _ => throw (IO.userError "expected receiver tag, written byte and buffer state")
  match bufferWriteByte receiver byte with
  | none => return "panic"
  | some result =>
    let state := result.state
    let tail := state.tail
    let marker := if tail.stored.isNone then "nil" else "slice"
    let stored := tail.stored.getD []
    return s!"ok:{state.lastRead}:{tail.offset}:{stored.length + tail.spare.length}:{marker}:" ++
      s!"{result.reused}:{wire stored}:{wire tail.spare}"
