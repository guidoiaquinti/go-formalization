namespace GoLean.Slices

/-- Go 1.27.0 src/slices/slices.go:224–235. The optional result represents
a bounds panic, while the second field records the original backing region. -/
def sliceDelete (zero : α) (xs : List α) (i j : Int) : Option (List α) × List α :=
  if 0 ≤ i ∧ i ≤ j ∧ j ≤ (xs.length : Int) then
    if i = j then (some xs, xs)
    else
      let kept := xs.take i.toNat ++ xs.drop j.toNat
      (some kept, kept ++ List.replicate (j - i).toNat zero)
  else (none, xs)

theorem sliceDelete_valid (zero : α) (xs : List α) (i j : Int)
    (bounds : 0 ≤ i ∧ i ≤ j ∧ j ≤ (xs.length : Int)) :
    let kept := xs.take i.toNat ++ xs.drop j.toNat
    sliceDelete zero xs i j = (some kept, kept ++ List.replicate (j - i).toNat zero) ∧
      kept.length + (j - i).toNat = xs.length := by
  have iBound : i.toNat ≤ xs.length := by omega
  have jBound : j.toNat ≤ xs.length := by omega
  constructor
  · by_cases empty : i = j
    · subst j
      simp [sliceDelete, bounds, List.take_append_drop]
    · simp [sliceDelete, bounds, empty]
  · simp only [List.length_append, List.length_take, List.length_drop, Nat.min_eq_left iBound]
    omega

theorem sliceDelete_invalid (zero : α) (xs : List α) (i j : Int)
    (bounds : ¬(0 ≤ i ∧ i ≤ j ∧ j ≤ (xs.length : Int))) :
    sliceDelete zero xs i j = (none, xs) := by
  simp [sliceDelete, bounds]

end GoLean.Slices
