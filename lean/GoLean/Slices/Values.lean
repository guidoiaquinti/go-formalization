import GoLean.Slices.All

namespace GoLean.Slices

/-- Go 1.27.0 src/slices/iter.go:37–45. Traverse the captured length,
    loading each value from the current backing state before yielding it. -/
def sliceValues (length : Nat) (read : σ → Nat → α)
    (yield : Option (α → σ → σ × Except ε Bool)) (state : σ) :
    List α × σ × Except (AllPanic ε) Unit :=
  let (indexed, final, outcome) := sliceAll length read (yield.map fun callback _ => callback) state
  (indexed.map Prod.snd, final, outcome)

theorem sliceValues_nil (length : Nat) (read : σ → Nat → α) (state : σ) :
    sliceValues (ε := ε) length read none state =
      ([], state, if length = 0 then .ok () else .error .nilYield) := by
  simp [sliceValues, sliceAll_nil]

theorem sliceValues_full (length : Nat) (read : σ → Nat → α) (state : σ) :
    sliceValues (ε := ε) length read (some fun _ s => (s, .ok true)) state =
      ((List.range length).map (read state), state, .ok ()) := by
  simp [sliceValues, sliceAll_full, List.map_map, Function.comp_def]

theorem sliceValues_firstStop (length : Nat) (read : σ → Nat → α)
    (yield : α → σ → σ × Except ε Bool) (state final : σ) (result : Except ε Bool)
    (stop : result ≠ .ok true) (first : yield (read state 0) state = (final, result)) :
    sliceValues (length + 1) read (some yield) state =
      ([read state 0], final, (result.map fun _ => ()).mapError AllPanic.callback) := by
  unfold sliceValues
  simp only [Option.map_some]
  rw [sliceAll_firstStop length read (fun _ => yield) state final result stop first]
  rfl

end GoLean.Slices
