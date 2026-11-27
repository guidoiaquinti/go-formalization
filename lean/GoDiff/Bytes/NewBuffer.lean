import GoLean.Bytes.NewBuffer
import GoDiff.Math.Bits.Tuple

private def wire (xs : List (BitVec 8)) : String :=
  String.intercalate "," (xs.map (toString ∘ BitVec.toNat))

def main : IO Unit := runTuple fun row => do
  let tag :: length :: raw := row | throw (IO.userError "expected slice marker and bytes")
  let payload ← raw.mapM (wordInput 8)
  let input := if tag == 0 then none else some (payload.take length.toNat)
  let result := GoLean.Bytes.newBuffer input (payload.drop length.toNat)
  let stored := result.tail.stored.getD []
  let marker := if result.tail.stored.isNone then "nil" else "slice"
  let capacity := stored.length + result.tail.spare.length
  return s!"{marker}:{stored.length}:{capacity}:{result.tail.offset}:{result.lastRead}:" ++
    s!"{wire stored}:{wire result.tail.spare}"
