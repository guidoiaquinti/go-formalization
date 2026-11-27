import GoLean.Bytes.Fields
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let 0 :: raw := row | throw (IO.userError "expected zero marker and slice bytes")
  let input ← raw.mapM (wordInput 8)
  match GoLean.Bytes.fields input with
  | none => return "nil"
  | some fields =>
    return "slice" ++ String.intercalate ";" (fields.map fun (field, capacity) =>
      "f" ++ toString capacity ++ ":" ++ String.intercalate "," (field.map toString))
