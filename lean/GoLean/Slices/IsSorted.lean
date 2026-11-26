import GoLean.Sort.IsSorted

namespace GoLean.Slices

/-- Go 1.27.0 src/slices/sort.go:41-49 scans for an adjacent inversion
using the supplied specialization of cmp.Less. -/
def isSorted (less : α → α → Bool) (fallback : α) (xs : List α) : Bool :=
  GoLean.Sort.isSortedList less fallback xs

theorem isSorted_spec (less : α → α → Bool) (fallback : α) (xs : List α) :
    isSorted less fallback xs = true ↔
      ∀ i < xs.length - 1, less (xs[i + 1]?.getD fallback) (xs[i]?.getD fallback) = false :=
  GoLean.Sort.isSortedList_spec less fallback xs


end GoLean.Slices
