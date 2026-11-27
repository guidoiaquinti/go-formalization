import GoLean.Strings.Substring

namespace GoLean.Strings

structure ConfiguredSplitRun where
  parts : List (List Nat)
  remaining : List Nat
  boundsPanic : Bool

private def configuredSplitStep (input separator : List Nat) (save : Int) :
    Except Unit (Option (List Nat × Option (List Nat))) :=
  if separator.isEmpty then
    match input with
    | [] => .ok none
    | lead :: rest =>
      let width := runeWidth lead rest
      .ok (some (input.take width, some (input.drop width)))
  else
    match index input separator with
    | none => .ok (some (input, none))
    | some i =>
      let endpoint := Int.ofNat i + save
      if endpoint < 0 || Int.ofNat input.length < endpoint then .error ()
      else .ok (some (input.take endpoint.toNat, some (input.drop (i + separator.length))))

/-- Go 1.27.0 src/strings/iter.go:36–60. A false callback or a slice
    bounds panic leaves that current fragment unconsumed. With legal native
    lengths and save values, a positive endpoint overflow wraps below zero
    and is already classified as a bounds panic by the mathematical check. -/
def splitSeqConfigured (input separator : List Nat) (save : Int) : Nat → ConfiguredSplitRun
  | 0 => ⟨[], input, false⟩
  | steps + 1 =>
    match configuredSplitStep input separator save with
    | .error _ => ⟨[], input, true⟩
    | .ok none => ⟨[], input, false⟩
    | .ok (some (part, none)) => ⟨[part], input, false⟩
    | .ok (some (part, some rest)) =>
      if steps = 0 then ⟨[part], input, false⟩
      else
        let result := splitSeqConfigured rest separator save steps
        {result with parts := part :: result.parts}

theorem splitSeqConfigured_stopped (input separator : List Nat) (save : Int) :
    (splitSeqConfigured input separator save 1).remaining = input := by
  unfold splitSeqConfigured
  split <;> simp

theorem splitSeqConfigured_partial :
    (splitSeqConfigured [97, 45, 45, 98, 45, 45] [45, 45] 1 9).parts =
      [[97, 45], [98, 45], []] := rfl

theorem splitSeqConfigured_panicAfterYield :
    let result := splitSeqConfigured [97, 98, 45, 45] [45] (-1) 8
    result.parts = [[97]] ∧ result.remaining = [45] ∧ result.boundsPanic = true :=
  ⟨rfl, rfl, rfl⟩

theorem splitSeqConfigured_emptyIgnoresSave (input : List Nat) (save : Int) (budget : Nat) :
    (splitSeqConfigured input [] save budget).boundsPanic = false := by
  induction budget generalizing input with
  | zero => rfl
  | succ budget ih =>
    cases input with
    | nil => rfl
    | cons lead rest =>
      simp only [splitSeqConfigured, configuredSplitStep, List.isEmpty_nil, ite_true]
      split <;> simp [ih]

end GoLean.Strings
