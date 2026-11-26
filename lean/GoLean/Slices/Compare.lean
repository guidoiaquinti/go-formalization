import Init.Data.List.Lex

namespace GoLean.Slices

/-- Go 1.27.0 src/slices/slices.go:50-71. Compare corresponding elements
until one differs or a slice ends. -/
def sliceCompare (cmp : α → α → Int) : List α → List α → Int
  | [], [] => 0
  | [], _ :: _ => -1
  | _ :: _, [] => 1
  | x :: xs, y :: ys =>
    let order := cmp x y
    if order = 0 then sliceCompare cmp xs ys else order

/-- Go's cmp.Compare specialized to int (src/cmp/cmp.go:40-60). -/
def compareInt (x y : Int) : Int :=
  if x < y then -1 else if x > y then 1 else 0

theorem sliceCompare_int_spec (xs ys : List Int) :
    (sliceCompare compareInt xs ys = -1 ↔ xs < ys) ∧
    (sliceCompare compareInt xs ys = 0 ↔ xs = ys) ∧
    (sliceCompare compareInt xs ys = 1 ↔ ys < xs) := by
  induction xs generalizing ys with
  | nil => cases ys <;> simp [sliceCompare]
  | cons x xs ih =>
    cases ys with
    | nil => simp [sliceCompare]
    | cons y ys =>
      by_cases lt : x < y
      · simp [sliceCompare, compareInt, List.cons_lt_cons_iff, lt,
          show ¬y < x by omega, show x ≠ y by omega, show y ≠ x by omega]
      · by_cases gt : y < x
        · simp [sliceCompare, compareInt, List.cons_lt_cons_iff, lt, gt,
            show x ≠ y by omega]
        · have eq : x = y := by omega
          subst y
          simpa [sliceCompare, compareInt, List.cons_lt_cons_iff] using ih ys

end GoLean.Slices
