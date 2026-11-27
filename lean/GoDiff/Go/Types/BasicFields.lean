import GoLean.Go.Types.BasicFields
import GoDiff.Math.Bits.Tuple

open GoLean.Go.Types

private def observe (render : α → String) : Except ReadFault α → String
  | .error .nilReceiver => "panic"
  | .ok value => render value

private def bytesWire (bytes : List Nat) : String :=
  String.intercalate "," (bytes.map toString)

def main : IO Unit := runTuple fun fields => do
  let presence :: kind :: info :: bytes := fields | throw (IO.userError "expected Basic fields")
  let name ← bytes.mapM fun b => return (← wordInput 8 b).toNat
  let value : Option BasicType := if presence == 0 then none else some ⟨31, kind, info, name⟩
  let identity := match basicUnderlying value with
    | .typedNil 6 => "typedNil"
    | .value 31 => "same"
    | _ => "wrong"
  return String.intercalate "|" [
    observe toString (basicKind value), observe toString (basicInfo value),
    observe bytesWire (basicName value), identity, observe bytesWire (basicString value)]
