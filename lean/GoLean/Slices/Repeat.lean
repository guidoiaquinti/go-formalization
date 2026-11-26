namespace GoLean.Slices

/-- Go 1.27.0 src/slices/slices.go:509-531. The exact product represents
the two-word multiplication and signed length check. The successful output
represents the initialized slice contents after the doubling copies. -/
def sliceRepeat (maxLength : Nat) (xs : List α) (count : Int) : Except String (List α) :=
  if count < 0 then .error "cannot be negative"
  else if xs.length * count.toNat > maxLength then
    .error "the result of (len(x) * count) overflows"
  else if xs.isEmpty then .ok []
  else .ok ((List.replicate count.toNat xs).flatten)

theorem sliceRepeat_spec (maxLength : Nat) (xs : List α) (count : Int)
    (nonnegative : 0 ≤ count) (fits : xs.length * count.toNat ≤ maxLength) :
    sliceRepeat maxLength xs count = .ok ((List.replicate count.toNat xs).flatten) ∧
      (List.replicate count.toNat xs).flatten.length = xs.length * count.toNat := by
  constructor
  · simp only [sliceRepeat, Int.not_lt.mpr nonnegative, Nat.not_lt.mpr fits, ite_false]
    cases xs <;> simp
  · simp [Nat.mul_comm]

theorem sliceRepeat_error (maxLength : Nat) (xs : List α) (count : Int) :
    (sliceRepeat maxLength xs count = .error "cannot be negative" ↔ count < 0) ∧
    (sliceRepeat maxLength xs count = .error "the result of (len(x) * count) overflows" ↔
      0 ≤ count ∧ maxLength < xs.length * count.toNat) := by
  by_cases h : count < 0
  · simp [sliceRepeat, h, show ¬0 ≤ count by omega]
  · have nonnegative : 0 ≤ count := by omega
    by_cases overflow : maxLength < xs.length * count.toNat
    · simp [sliceRepeat, h, overflow, nonnegative]
    · simp only [sliceRepeat, h, overflow, ite_false, and_false]
      cases xs <;> simp

end GoLean.Slices
