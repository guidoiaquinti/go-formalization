namespace GoLean.Sort

/-- Read both indexed operands in Go evaluation order. -/
def indexedLess (less : α → α → Bool) (xs : List α) (i j : Int) : Except Int Bool := do
  let some a := if i < 0 then none else xs[i.toNat]?
    | throw i
  let some b := if j < 0 then none else xs[j.toNat]?
    | throw j
  return less a b

theorem indexedLess_spec (less : α → α → Bool) (fallback : α) (xs : List α) (i j : Int) :
    indexedLess less xs i j =
      if i < 0 ∨ xs.length ≤ i.toNat then .error i
      else if j < 0 ∨ xs.length ≤ j.toNat then .error j
      else .ok (less (xs[i.toNat]?.getD fallback) (xs[j.toNat]?.getD fallback)) := by
  by_cases hi : i < 0 <;> by_cases hib : i.toNat < xs.length <;>
    by_cases hj : j < 0 <;> by_cases hjb : j.toNat < xs.length <;>
    simp [indexedLess, hi, hib, hj, hjb, ← Nat.not_lt] <;> rfl

end GoLean.Sort
