import GoLean.Bytes.NewBufferString
import GoDiff.Math.Bits.Tuple

private def wire (xs : List (BitVec 8)) : String :=
  String.intercalate "," (xs.map (toString ∘ BitVec.toNat))

def main : IO Unit := runTuple fun row => do
  let _ :: raw := row | throw (IO.userError "expected row marker")
  let bytes ← raw.mapM (wordInput 8)
  let result := GoLean.Bytes.newBufferString bytes
  let stored := result.tail.stored.getD []
  let marker := if result.tail.stored.isNone then "nil" else "slice"
  let capacity := stored.length + result.tail.spare.length
  return s!"{marker}:{stored.length}:{capacity}:{result.tail.offset}:{result.lastRead}:" ++
    s!"{wire stored}:{wire result.tail.spare}"
