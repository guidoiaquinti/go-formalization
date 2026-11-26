import GoLean.Sort.SearchInterval
import Init.Data.Float.Float

namespace GoLean.Sort

/-- Go 1.27.0 src/sort/search.go:127-133. The loop accesses only valid indices. -/
def searchFloat64s (xs : List Float) (key : Float) : Nat :=
  searchInterval (fun i => decide (key ≤ xs[i]?.getD 0)) 0 xs.length

theorem searchFloat64s_spec (xs : List Float) (key : Float) (cut : Nat)
    (bounded : cut ≤ xs.length)
    (partition : ∀ i < xs.length, key ≤ xs[i]?.getD 0 ↔ cut ≤ i) :
    searchFloat64s xs key = cut := by
  apply searchInterval_partition _ 0 xs.length cut (Nat.zero_le _) bounded
  intro i _ hi
  simpa using partition i hi


end GoLean.Sort
