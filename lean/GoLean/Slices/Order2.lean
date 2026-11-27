import Init.Data.BitVec.Basic

namespace GoLean.Slices

inductive Order2Failure (ε : Type) where
  | nilCounter
  | nilComparator
  | callback (payload : ε)

structure Order2Observation (α ε : Type) where
  result : Except (Order2Failure ε) (Nat × Nat)
  counter : Option Int
  calls : List (α × α)

private def reorderPair (width a b : Nat) (counter : Option Int) (inverted : Bool) :
    Except (Order2Failure ε) (Nat × Nat) × Option Int :=
  if inverted then
    match counter with
    | none => (.error .nilCounter, none)
    | some old => (.ok (b, a), some (BitVec.ofInt width (old + 1)).toInt)
  else (.ok (a, b), counter)

/-- Go 1.27.0 src/slices/zsortordered.go:300–306. Reads the two supplied
    values, tests the reversed comparison and changes the counter only when
    exchanging the returned indices. The data is never changed. -/
def order2Ordered (width a b : Nat) (left right : α) (counter : Option Int)
    (less : α → α → Bool) : Order2Observation α ε :=
  let (result, count) := reorderPair width a b counter (less right left)
  ⟨result, count, []⟩

/-- Go 1.27.0 src/slices/zsortanyfunc.go:298–304. A callable comparator
    receives (data[b], data[a]) once, before reading the counter pointer. -/
def order2CmpFunc (width a b : Nat) (left right : α) (counter : Option Int)
    (compare : Option (α → α → Except ε Int)) : Order2Observation α ε :=
  match compare with
  | none => ⟨.error .nilComparator, counter, []⟩
  | some callback =>
    let calls := [(right, left)]
    match callback right left with
    | .error payload => ⟨.error (.callback payload), counter, calls⟩
    | .ok reply =>
      let (result, count) := reorderPair width a b counter (decide (reply < 0))
      ⟨result, count, calls⟩

theorem order2Ordered_equal (width a b : Nat) (left right : α) (counter : Option Int)
    (less : α → α → Bool) (keeps : less right left = false) :
    order2Ordered (ε := ε) width a b left right counter less =
      ⟨.ok (a, b), counter, []⟩ := by simp [order2Ordered, reorderPair, keeps]

theorem order2Ordered_inverted (width a b : Nat) (left right : α) (counter : Int)
    (less : α → α → Bool) (inverted : less right left = true) :
    order2Ordered (ε := ε) width a b left right (some counter) less =
      ⟨.ok (b, a), some (BitVec.ofInt width (counter + 1)).toInt, []⟩ := by
  simp [order2Ordered, reorderPair, inverted]

theorem order2CmpFunc_panic (width a b : Nat) (left right : α) (counter : Option Int)
    (compare : α → α → Except ε Int) (payload : ε) (fails : compare right left = .error payload) :
    order2CmpFunc width a b left right counter (some compare) =
      ⟨.error (.callback payload), counter, [(right, left)]⟩ := by simp [order2CmpFunc, fails]

theorem order2CmpFunc_nil (width a b : Nat) (left right : α) (counter : Option Int) :
    order2CmpFunc (ε := ε) width a b left right counter none =
      ⟨.error .nilComparator, counter, []⟩ := rfl

end GoLean.Slices
