import GoLean.Slices.All

namespace GoLean.Slices

/-- Go 1.27.0 src/slices/iter.go:24–34. Length is captured at construction
    and each index is read from the backing state when it is visited. -/
def sliceBackward (length : Nat) (read : σ → Nat → α)
    (yield : Option (Nat → α → σ → σ × Except ε Bool)) (state : σ) :
    List (Nat × α) × σ × Except (AllPanic ε) Unit :=
  iterateIndices read yield state (List.range length).reverse

theorem sliceBackward_full (length : Nat) (read : σ → Nat → α) (state : σ) :
    sliceBackward (ε := ε) length read (some fun _ _ s => (s, .ok true)) state =
      (((List.range length).map (fun i => (i, read state i))).reverse, state, .ok ()) := by
  simpa [sliceBackward, iterateIndices, List.map_reverse] using
    (iterateIndices_prefix (ε := ε) read (fun _ _ s => (s, .ok true)) state
      (List.range length).reverse [] (by simp))

theorem sliceBackward_nil (length : Nat) (read : σ → Nat → α) (state : σ) :
    sliceBackward (ε := ε) length read none state =
      ([], state, if length = 0 then .ok () else .error .nilYield) := by
  cases length with
  | zero => rfl
  | succ length => simp [sliceBackward, List.range_succ, List.reverse_append, iterateIndices]

theorem sliceBackward_firstStop (last : Nat) (read : σ → Nat → α)
    (yield : Nat → α → σ → σ × Except ε Bool) (state updated : σ)
    (result : Except ε Bool) (stop : result ≠ .ok true)
    (first : yield last (read state last) state = (updated, result)) :
    sliceBackward (last + 1) read (some yield) state =
      ([(last, read state last)], updated, (result.map fun _ => ()).mapError AllPanic.callback) := by
  cases result with
  | error payload =>
    simp [sliceBackward, List.range_succ, List.reverse_append, iterateIndices, first, Except.map,
      Except.mapError]
  | ok value =>
    cases value with
    | false =>
      simp [sliceBackward, List.range_succ, List.reverse_append, iterateIndices, first, Except.map,
        Except.mapError]
    | true => exact (stop rfl).elim

end GoLean.Slices
