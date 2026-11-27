namespace GoLean.Slices

/-- Go 1.27.0 src/slices/slices.go:455–469. Addresses are mathematical
    byte addresses of contiguous storage. The model assumes the additions
    for the last element and its last byte are representable as uintptr. -/
def sliceOverlaps (aStart aLength bStart bLength size : Nat) : Bool :=
  if aLength = 0 ∨ bLength = 0 ∨ size = 0 then false
  else decide
    (aStart ≤ bStart + (bLength - 1) * size + (size - 1) ∧
     bStart ≤ aStart + (aLength - 1) * size + (size - 1))

private theorem endAddress (n size : Nat) (hn : n ≠ 0) (hs : size ≠ 0) :
    (n - 1) * size + (size - 1) + 1 = n * size := by
  rw [Nat.add_assoc, Nat.sub_add_cancel (show 1 ≤ size by omega),
    ← Nat.succ_mul, show (n - 1).succ = n by omega]

theorem sliceOverlaps_spec (aStart aLength bStart bLength size : Nat) :
    sliceOverlaps aStart aLength bStart bLength size = true ↔
      ∃ address : Nat, aStart ≤ address ∧ address < aStart + aLength * size ∧
        bStart ≤ address ∧ address < bStart + bLength * size := by
  by_cases ha : aLength = 0
  · simp [sliceOverlaps, ha]
    intro address lower upper
    omega
  by_cases hb : bLength = 0
  · simp [sliceOverlaps, hb]
  by_cases hs : size = 0
  · simp [sliceOverlaps, hs]
  simp [sliceOverlaps, ha, hb, hs]
  have endA := endAddress aLength size ha hs
  have endB := endAddress bLength size hb hs
  constructor
  · intro overlap
    have bounds := overlap.1
    have other := overlap.2
    have nonemptyA : 0 < aLength * size := Nat.mul_pos (by omega) (by omega)
    have nonemptyB : 0 < bLength * size := Nat.mul_pos (by omega) (by omega)
    refine ⟨max aStart bStart, ?_⟩
    omega
  · rintro ⟨address, aLower, aUpper, bLower, bUpper⟩
    exact ⟨by omega, by omega⟩

end GoLean.Slices
