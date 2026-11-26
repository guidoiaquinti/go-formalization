namespace GoLean.Sort

/-- Go 1.27.0 src/sort/sort.go:127 checks destinations in index order before either write.
A failed bounds check leaves the input list untouched. -/
def indexedSwap (xs : List α) (i j : Int) : Except Int (List α) := do
  let some a := if i < 0 then none else xs[i.toNat]?
    | throw i
  let some b := if j < 0 then none else xs[j.toNat]?
    | throw j
  return (xs.set i.toNat b).set j.toNat a

theorem indexedSwap_spec (fallback : α) (xs : List α) (i j : Int) :
    indexedSwap xs i j =
      if i < 0 ∨ xs.length ≤ i.toNat then .error i
      else if j < 0 ∨ xs.length ≤ j.toNat then .error j
      else .ok ((xs.set i.toNat (xs[j.toNat]?.getD fallback)).set j.toNat (xs[i.toNat]?.getD fallback)) := by
  by_cases hi : i < 0 <;> by_cases hib : i.toNat < xs.length <;>
    by_cases hj : j < 0 <;> by_cases hjb : j.toNat < xs.length <;>
    simp [indexedSwap, hi, hib, hj, hjb, ← Nat.not_lt] <;> rfl

theorem indexedSwap_index_spec (fallback : α) (xs : List α) (i j : Int)
    (hi : 0 ≤ i) (hj : 0 ≤ j) (hib : i.toNat < xs.length) (hjb : j.toNat < xs.length)
    (output : List α) (h : indexedSwap xs i j = .ok output) (k : Nat) :
    output.length = xs.length ∧
    output[k]? = xs[if k = j.toNat then i.toNat else if k = i.toNat then j.toNat else k]? := by
  rw [indexedSwap_spec fallback] at h
  simp [Int.not_lt.mpr hi, Int.not_lt.mpr hj, Nat.not_le.mpr hib, Nat.not_le.mpr hjb] at h
  subst output
  by_cases kj : k = j.toNat
  · subst k
    simp [hjb, hib]
  · by_cases ki : k = i.toNat
    · subst k
      simp [kj, Ne.symm kj, hib, hjb]
    · simp [Ne.symm kj, Ne.symm ki, kj, ki]

end GoLean.Sort
