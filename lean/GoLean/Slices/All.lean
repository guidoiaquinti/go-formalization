namespace GoLean.Slices

inductive AllPanic (ε : Type) where
  | nilYield
  | callback (payload : ε)

def iterateIndices (read : σ → Nat → α)
    (yield : Option (Nat → α → σ → σ × Except ε Bool)) (state : σ) :
    List Nat → List (Nat × α) × σ × Except (AllPanic ε) Unit
  | [] => ([], state, .ok ())
  | i :: tail =>
    match yield with
    | none => ([], state, .error .nilYield)
    | some callback =>
      let v := read state i
      let (updated, result) := callback i v state
      match result with
      | .error payload => ([(i, v)], updated, .error (.callback payload))
      | .ok false => ([(i, v)], updated, .ok ())
      | .ok true =>
        let (trace, finalState, finished) := iterateIndices read yield updated tail
        ((i, v) :: trace, finalState, finished)

/-- Go 1.27.0 src/slices/iter.go:12–22. Capture the length, then read
    each value when the returned iterator reaches its index. Callback state
    can include changes to cells in the captured backing array. -/
def sliceAll (length : Nat) (read : σ → Nat → α)
    (yield : Option (Nat → α → σ → σ × Except ε Bool)) (state : σ) :
    List (Nat × α) × σ × Except (AllPanic ε) Unit :=
  iterateIndices read yield state (List.range length)

theorem iterateIndices_prefix (read : σ → Nat → α)
    (yield : Nat → α → σ → σ × Except ε Bool) (state : σ) (before after : List Nat)
    (keepsGoing : ∀ i ∈ before, yield i (read state i) state = (state, .ok true)) :
    let (trace, finalState, result) := iterateIndices read (some yield) state after
    iterateIndices read (some yield) state (before ++ after) =
      (before.map (fun i => (i, read state i)) ++ trace, finalState, result) := by
  induction before with
  | nil => simp
  | cons i tail ih =>
    have next := keepsGoing i (by simp)
    have rest : ∀ j ∈ tail, yield j (read state j) state = (state, .ok true) :=
      fun j member => keepsGoing j (List.mem_cons_of_mem i member)
    simp only [List.cons_append, iterateIndices, next, List.map_cons, ih rest]

theorem sliceAll_full (length : Nat) (read : σ → Nat → α) (state : σ) :
    sliceAll (ε := ε) length read (some fun _ _ s => (s, .ok true)) state =
      ((List.range length).map (fun i => (i, read state i)), state, .ok ()) := by
  simpa [sliceAll, iterateIndices] using
    (iterateIndices_prefix (ε := ε) read (fun _ _ s => (s, .ok true)) state (List.range length) [] (by simp))

theorem sliceAll_nil (length : Nat) (read : σ → Nat → α) (state : σ) :
    sliceAll (ε := ε) length read none state =
      ([], state, if length = 0 then .ok () else .error .nilYield) := by
  cases length with
  | zero => rfl
  | succ n => simp [sliceAll, List.range_succ_eq_map, iterateIndices]

theorem sliceAll_firstStop (length : Nat) (read : σ → Nat → α)
    (yield : Nat → α → σ → σ × Except ε Bool) (state updated : σ)
    (result : Except ε Bool) (stop : result ≠ .ok true)
    (first : yield 0 (read state 0) state = (updated, result)) :
    sliceAll (length + 1) read (some yield) state =
      ([(0, read state 0)], updated, (result.map fun _ => ()).mapError AllPanic.callback) := by
  cases result with
  | error payload => simp [sliceAll, List.range_succ_eq_map, iterateIndices, first, Except.map, Except.mapError]
  | ok value =>
    cases value with
    | false => simp [sliceAll, List.range_succ_eq_map, iterateIndices, first, Except.map, Except.mapError]
    | true => exact (stop rfl).elim

end GoLean.Slices
