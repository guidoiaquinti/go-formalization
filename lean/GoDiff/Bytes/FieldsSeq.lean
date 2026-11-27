import GoLean.Bytes.FieldsSeq
import GoDiff.Math.Bits.Tuple

private def describe (fields : List (List Nat × Nat)) : String :=
  String.intercalate ";" (fields.map fun field =>
    toString field.2 ++ ":" ++ String.intercalate "," (field.1.map toString))

def main : IO Unit := runTuple fun row => do
  let budget :: input := row | throw (IO.userError "expected a callback budget and string bytes")
  let bytes ← input.mapM (wordInput 8)
  return describe (GoLean.Bytes.fieldsSeq bytes (some budget.toNat)) ++ "|" ++
    describe (GoLean.Bytes.fieldsSeq bytes none)
