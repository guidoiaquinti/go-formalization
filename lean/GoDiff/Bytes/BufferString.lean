import GoLean.Bytes.BufferString
import GoDiff.Math.Bits.Tuple

private def wireBytes (bytes : List (BitVec 8)) : String :=
  String.intercalate "," (bytes.map fun byte => toString byte.toNat)

def main : IO Unit := runTuple fun row => do
  match row with
  | [0] => return wireBytes (GoLean.Bytes.bufferString none)
  | tag :: length :: offset :: raw =>
    let bytes ← raw.mapM (wordInput 8)
    let (stored, spare) := (bytes.take length.toNat, bytes.drop length.toNat)
    let state : GoLean.Bytes.BufferTailState := ⟨if tag == 1 then none else some stored, spare, offset.toNat⟩
    let original := GoLean.Bytes.bufferString (some state)
    let overwrite := original.map fun byte => byte ^^^ 1
    let changed := { state with stored := state.stored.map fun all => all.take state.offset ++ overwrite }
    let changedSnapshot := GoLean.Bytes.bufferString (some changed)
    let resetBytes := List.replicate bytes.length (BitVec.ofNat 8 165)
    let reset : GoLean.Bytes.BufferTailState := ⟨some resetBytes, [], 0⟩
    let current := GoLean.Bytes.bufferString (some reset)
    return s!"{wireBytes original}:{wireBytes changedSnapshot}:{wireBytes current}"
  | _ => throw (IO.userError "expected receiver tag, length, offset and raw bytes")
