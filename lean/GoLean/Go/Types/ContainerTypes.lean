import GoLean.Go.Types.ArrayChan

namespace GoLean.Go.Types

structure SliceType where
  identity : Nat
  elem : TypeRef

structure PointerType where
  identity : Nat
  base : TypeRef

structure MapType where
  identity : Nat
  key : TypeRef
  elem : TypeRef

private def readField (access : α → TypeRef) (receiver : Option α) : Except ReadFault TypeRef :=
  match receiver with
  | none => .error .nilReceiver
  | some value => .ok (access value)

private def asType (identity : α → Nat) (tag : Nat) (receiver : Option α) : TypeRef :=
  match receiver with
  | none => .typedNil tag
  | some value => .value (identity value)

/-- Go 1.27.0 src/go/types/slice.go:16–21. Identities stand for pointer equality. -/
def newSlice (id : Nat) (element : TypeRef) : SliceType := ⟨id, element⟩
def sliceElem := readField SliceType.elem
def sliceUnderlying := asType SliceType.identity 3

/-- Go 1.27.0 src/go/types/pointer.go:16–21. -/
def newPointer (id : Nat) (base : TypeRef) : PointerType := ⟨id, base⟩
def pointerElem := readField PointerType.base
def pointerUnderlying := asType PointerType.identity 4

/-- Go 1.27.0 src/go/types/map.go:16–26. Construction stores both interfaces. -/
def newMap (id : Nat) (key element : TypeRef) : MapType := ⟨id, key, element⟩
def mapKey := readField MapType.key
def mapElem := readField MapType.elem
def mapUnderlying := asType MapType.identity 5

theorem newSlice_fields (id : Nat) (elem : TypeRef) :
    (newSlice id elem).identity = id ∧
      (newSlice id elem).elem = elem := ⟨rfl, rfl⟩

theorem sliceElem_reads (c : SliceType) :
    sliceElem (some c) = .ok c.elem ∧
      sliceElem none = .error .nilReceiver := ⟨rfl, rfl⟩

theorem sliceUnderlying_self (c : SliceType) :
    sliceUnderlying (some c) = .value c.identity ∧
      sliceUnderlying none = .typedNil 3 := ⟨rfl, rfl⟩

theorem newPointer_fields (id : Nat) (base : TypeRef) :
    (newPointer id base).identity = id ∧
      (newPointer id base).base = base := ⟨rfl, rfl⟩

theorem pointerElem_reads (c : PointerType) :
    pointerElem (some c) = .ok c.base ∧
      pointerElem none = .error .nilReceiver := ⟨rfl, rfl⟩

theorem pointerUnderlying_self (c : PointerType) :
    pointerUnderlying (some c) = .value c.identity ∧
      pointerUnderlying none = .typedNil 4 := ⟨rfl, rfl⟩

theorem newMap_fields (id : Nat) (key elem : TypeRef) :
    (newMap id key elem).identity = id ∧
      (newMap id key elem).key = key ∧
      (newMap id key elem).elem = elem := ⟨rfl, rfl, rfl⟩

theorem mapKey_reads (c : MapType) :
    mapKey (some c) = .ok c.key ∧
      mapKey none = .error .nilReceiver := ⟨rfl, rfl⟩

theorem mapElem_reads (c : MapType) :
    mapElem (some c) = .ok c.elem ∧
      mapElem none = .error .nilReceiver := ⟨rfl, rfl⟩

theorem mapUnderlying_self (c : MapType) :
    mapUnderlying (some c) = .value c.identity ∧
      mapUnderlying none = .typedNil 5 := ⟨rfl, rfl⟩

end GoLean.Go.Types
