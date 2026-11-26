namespace GoLean.Slices

/-- Go 1.27.0 src/slices/slices.go:14-30 checks length before comparing elements. -/
def sliceEqual (equal : α → α → Bool) (xs ys : List α) : Bool :=
  (xs.length == ys.length) && xs.isEqv ys equal

theorem sliceEqual_spec [BEq α] [LawfulBEq α] (xs ys : List α) :
    sliceEqual (· == ·) xs ys = true ↔ xs = ys := by
  simp [sliceEqual, List.isEqv_eq]
  exact congrArg List.length


end GoLean.Slices
