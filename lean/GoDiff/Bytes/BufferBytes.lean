import GoLean.Bytes.BufferBytes
import GoDiff.Math.Bits.Tuple

private def wireBytes (bytes : List (BitVec 8)) : String :=
  String.intercalate "," (bytes.map fun byte => toString byte.toNat)

def main : IO Unit := runTuple fun row => do
  match row with
  | [0] =>
    match GoLean.Bytes.bufferBytes none with
    | none => return "panic"
    | some _ => return "unexpected successful nil receiver"
  | tag :: length :: offset :: raw =>
    let bytes ← raw.mapM (wordInput 8)
    let (stored, spare) := (bytes.take length.toNat, bytes.drop length.toNat)
    let state : GoLean.Bytes.BufferTailState := ⟨if tag == 1 then none else some stored, spare, offset.toNat⟩
    let result := GoLean.Bytes.bufferBytes (some state)
    let view := result.join.getD ⟨[], []⟩
    let mutated := view.bytes.map fun byte => byte ^^^ 1
    let backing := stored.take state.offset ++ mutated ++ view.spare
    let marker := if result.join.isNone then "nil" else "slice"
    return s!"{marker}:{view.bytes.length + view.spare.length}:{wireBytes view.bytes}:" ++
      s!"{wireBytes view.spare}:{wireBytes mutated}:{wireBytes backing}"
  | _ => throw (IO.userError "expected receiver tag, length, offset and raw bytes")
