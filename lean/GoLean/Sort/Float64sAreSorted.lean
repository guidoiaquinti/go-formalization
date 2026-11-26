import GoLean.Sort.IsSorted
import Init.Data.Float.Float

namespace GoLean.Sort

/-- Go 1.27.0 src/sort/sort.go:192 delegates to src/slices/sort.go:42-49.
The float64 specialization of src/cmp/cmp.go:28-30 uses self-inequality to detect NaNs. -/
def float64sAreSorted (xs : List Float) : Bool :=
  isSortedList (fun a b => (a != a && !(b != b)) || a.lt b) 0 xs

theorem float64sAreSorted_spec (xs : List Float) :
    float64sAreSorted xs = true ↔
      ∀ i < xs.length - 1,
        let a := xs[i]?.getD 0
        let b := xs[i + 1]?.getD 0
        ((b == b) = false → (a == a) = false) ∧ b.lt a = false := by
  simp [float64sAreSorted, isSortedList_spec, bne, Bool.or_eq_false_iff]


end GoLean.Sort
