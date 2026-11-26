import GoLean.Sort.IndexedLess

namespace GoLean.Sort

/-- Go 1.27.0 src/sort/sort.go:122-126. An error carries the index that
failed its bounds check. Slice elements and indices are signed integers. -/
def intSliceLess (xs : List Int) (i j : Int) : Except Int Bool :=
  indexedLess (fun a b => decide (a < b)) xs i j

theorem intSliceLess_spec (xs : List Int) (i j : Int) :
    intSliceLess xs i j =
      if i < 0 ∨ xs.length ≤ i.toNat then .error i
      else if j < 0 ∨ xs.length ≤ j.toNat then .error j
      else .ok (decide (xs[i.toNat]?.getD 0 < xs[j.toNat]?.getD 0)) :=
  indexedLess_spec _ 0 xs i j

end GoLean.Sort
