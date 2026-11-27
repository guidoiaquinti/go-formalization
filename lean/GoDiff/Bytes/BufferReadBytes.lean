import GoLean.Bytes.BufferReadBytes
import GoDiff.Math.Bits.Tuple

open GoLean.Bytes

private def wire (bytes : List (BitVec 8)) : String :=
  String.intercalate "," (bytes.map (toString ∘ BitVec.toNat))

private def snapshot (state : BufferReadState) : String :=
  s!"{state.tail.offset}:{state.lastRead}:{wire (state.tail.stored.getD [])}:{wire state.tail.spare}"

def main : IO Unit := runTuple fun row => do
  let (receiver, count) ← match row with
    | [0, count] => pure (none, count)
    | tag :: count :: length :: offset :: flag :: raw => do
      let bytes ← raw.mapM (wordInput 8)
      let stored := if tag == 1 then none else some (bytes.take length.toNat)
      let tail : BufferTailState := ⟨stored, bytes.drop length.toNat, offset.toNat⟩
      pure (some (⟨tail, flag⟩ : BufferReadState), count)
    | _ => throw (IO.userError "expected receiver geometry and count")
  match bufferReadBytes receiver (BitVec.ofInt 8 count) with
  | none => return "nil-receiver"
  | some ⟨value, eof, state⟩ =>
    let status := if eof then "eof" else "ok"
    match value with
    | none => return s!"{status}:nil:{snapshot state}"
    | some view =>
      return s!"{status}:slice:{view.bytes.length + view.spare.length}:{wire view.bytes}:" ++
        s!"{wire view.spare}:{snapshot state}"
