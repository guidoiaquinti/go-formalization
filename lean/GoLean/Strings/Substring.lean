import GoLean.Strings.Count

namespace GoLean.Strings

/-- Compare candidate starts in byte order. The source optimizations select
    candidates using length checks, byte search, and rolling hashes.
    https://github.com/golang/go/blob/go1.27.0/src/strings/strings.go#L1261-L1265
    https://github.com/golang/go/blob/go1.27.0/src/internal/stringslite/strings.go#L28-L78 -/
def index (input needle : List Nat) : Option Nat :=
  if needle.isEmpty then some 0 else countFirstMatch needle input

/-- https://github.com/golang/go/blob/go1.27.0/src/strings/strings.go#L61-L64 -/
def contains (input needle : List Nat) : Bool :=
  (index input needle).isSome

/-- https://github.com/golang/go/blob/go1.27.0/src/strings/strings.go#L1267-L1273
    https://github.com/golang/go/blob/go1.27.0/src/internal/stringslite/strings.go#L80-L85 -/
def cut (input needle : List Nat) : List Nat × List Nat × Bool :=
  match index input needle with
  | some i => (input.take i, input.drop (i + needle.length), true)
  | none => (input, [], false)

/-- A candidate matches when the exact byte window beginning there equals
    the requested needle. -/
def matchesAt (input needle : List Nat) (i : Nat) : Prop :=
  (input.drop i).take needle.length = needle

def FirstMatch (input needle : List Nat) : Option Nat → Prop
  | none => ∀ i, ¬ matchesAt input needle i
  | some i => matchesAt input needle i ∧ ∀ j, j < i → ¬ matchesAt input needle j

private theorem scan_first (input needle : List Nat) (nonempty : needle ≠ []) :
    FirstMatch input needle (countFirstMatch needle input) := by
  induction input with
  | nil => simp [FirstMatch, countFirstMatch, matchesAt, Ne.symm nonempty]
  | cons b rest ih =>
    by_cases head : matchesAt (b :: rest) needle 0
    · have same : (b :: rest).take needle.length = needle := head
      simpa [countFirstMatch, same, FirstMatch] using head
    · have diff : (b :: rest).take needle.length ≠ needle := head
      simp only [countFirstMatch, beq_eq_false_iff_ne.mpr diff, Bool.false_eq_true, ↓reduceIte]
      cases outcome : countFirstMatch needle rest with
      | none =>
        simp only [outcome, FirstMatch] at ih
        intro i
        cases i with
        | zero => exact head
        | succ i => exact ih i
      | some k =>
        simp only [outcome, FirstMatch] at ih
        refine ⟨ih.1, ?_⟩
        intro j prior
        cases j with
        | zero => exact head
        | succ j => exact ih.2 j (by omega)

theorem index_spec (input needle : List Nat) :
    FirstMatch input needle (index input needle) := by
  cases needle with
  | nil => simp [index, FirstMatch, matchesAt]
  | cons b rest => exact scan_first input (b :: rest) (by simp)

theorem contains_spec (input needle : List Nat) :
    contains input needle = true ↔ ∃ i, matchesAt input needle i := by
  have first := index_spec input needle
  cases result : index input needle with
  | none =>
    simp only [result, FirstMatch] at first
    simp [contains, result, first]
  | some i =>
    simp only [result, FirstMatch] at first
    simp [contains, result, show ∃ i, matchesAt input needle i from ⟨i, first.1⟩]

theorem cut_spec (input needle : List Nat) :
    match index input needle with
    | none => cut input needle = (input, [], false)
    | some i =>
      cut input needle = (input.take i, input.drop (i + needle.length), true) ∧
      input = (input.take i) ++ needle ++ (input.drop (i + needle.length)) := by
  have first := index_spec input needle
  cases result : index input needle with
  | none => simp [cut, result]
  | some i =>
    simp only [result, FirstMatch] at first
    refine ⟨by simp [cut, result], ?_⟩
    have suffix : input.drop i = needle ++ input.drop (i + needle.length) := by
      rw [← List.drop_drop]
      exact (List.take_append_drop needle.length (input.drop i)).symm.trans
        (congrArg (· ++ (input.drop i).drop needle.length) first.1)
    rw [List.append_assoc, ← suffix, List.take_append_drop]

end GoLean.Strings
