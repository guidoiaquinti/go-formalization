import GoLean.Strings.Substring

namespace GoLean.Strings

private def scanLast (needle : List Nat) : List Nat → Option Nat
  | [] => none
  | b :: rest =>
    match scanLast needle rest with
    | some i => some (i + 1)
    | none => if (b :: rest).take needle.length == needle then some 0 else none

/-- Go matches in descending offset order. The scan abstracts the rolling
    hash filter by directly checking each byte window.
    https://github.com/golang/go/blob/go1.27.0/src/strings/strings.go#L81-L99
    https://github.com/golang/go/blob/go1.27.0/src/internal/bytealg/bytealg.go#L95-L117 -/
def lastIndex (input needle : List Nat) : Option Nat :=
  if needle.isEmpty then some input.length else scanLast needle input

/-- The result is a matching offset with no later match. For an empty
    needle every offset through the end is a candidate. -/
def LastMatch (input needle : List Nat) : Option Nat → Prop
  | none => ∀ i, i ≤ input.length → ¬ matchesAt input needle i
  | some i => i ≤ input.length ∧ matchesAt input needle i ∧
      ∀ j, i < j → j ≤ input.length → ¬ matchesAt input needle j

private theorem scanLast_spec (input needle : List Nat) (nonempty : needle ≠ []) :
    LastMatch input needle (scanLast needle input) := by
  induction input with
  | nil => simp [LastMatch, scanLast, matchesAt, Ne.symm nonempty]
  | cons b rest ih =>
    unfold scanLast
    cases result : scanLast needle rest with
    | some i =>
      simp only [result, LastMatch] at ih
      refine ⟨by simp; omega, ih.2.1, ?_⟩
      intro j after bound
      cases j with
      | zero => omega
      | succ j => exact ih.2.2 j (by omega) (by simp at bound; omega)
    | none =>
      simp only [result, LastMatch] at ih
      by_cases candidate : matchesAt (b :: rest) needle 0
      · have same : (b :: rest).take needle.length = needle := candidate
        simp only [beq_iff_eq, same, ↓reduceIte]
        refine ⟨by omega, candidate, ?_⟩
        intro j after bound
        cases j with
        | zero => omega
        | succ j => exact ih j (by simp at bound; omega)
      · have diff : (b :: rest).take needle.length ≠ needle := candidate
        simp only [beq_iff_eq, diff, ↓reduceIte]
        intro j bound
        cases j with
        | zero => exact candidate
        | succ j => exact ih j (by simp at bound; omega)

theorem lastIndex_spec (input needle : List Nat) :
    LastMatch input needle (lastIndex input needle) := by
  cases needle with
  | nil => simp [lastIndex, LastMatch, matchesAt]
  | cons b rest => exact scanLast_spec input (b :: rest) (by simp)

/-- Go slices before and after the last match.
    https://github.com/golang/go/blob/go1.27.0/src/strings/strings.go#L1291-L1300 -/
def cutLast (input needle : List Nat) : List Nat × List Nat × Bool :=
  match lastIndex input needle with
  | none => (input, [], false)
  | some i => (input.take i, input.drop (i + needle.length), true)

theorem cutLast_spec (input needle : List Nat) :
    match lastIndex input needle with
    | none => cutLast input needle = (input, [], false)
    | some i =>
      cutLast input needle = (input.take i, input.drop (i + needle.length), true) ∧
      input = input.take i ++ needle ++ input.drop (i + needle.length) := by
  have last := lastIndex_spec input needle
  cases result : lastIndex input needle with
  | none => simp [cutLast, result]
  | some i =>
    simp only [result, LastMatch] at last
    refine ⟨by simp [cutLast, result], ?_⟩
    have suffix : input.drop i = needle ++ input.drop (i + needle.length) := by
      rw [← List.drop_drop]
      exact (List.take_append_drop needle.length (input.drop i)).symm.trans
        (congrArg (· ++ (input.drop i).drop needle.length) last.2.1)
    rw [List.append_assoc, ← suffix, List.take_append_drop]

end GoLean.Strings
