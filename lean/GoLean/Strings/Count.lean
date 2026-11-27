import Init.Data.List.Count
import Init.Data.List.TakeDrop

namespace GoLean.Strings

private def continuation (b : Nat) : Bool :=
  0x80 ≤ b && b ≤ 0xbf

/-- Decode only the width. Invalid or incomplete sequences consume a single byte.
    The lead-byte and continuation bounds reproduce Go's width and acceptance tables.
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf8/utf8.go#L65-L108
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf8/utf8.go#L219-L266 -/
def runeWidth (lead : Nat) (rest : List Nat) : Nat :=
  match rest with
  | [] => 1
  | b :: bs =>
    if 0xc2 ≤ lead && lead ≤ 0xdf && continuation b then 2
    else
      match bs with
      | [] => 1
      | c :: cs =>
        if 0xe0 ≤ lead && lead ≤ 0xef && continuation b &&
            (lead != 0xe0 || b ≥ 0xa0) && (lead != 0xed || b ≤ 0x9f) && continuation c then 3
        else
          match cs with
          | [] => 1
          | d :: _ =>
            if 0xf0 ≤ lead && lead ≤ 0xf4 && continuation b &&
                (lead != 0xf0 || b ≥ 0x90) && (lead != 0xf4 || b ≤ 0x8f) &&
                continuation c && continuation d then 4
            else 1

/-- Count the iterations of Go's string range over raw UTF-8 bytes.
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf8/utf8.go#L447-L452
    https://github.com/golang/go/blob/go1.27.0/src/runtime/utf8.go#L54-L98 -/
def countRunes : List Nat → Nat
  | [] => 0
  | lead :: rest => 1 + countRunes (rest.drop (runeWidth lead rest - 1))
termination_by xs => xs.length
decreasing_by simp [List.length_drop]; omega

/-- Observable result of Index on a nonempty pattern. Search strategies used
    to find the first match are abstracted by a direct candidate scan.
    https://github.com/golang/go/blob/go1.27.0/src/strings/strings.go#L1263-L1265
    https://github.com/golang/go/blob/go1.27.0/src/internal/stringslite/strings.go#L28-L78 -/
def countFirstMatch (needle : List Nat) : List Nat → Option Nat
  | [] => none
  | b :: bs =>
    if (b :: bs).take needle.length == needle then some 0
    else (countFirstMatch needle bs).map Nat.succ

/-- Count starts after the end of each successful match.
    https://github.com/golang/go/blob/go1.27.0/src/strings/strings.go#L37-L59 -/
def countMatches (needle : List Nat) (positive : 0 < needle.length) : List Nat → Nat
  | [] => 0
  | b :: bs =>
    match countFirstMatch needle (b :: bs) with
    | none => 0
    | some i => 1 + countMatches needle positive ((b :: bs).drop (i + needle.length))
termination_by xs => xs.length
decreasing_by simp [List.length_drop]; omega

/-- Go treats empty, one-byte and longer separators separately.
    https://github.com/golang/go/blob/go1.27.0/src/internal/bytealg/count_generic.go#L19-L27
    https://github.com/golang/go/blob/go1.27.0/src/strings/strings.go#L42-L59 -/
def count (text : List Nat) : (needle : List Nat) → Nat
  | [] => countRunes text + 1
  | [b] => text.count b
  | b :: c :: cs => countMatches (b :: c :: cs) (by simp) text

/-- An independent position-by-position specification. A match consumes its
    whole separator, while a mismatch consumes one byte. -/
def greedyCount (needle : List Nat) (positive : 0 < needle.length) : List Nat → Nat
  | [] => 0
  | b :: bs =>
    if (b :: bs).take needle.length = needle then
      1 + greedyCount needle positive ((b :: bs).drop needle.length)
    else greedyCount needle positive bs
termination_by xs => xs.length
decreasing_by all_goals simp [List.length_drop] <;> omega

private theorem countMatches_step (needle : List Nat) (positive : 0 < needle.length)
    (b : Nat) (bs : List Nat) :
    countMatches needle positive (b :: bs) =
      if (b :: bs).take needle.length = needle then
        1 + countMatches needle positive ((b :: bs).drop needle.length)
      else countMatches needle positive bs := by
  rw [countMatches.eq_def]
  by_cases h : (b :: bs).take needle.length = needle
  · simp [countFirstMatch, h]
  · simp only [countFirstMatch, beq_iff_eq, h, ↓reduceIte]
    cases bs with
    | nil => simp [countFirstMatch, countMatches]
    | cons c cs =>
      rw [countMatches.eq_def]
      cases hs : countFirstMatch needle (c :: cs) <;>
        simp [hs, Nat.succ_add, List.drop_succ_cons]

private theorem countMatches_greedy (needle : List Nat) (positive : 0 < needle.length)
    (text : List Nat) : countMatches needle positive text = greedyCount needle positive text := by
  cases text with
  | nil => rw [countMatches.eq_def, greedyCount.eq_def]
  | cons b bs =>
    rw [countMatches_step, greedyCount.eq_def]
    dsimp only
    split
    · rw [countMatches_greedy]
    · exact countMatches_greedy needle positive bs
termination_by text.length
decreasing_by all_goals simp [List.length_drop] <;> omega

theorem count_empty_spec (text : List Nat) : count text [] = countRunes text + 1 := rfl

theorem count_byte_spec (text : List Nat) (b : Nat) :
    count text [b] = text.countP (· == b) := by
  simp [count, List.count_eq_countP]

theorem count_nonoverlap_spec (text : List Nat) (b c : Nat) (bs : List Nat) :
    count text (b :: c :: bs) = greedyCount (b :: c :: bs) (by simp) text :=
  countMatches_greedy _ _ _

end GoLean.Strings
