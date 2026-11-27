import GoLean.Encoding.Json.CompatErrors
import GoDiff.Math.Bits.Tuple

open GoLean.Encoding.Json

private def bytes (row : List Int) : IO (ErrorBytes × List Int) := do
  let count :: rest := row | throw (IO.userError "expected byte count")
  return ((rest.take count.toNat).map Int.toNat, rest.drop count.toNat)

private def wire (v : ErrorBytes) := String.intercalate "," (v.map toString)

private def outcome (render : α → String) : ErrorOutcome α → String
  | .ok a => render a
  | .nilPanic => "nilpanic"
  | .nestedPanic => "nestedpanic"

def main : IO Unit := runTuple fun fields => do
  let present :: kind :: tag :: rest := fields | throw (IO.userError "expected receiver metadata")
  let (typeText, rest) ← bytes rest
  let (nestedText, rest) ← bytes rest
  let (value, rest) ← bytes rest
  let (structName, rest) ← bytes rest
  let (field, rest) ← bytes rest
  let (source, _) ← bytes rest
  let typ := if kind == 0 then none else some ⟨kind == 2, typeText⟩
  let nested := if tag == 0 then .absent else if tag == 2 then .panic 2 else .text 1 nestedText
  let error : Option CompatibilityError :=
    if present == 0 then none else some ⟨typ, nested, value, structName, field, source⟩
  return String.intercalate "|" [
    outcome wire (syntaxError error), outcome wire (unsupportedValueError error),
    outcome wire (unsupportedTypeError error), outcome wire (invalidUnmarshalError error),
    outcome wire (marshalerError error), outcome toString (marshalerErrorUnwrap error),
    outcome wire (unmarshalTypeError error), outcome toString (unmarshalTypeErrorUnwrap error)]
