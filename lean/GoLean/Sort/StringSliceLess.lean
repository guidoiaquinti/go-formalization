import GoLean.Sort.IndexedLess
import Init.Data.List.Lex

namespace GoLean.Sort

/-- Go 1.27.0 src/sort/sort.go:155-159. Go strings are byte sequences. -/
def stringSliceLess (xs : List (List UInt8)) (i j : Int) : Except Int Bool :=
  indexedLess (fun a b => a.lex b) xs i j

theorem stringSliceLess_spec (xs : List (List UInt8)) (i j : Int) :
    stringSliceLess xs i j =
      if i < 0 ∨ xs.length ≤ i.toNat then .error i
      else if j < 0 ∨ xs.length ≤ j.toNat then .error j
      else .ok (decide (xs[i.toNat]?.getD [] < xs[j.toNat]?.getD [])) := by
  have lexDecide (a b : List UInt8) : a.lex b = decide (a < b) := by
    exact Bool.eq_iff_iff.mpr (by simp)
  rw [stringSliceLess, indexedLess_spec (fun a b : List UInt8 => a.lex b (fun x y => decide (x < y))) [] xs i j]
  simp only [lexDecide]

end GoLean.Sort
