import GoLean.Go.Types.ArrayChan

namespace GoLean.Go.Types

structure BasicType where
  identity : Nat
  kind : Int
  info : Int
  name : List Nat

private def basicRead (field : BasicType → α) (receiver : Option BasicType) : Except ReadFault α :=
  match receiver with
  | none => .error .nilReceiver
  | some b => .ok (field b)

/-- Pinned Go 1.27.0 src/go/types/basic.go:76–85. -/
def basicKind := basicRead BasicType.kind
def basicInfo := basicRead BasicType.info
def basicName := basicRead BasicType.name

def basicUnderlying : Option BasicType → TypeRef
  | none => .typedNil 6
  | some b => .value b.identity

private def asciiBytes (s : String) : List Nat := s.toList.map Char.toNat

/-- The Basic arm of typestring.go uses the initialized unsafe scope.
    Its sole exported TypeName is Pointer. -/
def basicString (receiver : Option BasicType) : Except ReadFault (List Nat) :=
  (basicName receiver).map fun name =>
    if name == asciiBytes "Pointer" then asciiBytes "unsafe.Pointer" else name

theorem basicKind_reads (b : BasicType) :
    basicKind (some b) = .ok b.kind ∧ basicKind none = .error .nilReceiver :=
  ⟨rfl, rfl⟩

theorem basicInfo_reads (b : BasicType) :
    basicInfo (some b) = .ok b.info ∧ basicInfo none = .error .nilReceiver :=
  ⟨rfl, rfl⟩

theorem basicName_reads (b : BasicType) :
    basicName (some b) = .ok b.name ∧ basicName none = .error .nilReceiver :=
  ⟨rfl, rfl⟩

theorem basicUnderlying_identity (b : BasicType) :
    basicUnderlying (some b) = .value b.identity ∧ basicUnderlying none = .typedNil 6 :=
  ⟨rfl, rfl⟩

theorem basicString_pointer (identity : Nat) (kind info : Int) :
    basicString (some ⟨identity, kind, info, asciiBytes "Pointer"⟩) =
      .ok (asciiBytes "unsafe.Pointer") := rfl

theorem basicString_plain (b : BasicType) (plain : b.name ≠ asciiBytes "Pointer") :
    basicString (some b) = .ok b.name := by
  simp [basicString, basicName, basicRead, Except.map, plain]

end GoLean.Go.Types
