import GoLean.Slices.Compare
import Init.Data.List.Find

namespace GoLean.Slices

/-- Go 1.27.0 src/slices/slices.go:73-92. Count calls up to and including
the first nonzero or panicking comparison. -/
def sliceCompareFunc (cmp : α → β → Except Unit Int) :
    List α → List β → Except Unit Int × Nat
  | [], [] => (.ok 0, 0)
  | [], _ :: _ => (.ok (-1), 0)
  | _ :: _, [] => (.ok 1, 0)
  | x :: xs, y :: ys =>
    match cmp x y with
    | .error () => (.error (), 1)
    | .ok order =>
      if order != 0 then (.ok order, 1)
      else
        let (result, calls) := sliceCompareFunc cmp xs ys
        (result, calls + 1)

theorem sliceCompareFunc_spec (cmp : α → β → Except Unit Int) (xs : List α) (ys : List β) :
    (sliceCompareFunc cmp xs ys).1 =
      match (xs.zip ys).find? (fun (a, b) => match cmp a b with
        | .ok order => order != 0
        | .error () => true) with
      | some (a, b) => cmp a b
      | none => .ok (compareInt (Int.ofNat xs.length) (Int.ofNat ys.length)) := by
  induction xs generalizing ys with
  | nil => cases ys <;> simp [sliceCompareFunc, compareInt]
  | cons x xs ih =>
    cases ys with
    | nil => simp [sliceCompareFunc, compareInt] <;> omega
    | cons y ys =>
      cases h : cmp x y with
      | error e => cases e; simp [sliceCompareFunc, h]
      | ok order =>
        by_cases zero : order = 0
        · simpa [sliceCompareFunc, h, zero, compareInt] using ih ys
        · simp [sliceCompareFunc, h, zero]

end GoLean.Slices
