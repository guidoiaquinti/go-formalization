import GoLean.Go.Types.VariableState

namespace GoLean.Go.Types

structure TypedObject where
  identity : Nat
  obj : VarObject

structure ConstantObject where
  base : TypedObject
  value : TypeRef

structure PackageNameObject where
  base : TypedObject
  imported : Option Nat

structure LabelObject where
  base : TypedObject
  used : Bool

structure FunctionObject where
  base : TypedObject
  typeIsSignature : Bool
  signatureScope : Option Nat
  origin : Option Nat

def newConst (id : Nat) (obj : VarObject) (value : TypeRef) : ConstantObject := ⟨⟨id, obj⟩, value⟩

def newPkgName (id : Nat) (obj : VarObject) (imported : Option Nat) : PackageNameObject :=
  ⟨⟨id, {obj with typ := .value 9}⟩, imported⟩

def newTypeName (id : Nat) (obj : VarObject) : TypedObject := ⟨id, obj⟩

def newLabel (id : Nat) (obj : VarObject) : LabelObject :=
  ⟨⟨id, {obj with typ := .value 9}⟩, false⟩

def newFunc (id : Nat) (obj : VarObject) (signature scope : Option Nat) : FunctionObject :=
  ⟨⟨id, {obj with typ := (signature.map TypeRef.value).getD .nilInterface}⟩, signature.isSome, scope, none⟩

private def readObj (get : α → β) (obj : Option α) : Except ReadFault β :=
  match obj with
  | none => .error .nilReceiver
  | some obj => .ok (get obj)

def constVal := readObj ConstantObject.value

def pkgNameImported := readObj PackageNameObject.imported

def funcPkg := readObj (fun f : FunctionObject => f.base.obj.pkg)

def funcOrigin := readObj (fun f : FunctionObject => f.origin.getD f.base.identity)

inductive SignatureObservation where
  | fresh
  | reference (id : Nat)
  | typedNil
  | panic

def funcSignature (obj : Option FunctionObject) : SignatureObservation :=
  match obj with
  | none => .panic
  | some obj =>
    match obj.base.obj.typ with
    | .nilInterface => .fresh
    | .typedNil _ => if obj.typeIsSignature then .typedNil else .panic
    | .value id => if obj.typeIsSignature then .reference id else .panic

def funcScope (obj : Option FunctionObject) : Except ReadFault (Option Nat) :=
  match obj with
  | some obj =>
    match obj.base.obj.typ with
    | .value _ => if obj.typeIsSignature then .ok obj.signatureScope else .error .nilReceiver
    | _ => .error .nilReceiver
  | none => .error .nilReceiver

theorem newConst_stores (id : Nat) (obj : VarObject) (val : TypeRef) :
    (newConst id obj val).value = val := rfl
theorem constVal_reads (obj : ConstantObject) : constVal (some obj) = .ok obj.value := rfl
theorem newPkgName_stores (id : Nat) (obj : VarObject) (p : Option Nat) :
    (newPkgName id obj p).imported = p := rfl
theorem pkgNameImported_reads (obj : PackageNameObject) :
    pkgNameImported (some obj) = .ok obj.imported := rfl
theorem newTypeName_stores (id : Nat) (obj : VarObject) : (newTypeName id obj).obj = obj := rfl
theorem newLabel_unused (id : Nat) (obj : VarObject) : (newLabel id obj).used = false := rfl
theorem newFunc_nilSignature (id : Nat) (obj : VarObject) (scope : Option Nat) :
    (newFunc id obj none scope).base.obj.typ = .nilInterface := rfl
theorem funcPkg_reads (obj : FunctionObject) : funcPkg (some obj) = .ok obj.base.obj.pkg := rfl
theorem funcOrigin_reads (obj : FunctionObject) :
    funcOrigin (some obj) = .ok (obj.origin.getD obj.base.identity) := rfl
theorem funcSignature_fresh (id : Nat) (obj : VarObject) (scope : Option Nat) :
    funcSignature (some (newFunc id obj none scope)) = .fresh := rfl
theorem funcScope_live (id : Nat) (obj : VarObject) (scope : Option Nat) :
    funcScope (some (newFunc id obj (some 10) scope)) = .ok scope := rfl

end GoLean.Go.Types
