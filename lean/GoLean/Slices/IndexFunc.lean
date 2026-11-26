import Init.Data.List.Find

namespace GoLean.Slices

/-- Go 1.27.0 src/slices/slices.go:105-114. The second component counts
predicate calls, including the call that returns true or panics. -/
def sliceIndexFunc (predicate : α → Except Unit Bool) :
    List α → Except Unit (Option Nat) × Nat
  | [] => (.ok none, 0)
  | x :: xs =>
    match predicate x with
    | .error () => (.error (), 1)
    | .ok true => (.ok (some 0), 1)
    | .ok false =>
      let (result, calls) := sliceIndexFunc predicate xs
      (result.map (Option.map Nat.succ), calls + 1)

theorem sliceIndexFunc_pure (predicate : α → Bool) (xs : List α) :
    (sliceIndexFunc (fun x => .ok (predicate x)) xs).1 = .ok (xs.findIdx? predicate) := by
  induction xs with
  | nil => rfl
  | cons x xs ih =>
    cases h : predicate x <;> simp [sliceIndexFunc, List.findIdx?_cons, h, ih, Except.map]

theorem sliceIndexFunc_calls (predicate : α → Except Unit Bool) (xs : List α) :
    (sliceIndexFunc predicate xs).2 =
      ((xs.findIdx? fun x => match predicate x with
        | .ok false => false
        | _ => true).map Nat.succ).getD xs.length := by
  induction xs with
  | nil => rfl
  | cons x xs ih =>
    cases h : predicate x with
    | error e => cases e; simp [sliceIndexFunc, List.findIdx?_cons, h]
    | ok b =>
      cases b
      · simp only [sliceIndexFunc, h, List.findIdx?_cons, Bool.false_eq_true, ite_false]
        rw [ih]
        generalize (xs.findIdx? fun x => match predicate x with | .ok false => false | _ => true) = index
        cases index <;> rfl
      · simp [sliceIndexFunc, List.findIdx?_cons, h]

end GoLean.Slices
