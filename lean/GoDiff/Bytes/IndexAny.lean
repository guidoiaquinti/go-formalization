import GoLean.Bytes.IndexAny
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let count :: raw := row | throw (IO.userError "expected input length and bytes")
  let bytes ← raw.mapM (wordInput 8)
  let offset := GoLean.Bytes.indexAny (bytes.take count.toNat) (bytes.drop count.toNat)
  return offset.map toString |>.getD "-1"
