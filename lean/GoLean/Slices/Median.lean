import GoLean.Slices.Order2

namespace GoLean.Slices

structure MedianObservation (α ε : Type) where
  result : Except (Order2Failure ε) Nat
  counter : Option Int
  calls : List (α × α)

private def median (step : Nat → Nat → Nat → Option Int → Order2Observation α ε)
    (a b c : Nat) (counter : Option Int) : MedianObservation α ε :=
  let first := step 0 a b counter
  match first.result with
  | .error failure => ⟨.error failure, first.counter, first.calls⟩
  | .ok (a, b) =>
    let second := step 1 b c first.counter
    match second.result with
    | .error failure => ⟨.error failure, second.counter, first.calls ++ second.calls⟩
    | .ok (b, _) =>
      let third := step 2 a b second.counter
      ⟨third.result.map Prod.snd, third.counter, first.calls ++ second.calls ++ third.calls⟩

/-- Go 1.27.0 src/slices/zsortordered.go:309–314, with valid indices. -/
def medianOrdered (width : Nat) (get : Nat → α) (a b c : Nat) (counter : Option Int)
    (less : α → α → Bool) : MedianObservation α ε :=
  median (fun _ x y => fun count => order2Ordered width x y (get x) (get y) count less)
    a b c counter

/-- Go 1.27.0 src/slices/zsortanyfunc.go:307–312, with valid indices.
The call index allows the comparator to observe that calls have occurred. -/
def medianCmpFunc (width : Nat) (get : Nat → α) (a b c : Nat) (counter : Option Int)
    (compare : Option (Nat → α → α → Except ε Int)) : MedianObservation α ε :=
  median (fun call x y count => order2CmpFunc width x y (get x) (get y) count
    (compare.map (fun f => f call))) a b c counter

/-- Go 1.27.0 src/slices/zsortordered.go:317–319, with a valid neighborhood. -/
def medianAdjacentOrdered (width : Nat) (get : Nat → α) (a : Nat) (counter : Option Int)
    (less : α → α → Bool) : MedianObservation α ε :=
  medianOrdered width get (a - 1) a (a + 1) counter less

/-- Go 1.27.0 src/slices/zsortanyfunc.go:315–317, with a valid neighborhood. -/
def medianAdjacentCmpFunc (width : Nat) (get : Nat → α) (a : Nat) (counter : Option Int)
    (compare : Option (Nat → α → α → Except ε Int)) : MedianObservation α ε :=
  medianCmpFunc width get (a - 1) a (a + 1) counter compare

theorem medianOrdered_allTies (width a b c : Nat) (get : Nat → α) (counter : Option Int) :
    medianOrdered (ε := ε) width get a b c counter (fun _ _ => false) =
      ⟨.ok b, counter, []⟩ := rfl

theorem medianCmpFunc_failsFirst (width a b c : Nat) (get : Nat → α) (counter : Option Int)
    (payload : ε) :
    medianCmpFunc width get a b c counter (some (fun _ _ _ => .error payload)) =
      ⟨.error (.callback payload), counter, [(get b, get a)]⟩ := rfl

theorem medianCmpFunc_retainsPriorCount :
    medianCmpFunc 64 (fun i => ([30, 20, 10] : List Int)[i]?.getD 0)
      0 1 2 (some 5) (some (fun call _ _ => if call = 2 then .error 2 else .ok (-1))) =
      ⟨.error (.callback (2 : Int)), some 7, [(20, 30), (10, 30), (10, 20)]⟩ := rfl

theorem medianAdjacentOrdered_center (width a : Nat) (get : Nat → α) (counter : Option Int) :
    medianAdjacentOrdered (ε := ε) width get a counter (fun _ _ => false) =
      ⟨.ok a, counter, []⟩ := rfl

theorem medianAdjacentCmpFunc_first (width a : Nat) (get : Nat → α) (counter : Option Int)
    (payload : ε) :
    medianAdjacentCmpFunc width get a counter (some (fun _ _ _ => .error payload)) =
      ⟨.error (.callback payload), counter, [(get a, get (a - 1))]⟩ := rfl

end GoLean.Slices
