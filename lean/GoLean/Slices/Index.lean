import Init.Data.List.Find

namespace GoLean.Slices

/-- Go 1.27.0 src/slices/slices.go:96-103 scans from the front. None
represents the result -1. The equality operation is the typed Go comparison. -/
def sliceIndex (equal : α → α → Bool) (xs : List α) (key : α) : Option Nat :=
  xs.findIdx? (equal key)

theorem sliceIndex_found (equal : α → α → Bool) (xs : List α) (key : α) (i : Nat) :
    sliceIndex equal xs key = some i ↔
      ∃ hi : i < xs.length, equal key xs[i] ∧
        ∀ j (hj : j < i), ¬equal key (xs[j]'(Nat.lt_trans hj hi)) :=
  List.findIdx?_eq_some_iff_getElem

theorem sliceIndex_missing (equal : α → α → Bool) (xs : List α) (key : α) :
    sliceIndex equal xs key = none ↔ ∀ x ∈ xs, equal key x = false :=
  List.findIdx?_eq_none_iff


end GoLean.Slices
