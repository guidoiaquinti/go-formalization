import GoLean.Sort.SearchInterval

namespace GoLean.Slices

/-- Go 1.27.0 src/slices/sort.go:152–168, with a pure element-to-target comparator. -/
def binarySearchFunc (cmp : α → β → Int) (fallback : α)
    (xs : List α) (target : β) : Nat × Bool :=
  let i := GoLean.Sort.searchInterval
    (fun j => !decide (cmp (xs[j]?.getD fallback) target < 0)) 0 xs.length
  (i, decide (i < xs.length) && (cmp (xs[i]?.getD fallback) target == 0))

theorem binarySearchFunc_spec (cmp : α → β → Int) (fallback : α)
    (xs : List α) (target : β) (cut : Nat) (bounded : cut ≤ xs.length)
    (partition : ∀ i < xs.length, 0 ≤ cmp (xs[i]?.getD fallback) target ↔ cut ≤ i) :
    binarySearchFunc cmp fallback xs target =
      (cut, decide (cut < xs.length) && (cmp (xs[cut]?.getD fallback) target == 0)) := by
  have lowerBound := GoLean.Sort.searchInterval_partition
    (fun i => !decide (cmp (xs[i]?.getD fallback) target < 0))
    0 xs.length cut (Nat.zero_le _) bounded
    (fun i _ hi => by simpa using partition i hi)
  simp [binarySearchFunc, lowerBound]

end GoLean.Slices
