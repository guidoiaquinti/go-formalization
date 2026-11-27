import Init.Data.BitVec.Lemmas

namespace GoLean.Strings

private abbrev hasSuffixCondition (input suffix : List (BitVec 8)) : Prop :=
  suffix.length ≤ input.length ∧ input.drop (input.length - suffix.length) = suffix

/-- Go compares the matching byte range only after checking its length.
    https://github.com/golang/go/blob/go1.27.0/src/strings/strings.go#L506-L508 -/
def hasSuffix (input suffix : List (BitVec 8)) : Bool :=
  decide (hasSuffixCondition input suffix)

/-- Go checks the length and matches the bytes after the boundary.
    https://github.com/golang/go/blob/go1.27.0/src/strings/strings.go#L1283-L1289
    https://github.com/golang/go/blob/go1.27.0/src/internal/stringslite/strings.go#L20-L22
    https://github.com/golang/go/blob/go1.27.0/src/internal/stringslite/strings.go#L94-L99 -/
def cutSuffix (input suffix : List (BitVec 8)) : List (BitVec 8) × Bool :=
  let boundary := input.length - suffix.length
  if hasSuffixCondition input suffix then
    (input.take boundary, true)
  else
    (input, false)

theorem cutSuffix_spec (input suffix : List (BitVec 8)) :
    (∀ before, cutSuffix input suffix = (before, true) ↔ input = before ++ suffix) ∧
    ((cutSuffix input suffix).2 = false → (cutSuffix input suffix).1 = input) := by
  constructor
  · intro before
    constructor
    · intro result
      simp only [cutSuffix] at result
      split at result
      next h =>
        have lead : input.take (input.length - suffix.length) = before :=
          congrArg Prod.fst result
        calc
          input = input.take (input.length - suffix.length) ++
              input.drop (input.length - suffix.length) :=
            (List.take_append_drop (input.length - suffix.length) input).symm
          _ = before ++ suffix := by rw [lead, h.2]
      next _ => simp at result
    · intro h
      subst input
      simp [cutSuffix, hasSuffixCondition]
  · simp only [cutSuffix]
    split <;> simp

theorem hasSuffix_spec (input suffix : List (BitVec 8)) :
    hasSuffix input suffix = true ↔ ∃ before, input = before ++ suffix := by
  rw [hasSuffix, decide_eq_true_eq]
  change (suffix.length ≤ input.length ∧ input.drop (input.length - suffix.length) = suffix) ↔ ∃ before, input = before ++ suffix
  constructor
  · rintro ⟨_, hmatch⟩
    let boundary := input.length - suffix.length
    refine ⟨input.take boundary, ?_⟩
    calc
      input = input.take boundary ++ input.drop boundary :=
        (List.take_append_drop boundary input).symm
      _ = input.take boundary ++ suffix := congrArg (input.take boundary ++ ·) hmatch
  · rintro ⟨before, rfl⟩
    simp

/-- Go returns the sliced bytes when the affix matches.
    https://github.com/golang/go/blob/go1.27.0/src/strings/strings.go#L1126-L1128
    https://github.com/golang/go/blob/go1.27.0/src/internal/stringslite/strings.go#L108-L113 -/
def trimSuffix (input suffix : List (BitVec 8)) : List (BitVec 8) :=
  if hasSuffix input suffix then input.take (input.length - suffix.length) else input

theorem trimSuffix_spec (input suffix : List (BitVec 8)) :
    trimSuffix input suffix = (cutSuffix input suffix).1 := by
  simp only [cutSuffix]
  by_cases h : hasSuffixCondition input suffix
  · rw [ite_eq_left h]
    simp [trimSuffix, hasSuffix, h]
  · rw [ite_eq_right h]
    simp [trimSuffix, hasSuffix, h]

end GoLean.Strings
