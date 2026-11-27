import GoLean.Bytes.ContainsRune
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let r :: input := fields | throw (IO.userError "expected a signed rune and raw input bytes")
  let bytes ← input.mapM (wordInput 8)
  return toString (GoLean.Bytes.containsRune (bytes.map BitVec.toNat) r)
