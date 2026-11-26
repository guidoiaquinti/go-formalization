import GoLean.Sort.SearchInterval

namespace GoLean.Sort

/-- Go 1.27.0 src/sort/search.go:119-125. The loop accesses only valid indices. -/
def searchInts (xs : List Int) (key : Int) : Nat :=
  searchInterval (fun i => decide (key ≤ xs[i]?.getD 0)) 0 xs.length

theorem searchInts_spec (xs : List Int) (key : Int) (cut : Nat)
    (bounded : cut ≤ xs.length)
    (partition : ∀ i < xs.length, key ≤ xs[i]?.getD 0 ↔ cut ≤ i) :
    searchInts xs key = cut := by
  apply searchInterval_partition _ 0 xs.length cut (Nat.zero_le _) bounded
  intro i _ hi
  simpa using partition i hi

end GoLean.Sort
