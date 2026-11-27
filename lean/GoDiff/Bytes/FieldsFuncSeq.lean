import GoLean.Bytes.FieldsFuncSeq
import GoDiff.Math.Bits.Tuple

open GoLean.Bytes

private def descriptor (kind : Nat) : Nat → Bool :=
  match kind with
  | 0 => fun c => c == 32 || c == 0xa0 || c == 0x3000
  | 1 => fun _ => false
  | 2 => fun _ => true
  | 3 => fun c => c == 0xfffd
  | _ => fun c => c % 5 == 2

private def describe (events : List FieldEvent) : String :=
  String.intercalate ";" (events.map fun event =>
    match event with
    | .predicate r => "p" ++ toString r
    | .yield field capacity => "y" ++ toString capacity ++ ":" ++ String.intercalate "," (field.map toString))

def main : IO Unit := runTuple fun row => do
  let kind :: budget :: input := row | throw (IO.userError "expected predicate, budget and bytes")
  let bytes ← input.mapM (wordInput 8)
  let predicate := descriptor kind.toNat
  return describe (fieldsFuncSeq bytes predicate (some budget.toNat)) ++ "|" ++
    describe (fieldsFuncSeq bytes predicate none)
