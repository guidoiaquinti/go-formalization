import GoLean.Strings.LastIndexAny
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let count :: raw := row | throw (IO.userError "expected input length and bytes")
  let bytes ← raw.mapM fun field => do
    return (← wordInput 8 field).toNat
  let offset := GoLean.Strings.lastIndexAny (bytes.take count.toNat) (bytes.drop count.toNat)
  return offset.map toString |>.getD "-1"
