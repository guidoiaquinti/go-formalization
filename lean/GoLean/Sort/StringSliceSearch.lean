import GoLean.Sort.SearchStrings

namespace GoLean.Sort

/-- Go 1.27.0 src/sort/search.go:149-150 delegates to SearchStrings. -/
def stringSliceSearch (xs : List (List UInt8)) (key : List UInt8) : Nat :=
  searchStrings xs key

theorem stringSliceSearch_spec (xs : List (List UInt8)) (key : List UInt8) (cut : Nat)
    (bounded : cut ≤ xs.length)
    (partition : ∀ i < xs.length, (¬xs[i]?.getD [] < key) ↔ cut ≤ i) :
    stringSliceSearch xs key = cut :=
  searchStrings_spec xs key cut bounded partition


end GoLean.Sort
