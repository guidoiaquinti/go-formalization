import GoLean.Sort.IsSorted
import Init.Data.List.Lex

namespace GoLean.Sort

/-- Go 1.27.0 src/sort/sort.go:197 and src/slices/sort.go:42-49.
The string specialization of cmp.Less compares unsigned bytes lexicographically. -/
def stringsAreSorted (xs : List (List UInt8)) : Bool :=
  isSortedList (fun a b => a.lex b) [] xs

theorem stringsAreSorted_spec (xs : List (List UInt8)) :
    stringsAreSorted xs = true ↔
      ∀ i < xs.length - 1, ¬xs[i + 1]?.getD [] < xs[i]?.getD [] := by
  simp [stringsAreSorted, isSortedList_spec, Bool.eq_false_iff]


end GoLean.Sort
