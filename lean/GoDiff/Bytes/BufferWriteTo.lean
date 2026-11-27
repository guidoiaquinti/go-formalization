import GoLean.Bytes.BufferWriteTo
import GoDiff.Math.Bits.Tuple

open GoLean.Bytes

private def wire (xs : List (BitVec 8)) : String :=
  String.intercalate "," (xs.map (toString ∘ BitVec.toNat))

def main : IO Unit := runTuple fun row => do
  let (receiver, accepted, failed) ← match row with
    | [0] => pure (none, 0, false)
    | tag :: accepted :: failed :: length :: offset :: flag :: raw => do
      let bytes ← raw.mapM (wordInput 8)
      let stored := if tag == 1 then none else some (bytes.take length.toNat)
      let tail : BufferTailState := ⟨stored, bytes.drop length.toNat, offset.toNat⟩
      pure (some (⟨tail, flag⟩ : BufferReadState), accepted.toNat, failed == 1)
    | _ => throw (IO.userError "expected receiver geometry and writer reply")
  match bufferWriteTo receiver accepted failed with
  | none => return "nil-receiver"
  | some result =>
    let status := match result.status with
      | .ok => "ok"
      | .shortWrite => "short"
      | .writerError => "writer"
      | .invalidCount => "invalid"
    let call := match result.call with
      | none => "-"
      | some view => wire view.bytes ++ "/" ++ wire view.spare
    let tail := result.state.tail
    let marker := if tail.stored.isNone then "nil" else "slice"
    return s!"{status}:{result.count}:{call}:{marker}:{tail.offset}:{result.state.lastRead}:" ++
      s!"{wire (tail.stored.getD [])}:{wire tail.spare}"
