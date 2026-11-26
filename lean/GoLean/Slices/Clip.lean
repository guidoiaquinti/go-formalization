namespace GoLean.Slices

structure SliceView (α : Type) where
  elements : List α
  capacity : Nat
  isNil : Bool

/-- Go 1.27.0 src/slices/slices.go:432-437 applies a full slice expression
with start zero and both bounds equal to length. -/
def sliceClip (s : SliceView α) : SliceView α :=
  { s with capacity := s.elements.length }

theorem sliceClip_spec (s : SliceView α) :
    (sliceClip s).elements = s.elements ∧
    (sliceClip s).capacity = s.elements.length ∧
    (sliceClip s).isNil = s.isNil := ⟨rfl, rfl, rfl⟩


end GoLean.Slices
