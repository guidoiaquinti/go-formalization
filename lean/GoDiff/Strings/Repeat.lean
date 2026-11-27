import GoLean.Strings.Repeat
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let maxInt :: count :: raw := row | throw (IO.userError "expected int maximum, count and bytes")
  let input ← raw.mapM (wordInput 8)
  match GoLean.Strings.repeatBytes (input.map BitVec.toNat) count maxInt.toNat with
  | .error .negative => return "negative"
  | .error .overflow => return "overflow"
  | .ok bytes => return "bytes" ++ String.intercalate "," (bytes.map toString)
