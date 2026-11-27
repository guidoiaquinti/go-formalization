import GoLean.Slices.Index

namespace GoLean.Slices

inductive StartIdxPanic where
  | emptyNeedle
  | needleNotFound

/-- Go 1.27.0 src/slices/slices.go:470–480. Pointer identities describe
    each haystack element and the first needle element, if present. -/
def sliceStartIdx (haystack : List Nat) (needle : Option Nat) : Except StartIdxPanic Nat :=
  match needle with
  | none => .error .emptyNeedle
  | some first =>
    match sliceIndex (· == ·) haystack first with
    | some i => .ok i
    | none => .error .needleNotFound

theorem sliceStartIdx_found (haystack : List Nat) (first i : Nat) :
    sliceStartIdx haystack (some first) = .ok i ↔
      ∃ hi : i < haystack.length, first = haystack[i] ∧
        ∀ j (hj : j < i), first ≠ haystack[j]'(Nat.lt_trans hj hi) := by
  have result :
      sliceStartIdx haystack (some first) = .ok i ↔ sliceIndex (· == ·) haystack first = some i := by
    unfold sliceStartIdx
    cases h : sliceIndex (· == ·) haystack first <;> simp [h]
  rw [result, sliceIndex_found]
  simp only [beq_iff_eq]

theorem sliceStartIdx_missing (haystack : List Nat) (first : Nat) :
    sliceStartIdx haystack (some first) = .error .needleNotFound ↔ first ∉ haystack := by
  have result :
      sliceStartIdx haystack (some first) = .error .needleNotFound ↔
        sliceIndex (· == ·) haystack first = none := by
    unfold sliceStartIdx
    cases h : sliceIndex (· == ·) haystack first <;> simp [h]
  rw [result, sliceIndex_missing]
  simp only [beq_eq_false_iff_ne]
  constructor
  · intro absent member
    exact absent first member rfl
  · intro absent x member eq
    exact absent (eq ▸ member)

theorem sliceStartIdx_empty (haystack : List Nat) :
    sliceStartIdx haystack none = .error .emptyNeedle := rfl

theorem sliceStartIdx_aligned (base length offset : Nat) (fits : offset < length) :
    sliceStartIdx ((List.range length).map (base + ·)) (some (base + offset)) = .ok offset := by
  apply (sliceStartIdx_found _ _ _).2
  refine ⟨by simpa using fits, by simp, ?_⟩
  intro j earlier
  simp only [List.getElem_map, List.getElem_range]
  omega

end GoLean.Slices
