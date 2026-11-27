import GoLean.Bytes.CopyReplace
import GoDiff.Math.Bits.Tuple

open GoLean.Bytes

private def render (r : ByteCopy) : String :=
  let cap := (r.exactCapacity.map toString).getD "append"
  let tag := if r.bytes.isNone then "nil" else "slice"
  s!"{tag}:{cap}:{String.intercalate "," ((r.bytes.getD []).map toString)}"

def main : IO Unit := runTuple fun row => do
  let tag :: size :: oldSize :: newSize :: limit :: raw := row |
    throw (IO.userError "expected replacement input lengths")
  let values ← raw.mapM (wordInput 8)
  let bytes := values.map BitVec.toNat
  let input := if tag == 0 then none else some (bytes.take size.toNat)
  let old := (bytes.drop size.toNat).take oldSize.toNat
  let replacement := (bytes.drop (size.toNat + oldSize.toNat)).take newSize.toNat
  return String.intercalate "|" [render (cloneBytes input), render (replaceBytes input old replacement limit),
    render (replaceAllBytes input old replacement)]
