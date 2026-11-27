import GoLean.Compress.Flate.ErrorText
import GoDiff.Math.Bits.Tuple

open GoLean.Compress.Flate

private def wire (bytes : ErrorBytes) :=
  String.intercalate "," (bytes.map (toString ∘ BitVec.toNat))

private def result : ErrorResult → String
  | .nilAccess => "nil:0"
  | .panicked count => s!"panic:{count}"
  | .returned count bytes => s!"ok:{count}:{wire bytes}"

def main : IO Unit := runTuple fun row => do
  let present :: offset :: mode :: raw := row | throw (IO.userError "expected error fixture")
  let word := BitVec.ofInt 64 offset
  let bytes ← raw.mapM (wordInput 8)
  let callback := if mode == 0 then .absent else if mode == 2 then .panics else ErrorCall.returns bytes
  let offset := if present == 0 then none else some word
  return String.intercalate "|" [wire (corruptInputError word), wire (internalError bytes),
    result (readError offset callback), result (writeError offset callback)]
