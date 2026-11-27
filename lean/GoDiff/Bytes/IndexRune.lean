import GoLean.Bytes.IndexRune
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let r :: input := fields | throw (IO.userError "expected a signed rune and raw input bytes")
  let bytes ← input.mapM (wordInput 8)
  match GoLean.Bytes.indexRune (bytes.map BitVec.toNat) r with
  | none => return "-1"
  | some offset => return toString offset
