import GoLean.Strings.FieldsSeq
import GoDiff.Math.Bits.Tuple

private def describe (fields : List (List Nat)) : String :=
  String.intercalate ";" (fields.map fun field =>
    String.intercalate "," (field.map toString))

def main : IO Unit := runTuple fun row => do
  let budget :: input := row | throw (IO.userError "expected a callback budget and string bytes")
  let bytes ← input.mapM (wordInput 8)
  return describe (GoLean.Strings.fieldsSeq bytes (some budget.toNat)) ++ "|" ++
    describe (GoLean.Strings.fieldsSeq bytes none)
