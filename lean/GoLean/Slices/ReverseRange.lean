import GoLean.Slices.Reverse

namespace GoLean.Slices

private def reverseSubrange (cells : List α) (a b : Nat) : List α :=
  cells.take a ++ sliceReverse ((cells.drop a).take (b - a)) ++ cells.drop b

/-- Go 1.27.0 src/slices/zsortordered.go:321–329. The legal half-open
    subrange is reversed by moving endpoints; other cells are preserved. -/
def reverseRangeOrdered (cells : Option (List α)) (a b : Nat) : Option (List α) :=
  cells.map fun values => reverseSubrange values a b

/-- Go 1.27.0 src/slices/zsortanyfunc.go:319–327. The comparator is unused. -/
def reverseRangeCmpFunc (cells : Option (List α)) (a b : Nat)
    (_ : Option (α → α → Int)) : Option (List α) :=
  cells.map fun values => reverseSubrange values a b

theorem reverseRangeOrdered_spec (cells : List α) (a b : Nat) :
    reverseRangeOrdered (some cells) a b =
      some (cells.take a ++ ((cells.drop a).take (b - a)).reverse ++ cells.drop b) := by
  simp [reverseRangeOrdered, reverseSubrange, sliceReverse_spec]

theorem reverseRangeOrdered_empty (cells : Option (List α)) (a : Nat) :
    reverseRangeOrdered cells a a = cells := by
  cases cells <;> simp [reverseRangeOrdered, reverseSubrange, sliceReverse, List.take_append_drop]

theorem reverseRangeCmpFunc_ignoresComparator (cells : Option (List α)) (a b : Nat)
    (first second : Option (α → α → Int)) :
    reverseRangeCmpFunc cells a b first = reverseRangeCmpFunc cells a b second := rfl

theorem reverseRangeCmpFunc_spec (cells : Option (List α)) (a b : Nat)
    (compare : Option (α → α → Int)) :
    reverseRangeCmpFunc cells a b compare = reverseRangeOrdered cells a b := rfl

end GoLean.Slices
