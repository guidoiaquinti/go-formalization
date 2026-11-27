import GoLean.Strings.Fields
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let 0 :: raw := row | throw (IO.userError "expected zero marker and string bytes")
  let input ← raw.mapM (wordInput 8)
  return String.intercalate ";" ((GoLean.Strings.fields input).map fun field =>
    "f" ++ String.intercalate "," (field.map toString))
