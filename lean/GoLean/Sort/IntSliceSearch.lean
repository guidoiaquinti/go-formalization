import GoLean.Sort.SearchInts

namespace GoLean.Sort

/-- Go 1.27.0 src/sort/search.go:143-144 delegates to SearchInts. -/
def intSliceSearch (xs : List Int) (key : Int) : Nat :=
  searchInts xs key

theorem intSliceSearch_spec (xs : List Int) (key : Int) (cut : Nat)
    (bounded : cut ≤ xs.length)
    (partition : ∀ i < xs.length, key ≤ xs[i]?.getD 0 ↔ cut ≤ i) :
    intSliceSearch xs key = cut :=
  searchInts_spec xs key cut bounded partition


end GoLean.Sort
