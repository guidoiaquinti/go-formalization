import GoLean.Sort.IsSorted

namespace GoLean.Slices

/-- Go 1.27.0 src/slices/sort.go:51-60. The callback receives the
right element first. Only a negative result is an inversion. -/
def isSortedFunc (xs : List α) (fallback : α) (cmp : α → α → Except Unit Int) :
    Except Unit Bool × List Nat :=
  GoLean.Sort.isSorted xs.length fun i j =>
    (cmp (xs[i]?.getD fallback) (xs[j]?.getD fallback)).map fun order => decide (order < 0)

theorem isSortedFunc_spec (xs : List α) (fallback : α) (cmp : α → α → Except Unit Int) :
    (isSortedFunc xs fallback cmp).1 = .ok true ↔
      ∀ i < xs.length - 1, ∃ order,
        cmp (xs[i + 1]?.getD fallback) (xs[i]?.getD fallback) = .ok order ∧ 0 ≤ order := by
  rw [isSortedFunc, GoLean.Sort.isSorted_spec]
  apply forall_congr'
  intro i
  apply forall_congr'
  intro _
  cases h : cmp (xs[i + 1]?.getD fallback) (xs[i]?.getD fallback) <;>
    simp [Except.map, Int.not_lt]

end GoLean.Slices
