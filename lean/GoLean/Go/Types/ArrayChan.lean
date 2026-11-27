import Init.Data.BitVec.Basic

namespace GoLean.Go.Types

/-- Opaque interface observations preserve nilness, dynamic typed nil, and identity. -/
inductive TypeRef where
  | nilInterface
  | typedNil (tag : Nat)
  | value (identity : Nat)

inductive ReadFault where
  | nilReceiver

structure ArrayType where
  identity : Nat
  elem : TypeRef
  len : BitVec 64

structure ChanType (width : Nat) where
  identity : Nat
  dir : BitVec width
  elem : TypeRef

/-- Go 1.27.0 src/go/types/array.go:18–27. The allocation identity is supplied. -/
def newArray (identity : Nat) (elem : TypeRef) (len : BitVec 64) : ArrayType :=
  ⟨identity, elem, len⟩

def arrayLen (receiver : Option ArrayType) : Except ReadFault (BitVec 64) :=
  match receiver with
  | none => .error .nilReceiver
  | some a => .ok a.len

def arrayElem (receiver : Option ArrayType) : Except ReadFault TypeRef :=
  match receiver with
  | none => .error .nilReceiver
  | some a => .ok a.elem

def arrayUnderlying (receiver : Option ArrayType) : TypeRef :=
  match receiver with
  | none => .typedNil 1
  | some a => .value a.identity

/-- Go 1.27.0 src/go/types/chan.go:27–37. Direction retains the native int word. -/
def newChan (identity : Nat) (dir : BitVec width) (elem : TypeRef) : ChanType width :=
  ⟨identity, dir, elem⟩

def chanDir (receiver : Option (ChanType width)) : Except ReadFault (BitVec width) :=
  match receiver with
  | none => .error .nilReceiver
  | some c => .ok c.dir

def chanElem (receiver : Option (ChanType width)) : Except ReadFault TypeRef :=
  match receiver with
  | none => .error .nilReceiver
  | some c => .ok c.elem

def chanUnderlying (receiver : Option (ChanType width)) : TypeRef :=
  match receiver with
  | none => .typedNil 2
  | some c => .value c.identity

theorem newArray_fields (i : Nat) (e : TypeRef) (n : BitVec 64) :
    (newArray i e n).identity = i ∧ (newArray i e n).elem = e ∧
      (newArray i e n).len = n := ⟨rfl, rfl, rfl⟩

theorem arrayLen_reads (a : ArrayType) :
    arrayLen (some a) = .ok a.len ∧ arrayLen none = .error .nilReceiver := ⟨rfl, rfl⟩

theorem arrayElem_reads (a : ArrayType) :
    arrayElem (some a) = .ok a.elem ∧ arrayElem none = .error .nilReceiver := ⟨rfl, rfl⟩

theorem arrayUnderlying_self (a : ArrayType) :
    arrayUnderlying (some a) = .value a.identity ∧ arrayUnderlying none = .typedNil 1 :=
  ⟨rfl, rfl⟩

theorem newChan_fields (i : Nat) (d : BitVec width) (e : TypeRef) :
    (newChan i d e).identity = i ∧ (newChan i d e).dir = d ∧
      (newChan i d e).elem = e := ⟨rfl, rfl, rfl⟩

theorem chanDir_reads (c : ChanType width) :
    chanDir (some c) = .ok c.dir ∧ chanDir (width := width) none = .error .nilReceiver :=
  ⟨rfl, rfl⟩

theorem chanElem_reads (c : ChanType width) :
    chanElem (some c) = .ok c.elem ∧ chanElem (width := width) none = .error .nilReceiver :=
  ⟨rfl, rfl⟩

theorem chanUnderlying_self (c : ChanType width) :
    chanUnderlying (some c) = .value c.identity ∧ chanUnderlying (width := width) none = .typedNil 2 :=
  ⟨rfl, rfl⟩

end GoLean.Go.Types
