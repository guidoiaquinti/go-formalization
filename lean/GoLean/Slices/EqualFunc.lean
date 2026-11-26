import GoLean.Slices.IndexFunc

namespace GoLean.Slices

/-- Go 1.27.0 src/slices/slices.go:32-48 checks length first, then
scans aligned pairs until the equality callback fails or panics. -/
def sliceEqualFunc (eq : α → β → Except Unit Bool) (xs : List α) (ys : List β) :
    Except Unit Bool × Nat :=
  if xs.length != ys.length then (.ok false, 0)
  else
    let (first, calls) := sliceIndexFunc (fun (a, b) => (eq a b).map (! ·)) (xs.zip ys)
    (first.map Option.isNone, calls)

theorem sliceEqualFunc_spec (eq : α → β → Bool) (xs : List α) (ys : List β) :
    (sliceEqualFunc (fun a b => .ok (eq a b)) xs ys).1 = .ok true ↔
      xs.length = ys.length ∧ ∀ p ∈ xs.zip ys, eq p.1 p.2 = true := by
  by_cases h : xs.length = ys.length <;>
    simp [sliceEqualFunc, h, Except.map, sliceIndexFunc_pure, List.findIdx?_isNone, List.all_eq_true]

theorem sliceEqualFunc_length_mismatch (eq : α → β → Except Unit Bool)
    (xs : List α) (ys : List β) (h : xs.length ≠ ys.length) :
    sliceEqualFunc eq xs ys = (.ok false, 0) := by
  simp [sliceEqualFunc, h]


end GoLean.Slices
