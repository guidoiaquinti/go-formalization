import GoLean.Go.Types.ArrayChan

namespace GoLean.Go.Types

structure VarObject where
  position : Int
  pkg : Option Nat
  name : List Nat
  typ : TypeRef

structure Variable where
  identity : Nat
  obj : VarObject
  kind : BitVec 8
  embedded : Bool
  origin : Option Nat

/-- Go 1.27.0 src/go/types/object.go:333–393. -/
def newVarRaw (id : Nat) (kind : BitVec 8) (obj : VarObject) : Variable :=
  ⟨id, obj, kind, false, none⟩

def newVar (id : Nat) (obj : VarObject) := newVarRaw id 1 obj
def newParam (id : Nat) (obj : VarObject) := newVarRaw id 4 obj
def newField (id : Nat) (obj : VarObject) (embedded : Bool) : Variable :=
  { newVarRaw id 6 obj with embedded }

private def readVar (field : Variable → α) (receiver : Option Variable) : Except ReadFault α :=
  match receiver with
  | none => .error .nilReceiver
  | some v => .ok (field v)

def varKind := readVar Variable.kind
def varEmbedded := readVar Variable.embedded
def varAnonymous := readVar Variable.embedded
def varIsField := readVar (fun v => v.kind == 6)
def varOrigin := readVar (fun v => v.origin.getD v.identity)

def varSetKind (receiver : Option Variable) (kind : BitVec 8) : Except ReadFault Variable :=
  readVar (fun v => { v with kind }) receiver

private def kindNames : Array String := #[
  "VarKind(0)", "PackageVar", "LocalVar", "RecvVar", "ParamVar", "ResultVar", "FieldVar"]

def varKindString (kind : BitVec 8) : String :=
  kindNames[kind.toNat]?.getD s!"VarKind({kind.toNat})"

theorem newVarRaw_initial (id : Nat) (kind : BitVec 8) (obj : VarObject) :
    newVarRaw id kind obj = ⟨id, obj, kind, false, none⟩ := rfl

theorem newVar_initial (id : Nat) (obj : VarObject) :
    newVar id obj = ⟨id, obj, 1, false, none⟩ := rfl

theorem newParam_initial (id : Nat) (obj : VarObject) :
    newParam id obj = ⟨id, obj, 4, false, none⟩ := rfl

theorem newField_initial (id : Nat) (obj : VarObject) (e : Bool) :
    newField id obj e = ⟨id, obj, 6, e, none⟩ := rfl

theorem varSetKind_changes_only_kind (v : Variable) (kind : BitVec 8) :
    varSetKind (some v) kind = .ok { v with kind } ∧
      varSetKind none kind = .error .nilReceiver := ⟨rfl, rfl⟩

theorem varKindString_boundaries :
    varKindString 0 = "VarKind(0)" ∧ varKindString 6 = "FieldVar" ∧
      varKindString 255 = "VarKind(255)" := ⟨rfl, rfl, rfl⟩

theorem varKind_reads (v : Variable) :
    varKind (some v) = .ok v.kind ∧ varKind none = .error .nilReceiver :=
  ⟨rfl, rfl⟩

theorem varAnonymous_reads (v : Variable) :
    varAnonymous (some v) = .ok v.embedded ∧ varAnonymous none = .error .nilReceiver :=
  ⟨rfl, rfl⟩

theorem varEmbedded_reads (v : Variable) :
    varEmbedded (some v) = .ok v.embedded ∧ varEmbedded none = .error .nilReceiver :=
  ⟨rfl, rfl⟩

theorem varIsField_reads (v : Variable) :
    varIsField (some v) = .ok (v.kind == 6) ∧ varIsField none = .error .nilReceiver :=
  ⟨rfl, rfl⟩

theorem varOrigin_reads (v : Variable) :
    varOrigin (some v) = .ok (v.origin.getD v.identity) ∧ varOrigin none = .error .nilReceiver :=
  ⟨rfl, rfl⟩

end GoLean.Go.Types
