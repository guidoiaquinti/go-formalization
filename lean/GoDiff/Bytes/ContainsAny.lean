import GoLean.Bytes.ContainsAny
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let count :: raw := row | throw (IO.userError "expected input length and bytes")
  let bytes ← raw.mapM (wordInput 8)
  return toString (GoLean.Bytes.containsAny (bytes.take count.toNat) (bytes.drop count.toNat))
