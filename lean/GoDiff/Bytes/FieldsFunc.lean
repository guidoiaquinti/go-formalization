import GoLean.Bytes.FieldsFunc
import GoDiff.Math.Bits.Tuple

private def descriptor (kind : Nat) : Nat → Bool :=
  match kind with
  | 0 => fun c => c == 32 || c == 0xa0 || c == 0x3000
  | 1 => fun _ => false
  | 2 => fun _ => true
  | 3 => fun c => c == 0xfffd
  | _ => fun c => c % 5 == 2

def main : IO Unit := runTuple fun row => do
  let kind :: raw := row | throw (IO.userError "expected predicate and bytes")
  let input ← raw.mapM (wordInput 8)
  let (calls, fields) := GoLean.Bytes.fieldsFunc input (descriptor kind.toNat)
  let encoded := fields.map fun (field, cap) => "f" ++ toString cap ++ ":" ++ String.intercalate "," (field.map toString)
  return String.intercalate "," (calls.map toString) ++ "|" ++ String.intercalate ";" encoded
