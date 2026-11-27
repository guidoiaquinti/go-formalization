import GoLean.Strings.ToValidUTF8
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let count :: bytes := fields | throw (IO.userError "expected replacement length and bytes")
  let raw ← bytes.mapM (wordInput 8)
  let raw := raw.map BitVec.toNat
  let output := GoLean.Strings.toValidUTF8 (raw.drop count.toNat) (raw.take count.toNat)
  return "bytes:" ++ String.intercalate "," (output.map toString)
