import GoLean.Slices.Index

namespace GoLean.Slices

/-- Go 1.27.0 src/slices/slices.go:116-119. A successful index is nonnegative. -/
def sliceContains (equal : α → α → Bool) (xs : List α) (key : α) : Bool :=
  (sliceIndex equal xs key).isSome

theorem sliceContains_spec (equal : α → α → Bool) (xs : List α) (key : α) :
    sliceContains equal xs key = true ↔ ∃ x ∈ xs, equal key x = true := by
  simp [sliceContains, sliceIndex, List.findIdx?_isSome, List.any_eq_true]


end GoLean.Slices
