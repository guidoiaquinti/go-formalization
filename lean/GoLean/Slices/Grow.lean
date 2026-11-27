namespace GoLean.Slices

structure GrowView (α : Type) where
  cells : List α
  length : Nat
  isNil : Bool

structure GrowObservation (α : Type) where
  view : GrowView α
  allocated : Bool

/-- The cell list represents the entire old capacity and the requested
    extension. Any allocator-rounded extra capacity is outside the model.
    Go 1.27.0 src/slices/slices.go:417–430. -/
def sliceGrow (zero : α) (s : GrowView α) (n : Int) : Except String (GrowObservation α) :=
  if n < 0 then .error "cannot be negative"
  else
    let extra := n.toNat - (s.cells.length - s.length)
    if extra = 0 then .ok ⟨s, false⟩
    else .ok ⟨⟨s.cells ++ List.replicate extra zero, s.length, false⟩, true⟩

theorem sliceGrow_negative (zero : α) (s : GrowView α) (n : Int) (negative : n < 0) :
    sliceGrow zero s n = .error "cannot be negative" := by
  simp [sliceGrow, negative]

theorem sliceGrow_reuse (zero : α) (s : GrowView α) (n : Nat)
    (fits : s.length + n ≤ s.cells.length) :
    sliceGrow zero s n = .ok ⟨s, false⟩ := by
  have noExtra : n - (s.cells.length - s.length) = 0 := by omega
  simp [sliceGrow, noExtra]

theorem sliceGrow_allocate (zero : α) (s : GrowView α) (n : Nat)
    (valid : s.length ≤ s.cells.length) (needsSpace : s.cells.length < s.length + n) :
    sliceGrow zero s n =
      .ok ⟨⟨s.cells ++ List.replicate (s.length + n - s.cells.length) zero,
        s.length, false⟩, true⟩ := by
  have extra : n - (s.cells.length - s.length) = s.length + n - s.cells.length := by omega
  have positive : s.length + n - s.cells.length ≠ 0 := by omega
  simp [sliceGrow, extra, positive]

end GoLean.Slices
