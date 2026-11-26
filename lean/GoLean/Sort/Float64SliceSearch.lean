import GoLean.Sort.SearchFloat64s

namespace GoLean.Sort

/-- Go 1.27.0 src/sort/search.go:146-147 delegates to SearchFloat64s. -/
def float64SliceSearch (xs : List Float) (key : Float) : Nat :=
  searchFloat64s xs key

theorem float64SliceSearch_spec (xs : List Float) (key : Float) (cut : Nat)
    (bounded : cut ≤ xs.length)
    (partition : ∀ i < xs.length, key ≤ xs[i]?.getD 0 ↔ cut ≤ i) :
    float64SliceSearch xs key = cut :=
  searchFloat64s_spec xs key cut bounded partition


end GoLean.Sort
