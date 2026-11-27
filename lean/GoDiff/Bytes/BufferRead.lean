import GoLean.Bytes.BufferRead
import GoDiff.Math.Bits.Tuple

open GoLean.Bytes

private def wire (xs : List (BitVec 8)) : String :=
  String.intercalate "," (xs.map (toString ∘ BitVec.toNat))

def main : IO Unit := runTuple fun row => do
  let (receiver, destination) ← match row with
    | [0] => pure (none, [])
    | tag :: length :: spare :: offset :: flag :: raw => do
      let bytes ← raw.mapM (wordInput 8)
      let stored := if tag == 1 then none else some (bytes.take length.toNat)
      let tail : BufferTailState := ⟨stored, (bytes.drop length.toNat).take spare.toNat, offset.toNat⟩
      pure (some (⟨tail, flag⟩ : BufferReadState), bytes.drop (length + spare).toNat)
    | _ => throw (IO.userError "expected receiver geometry and destination")
  match bufferRead receiver destination with
  | none => return "panic"
  | some result =>
    let tail := result.state.tail
    let status := if result.eof then "eof" else "ok"
    let marker := if tail.stored.isNone then "nil" else "slice"
    return s!"{status}:{result.count}:{marker}:{tail.offset}:{result.state.lastRead}:" ++
      s!"{wire result.destination}:{wire (tail.stored.getD [])}:{wire tail.spare}"
