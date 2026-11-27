import GoLean.Go.Types.ObjectConstructors
import GoDiff.Math.Bits.Tuple

open GoLean.Go.Types

private def out (f : α → String) : Except ReadFault α → String
  | .error _ => "panic"
  | .ok a => f a

private def opt (n : Option Nat) := (n.map toString).getD "nil"
private def typ (t : TypeRef) : String :=
  match t with
  | .nilInterface => "nil"
  | .typedNil n => s!"typedNil{n}"
  | .value n => toString n

def main : IO Unit := runTuple fun row => do
  let mode :: pos :: pkg :: kind :: value :: imp :: sig :: sc :: raw := row |
    throw (IO.userError "expected receiver and object constructor fields")
  let pkgID := if pkg == 0 then none else some 21
  let imported := if imp == 0 then none else some 22
  let scope := if sc == 0 then none else some 50
  let name := raw.map Int.toNat
  let t := if kind == 0 then TypeRef.nilInterface else if kind == 1 then .typedNil 1 else .value 11
  let v := if value == 0 then TypeRef.nilInterface else if value == 1 then .typedNil 2 else .value 12
  let o : VarObject := ⟨pos, pkgID, name, t⟩
  let c := newConst 31 o v
  let p := newPkgName 32 o imported
  let n := newTypeName 33 o
  let l := newLabel 34 o
  let f := newFunc 35 o (if sig == 0 then none else some 10) scope
  let syntheticType := if sig == 2 then TypeRef.typedNil 10
    else if sig == 3 then .value 11 else f.base.obj.typ
  let synthetic := { f with
    base := { f.base with obj := { f.base.obj with typ := syntheticType } },
    typeIsSignature := sig == 1 || sig == 2,
    origin := if mode == 2 then some 36 else none }
  let recv := if mode == 0 then none else some synthetic
  let signature := match funcSignature recv with
    | .fresh => "fresh"
    | .reference id => toString id
    | .typedNil => "typedNil10"
    | .panic => "panic"
  return String.intercalate "|" [
    typ c.base.obj.typ ++ ":" ++ typ c.value, typ p.base.obj.typ ++ ":" ++ opt p.imported,
    typ n.obj.typ, typ l.base.obj.typ ++ ":" ++ toString l.used, typ f.base.obj.typ,
    out typ (constVal (if mode == 0 then none else some c)),
    out opt (pkgNameImported (if mode == 0 then none else some p)),
    signature, out opt (funcScope recv), out toString (funcOrigin recv), out opt (funcPkg recv)]
