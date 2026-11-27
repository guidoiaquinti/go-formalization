import GoLean.Strings.ContainsAny
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let count :: raw := row | throw (IO.userError "expected input length and bytes")
  let bytes ← raw.mapM fun field => do
    return (← wordInput 8 field).toNat
  return toString (GoLean.Strings.containsAny (bytes.take count.toNat) (bytes.drop count.toNat))
