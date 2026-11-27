import GoLean.Strings.Substring

namespace GoLean.Strings

private def nextSplit (input separator : List Nat) (save : Nat) : Option (List Nat × Option (List Nat)) :=
  if separator.isEmpty then
    match input with
    | [] => none
    | lead :: rest =>
      let width := runeWidth lead rest
      some (input.take width, some (input.drop width))
  else
    match index input separator with
    | none => some (input, none)
    | some i => some (input.take (i + save), some (input.drop (i + separator.length)))

/--
Go 1.27.0 strings/iter.go:36–69. The returned state is the captured
suffix. The callback at the budget returns false, before the source
advances that suffix. A final yield also leaves the suffix unchanged.
-/
private def runSplitSeq (input separator : List Nat) (save : Nat) : Nat → List (List Nat) × List Nat
  | 0 => ([], input)
  | budget + 1 =>
    match nextSplit input separator save with
    | none => ([], input)
    | some (part, none) => ([part], input)
    | some (part, some rest) =>
      if budget = 0 then ([part], input)
      else
        let (parts, remaining) := runSplitSeq rest separator save budget
        (part :: parts, remaining)

/-- Go strings/iter.go:67–69 uses a zero-length retained separator. -/
def splitSeq (input separator : List Nat) (budget : Nat) : List (List Nat) × List Nat :=
  runSplitSeq input separator 0 budget

/-- Go strings/iter.go:75–77 retains the complete separator. -/
def splitAfterSeq (input separator : List Nat) (budget : Nat) : List (List Nat) × List Nat :=
  runSplitSeq input separator separator.length budget

/-- Stopping at the first yield preserves the entire captured input. -/
theorem splitSeq_stop_first (input separator : List Nat) :
    (splitSeq input separator 1).2 = input := by
  unfold splitSeq runSplitSeq
  split <;> simp

/-- A nonempty separator yields one final empty fragment on empty input. -/
theorem splitSeq_empty (lead : Nat) (rest : List Nat) (budget : Nat) :
    splitSeq [] (lead :: rest) (budget + 1) = ([[]], []) := by
  simp [splitSeq, runSplitSeq, nextSplit, index, countFirstMatch]

/-- The first match yields through the separator, and a stop leaves it unconsumed. -/
theorem splitAfterSeq_first (input separator : List Nat) (i : Nat)
    (nonempty : separator ≠ []) (found : index input separator = some i) :
    splitAfterSeq input separator 1 = ([input.take i ++ separator], input) := by
  have spec := index_spec input separator
  rw [found] at spec
  have matched : (input.drop i).take separator.length = separator := spec.1
  simp [splitAfterSeq, runSplitSeq, nextSplit, List.isEmpty_iff, nonempty,
    found, List.take_add, matched]

end GoLean.Strings
