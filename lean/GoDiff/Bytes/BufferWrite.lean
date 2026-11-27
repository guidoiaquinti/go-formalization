import GoLean.Bytes.BufferWrite
import GoDiff.Math.Bits.Tuple

open GoLean.Bytes

private def wire (bytes : List (BitVec 8)) : String :=
  String.intercalate "," (bytes.map fun byte => toString byte.toNat)

def main : IO Unit := runTuple fun row => do
  let (receiver, source) ← match row with
    | [0] => pure (none, .detached [])
    | tag :: length :: spare :: offset :: lastRead :: aliasStart :: aliasLen :: raw => do
      let bytes ← raw.mapM (wordInput 8)
      let stored := if tag == 1 then none else some (bytes.take length.toNat)
      let tail : BufferTailState := ⟨stored, (bytes.drop length.toNat).take spare.toNat, offset.toNat⟩
      let source := if aliasStart < 0 then BufferWriteInput.detached (bytes.drop (length + spare).toNat)
        else BufferWriteInput.backing aliasStart.toNat aliasLen.toNat
      pure (some (⟨tail, lastRead⟩ : BufferReadState), source)
    | _ => throw (IO.userError "expected receiver geometry and write input")
  match bufferWrite receiver source with
  | none => return "panic"
  | some (count, result) =>
    let state := result.state
    let tail := state.tail
    let marker := if tail.stored.isNone then "nil" else "slice"
    let stored := tail.stored.getD []
    return s!"ok:{count}:{state.lastRead}:{tail.offset}:{stored.length + tail.spare.length}:{marker}:" ++
      s!"{result.reused}:{wire stored}:{wire tail.spare}"
