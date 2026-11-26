import GoLean.Sort.IsSorted

namespace GoLean.Sort

/-- Go 1.27.0 src/sort/sort.go:186 and src/slices/sort.go:42-49.
For int, cmp.Less reduces to signed less-than. The reverse adjacent scan
shares its control flow with sort.IsSorted. -/
def intsAreSorted (xs : List Int) : Bool :=
  isSortedList (fun a b => decide (a < b)) 0 xs

theorem intsAreSorted_spec (xs : List Int) :
    intsAreSorted xs = true ↔
      ∀ i < xs.length - 1, xs[i]?.getD 0 ≤ xs[i + 1]?.getD 0 := by
  simp [intsAreSorted, isSortedList_spec, Int.not_lt]

end GoLean.Sort
