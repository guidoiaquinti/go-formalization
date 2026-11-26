import GoLean.Sort.SearchInterval

namespace GoLean.Slices

/-- Go 1.27.0 src/slices/sort.go:126-143. The comparison parameters
represent cmp.Less and the final equality check for the element type. -/
def binarySearch (less equal : α → α → Bool) (fallback : α)
    (xs : List α) (target : α) : Nat × Bool :=
  let i := GoLean.Sort.searchInterval (fun j => !less (xs[j]?.getD fallback) target) 0 xs.length
  (i, decide (i < xs.length) && equal (xs[i]?.getD fallback) target)

theorem binarySearch_spec (less equal : α → α → Bool) (fallback : α)
    (xs : List α) (target : α) (cut : Nat) (bounded : cut ≤ xs.length)
    (partition : ∀ i < xs.length, less (xs[i]?.getD fallback) target = false ↔ cut ≤ i) :
    binarySearch less equal fallback xs target =
      (cut, decide (cut < xs.length) && equal (xs[cut]?.getD fallback) target) := by
  have lowerBound := GoLean.Sort.searchInterval_partition
    (fun i => !less (xs[i]?.getD fallback) target) 0 xs.length cut (Nat.zero_le _) bounded
    (fun i _ hi => by simpa using partition i hi)
  simp [binarySearch, lowerBound]


end GoLean.Slices
