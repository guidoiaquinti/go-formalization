import Init.Data.BitVec.Lemmas

namespace GoLean.Strings

private abbrev hasPrefixCondition (input pre : List (BitVec 8)) : Prop :=
  pre.length ≤ input.length ∧ input.take pre.length = pre

/-- Go compares the matching byte range only after checking its length.
    https://github.com/golang/go/blob/go1.27.0/src/strings/strings.go#L501-L503 -/
def hasPrefix (input pre : List (BitVec 8)) : Bool :=
  decide (hasPrefixCondition input pre)

/-- The length guard makes Go's byte slice safe before checking equality.
    https://github.com/golang/go/blob/go1.27.0/src/strings/strings.go#L1276-L1281
    https://github.com/golang/go/blob/go1.27.0/src/internal/stringslite/strings.go#L16-L18
    https://github.com/golang/go/blob/go1.27.0/src/internal/stringslite/strings.go#L87-L92 -/
def cutPrefix (input pre : List (BitVec 8)) : List (BitVec 8) × Bool :=
  if hasPrefixCondition input pre then
    (input.drop pre.length, true)
  else
    (input, false)

theorem cutPrefix_spec (input pre : List (BitVec 8)) :
    (∀ rest, cutPrefix input pre = (rest, true) ↔ input = pre ++ rest) ∧
    ((cutPrefix input pre).2 = false → (cutPrefix input pre).1 = input) := by
  constructor
  · intro rest
    constructor
    · intro result
      unfold cutPrefix at result
      split at result
      next h =>
        have suffix : input.drop pre.length = rest := congrArg Prod.fst result
        calc
          input = input.take pre.length ++ input.drop pre.length :=
            (List.take_append_drop pre.length input).symm
          _ = pre ++ rest := by rw [suffix, h.2]
      next _ => simp at result
    · intro h
      subst input
      simp [cutPrefix, hasPrefixCondition]
  · unfold cutPrefix
    split <;> simp

theorem hasPrefix_spec (input pre : List (BitVec 8)) :
    hasPrefix input pre = true ↔ ∃ rest, input = pre ++ rest := by
  rw [hasPrefix, decide_eq_true_eq]
  change (pre.length ≤ input.length ∧ input.take pre.length = pre) ↔ ∃ rest, input = pre ++ rest
  constructor
  · rintro ⟨_, hmatch⟩
    refine ⟨input.drop pre.length, ?_⟩
    calc
      input = input.take pre.length ++ input.drop pre.length :=
        (List.take_append_drop pre.length input).symm
      _ = pre ++ input.drop pre.length := congrArg (· ++ input.drop pre.length) hmatch
  · rintro ⟨rest, rfl⟩
    simp

/-- Go returns the sliced bytes when the affix matches.
    https://github.com/golang/go/blob/go1.27.0/src/strings/strings.go#L1120-L1122
    https://github.com/golang/go/blob/go1.27.0/src/internal/stringslite/strings.go#L101-L106 -/
def trimPrefix (input pre : List (BitVec 8)) : List (BitVec 8) :=
  if hasPrefix input pre then input.drop pre.length else input

theorem trimPrefix_spec (input pre : List (BitVec 8)) :
    trimPrefix input pre = (cutPrefix input pre).1 := by
  simp only [cutPrefix]
  by_cases h : hasPrefixCondition input pre
  · rw [ite_eq_left h]
    simp [trimPrefix, hasPrefix, h]
  · rw [ite_eq_right h]
    simp [trimPrefix, hasPrefix, h]

end GoLean.Strings
