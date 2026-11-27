import GoLean.Bytes.BufferReset
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
  match GoLean.Bytes.bufferReset receiver with
  | none => return "panic"
  | some state =>
    let tail := state.tail
    let marker := if tail.stored.isNone then "nil" else "slice"
    let length := (tail.stored.getD []).length - tail.offset
    let backing := (tail.stored.getD [] ++ tail.spare).drop tail.offset
    let reused := if backing.isEmpty then [] else 0xa5 :: backing.drop 1
    return s!"{marker}:{length}:{backing.length}:{wireBytes backing}:" ++
      s!"{state.lastRead == 0}:{decide (state.lastRead ≤ 0)}:{wireBytes reused}"
