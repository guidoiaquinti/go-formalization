import GoLean.Go.Types.VariableState
import GoDiff.Math.Bits.Tuple

open GoLean.Go.Types

private def render (encode : α → String) : Except ReadFault α → String
  | .error .nilReceiver => "panic"
  | .ok value => encode value

private def state (v : Variable) : String := s!"{v.kind.toNat}:{v.embedded}:{v.origin.isNone}"

def main : IO Unit := runTuple fun fields => do
  let mode :: raw :: replacement :: embedded :: pos :: pkg :: typ :: bytes := fields |
    throw (IO.userError "expected state and constructor arguments")
  let raw ← wordInput 8 raw
  let replacement ← wordInput 8 replacement
  let name ← bytes.mapM fun b => return (← wordInput 8 b).toNat
  let obj : VarObject := ⟨pos, if pkg == 0 then none else some 21, name,
    if typ == 0 then .nilInterface else if typ == 1 then .typedNil 1 else .value typ.toNat⟩
  let e := embedded == 1
  let v : Option Variable := if mode == 0 then none else some ⟨31, obj, raw, e,
    if mode == 2 then some 32 else none⟩
  return String.intercalate "|" [
    state (newVarRaw 41 raw obj), state (newVar 42 obj), state (newParam 43 obj), state (newField 44 obj e),
    render (toString ∘ BitVec.toNat) (varKind v), render toString (varAnonymous v),
    render toString (varEmbedded v), render toString (varIsField v), render toString (varOrigin v),
    render (fun updated => s!"{updated.kind.toNat}:{updated.embedded}:{updated.origin.getD updated.identity}")
      (varSetKind v replacement), varKindString raw]
