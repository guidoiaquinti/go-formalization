import GoLean.Slices.IndexFunc

namespace GoLean.Slices

/-- Go 1.27.0 src/slices/slices.go:121-126 tests the sign of the IndexFunc result. -/
def sliceContainsFunc (predicate : α → Except Unit Bool) (xs : List α) :
    Except Unit Bool × Nat :=
  let (index, calls) := sliceIndexFunc predicate xs
  (index.map Option.isSome, calls)

theorem sliceContainsFunc_pure (p : α → Bool) (xs : List α) :
    (sliceContainsFunc (fun x => .ok (p x)) xs).1 = .ok true ↔ ∃ x ∈ xs, p x = true := by
  simp [sliceContainsFunc, sliceIndexFunc_pure, Except.map, List.findIdx?_isSome, List.any_eq_true]

theorem sliceContainsFunc_calls (p : α → Except Unit Bool) (xs : List α) :
    (sliceContainsFunc p xs).2 =
      ((xs.findIdx? fun x => match p x with
        | .ok false => false
        | _ => true).map Nat.succ).getD xs.length :=
  sliceIndexFunc_calls p xs


end GoLean.Slices
