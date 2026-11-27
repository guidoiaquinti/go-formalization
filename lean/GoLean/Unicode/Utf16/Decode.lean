import GoLean.Unicode.Utf16.DecodeRune

namespace GoLean.Unicode.Utf16

/-- Go's loop advances once for a BMP word or invalid surrogate and twice
    for a valid pair. Represent the input uint16 words by natural numbers.
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf16/utf16.go#L114-L144 -/
def decode : List Nat → List Nat
  | [] => []
  | first :: rest =>
    if first < 0xd800 ∨ 0xe000 ≤ first then first :: decode rest
    else match rest with
      | [] => [0xfffd]
      | second :: tail =>
        if 0xd800 ≤ first ∧ first < 0xdc00 ∧ 0xdc00 ≤ second ∧ second < 0xe000 then
          decodeRune first second :: decode tail
        else
          0xfffd :: decode (second :: tail)

termination_by words => words.length
decreasing_by all_goals simp_wf <;> omega

/-- The recurrence determines the full sequence and the words consumed in each step. -/
theorem decode_spec (words : List Nat) :
    decode words =
      match words with
      | [] => []
      | first :: rest =>
        if first / 1024 ≠ 54 ∧ first / 1024 ≠ 55 then first :: decode rest
        else match rest with
          | [] => [65533]
          | second :: tail =>
            if first / 1024 = 54 ∧ second / 1024 = 55 then
              (65536 + 1024 * (first % 1024) + second % 1024) :: decode tail
            else
              65533 :: decode (second :: tail) := by
  cases words with
  | nil => rw [decode.eq_def]
  | cons first rest =>
    rw [decode.eq_def]
    dsimp only
    have bmp : (first < 0xd800 ∨ 0xe000 ≤ first) ↔
        (first / 1024 ≠ 54 ∧ first / 1024 ≠ 55) := by omega
    by_cases h : first < 0xd800 ∨ 0xe000 ≤ first
    · rw [ite_eq_left h, ite_eq_left (bmp.mp h)]
    · rw [ite_eq_right h, ite_eq_right (mt bmp.mpr h)]
      cases rest with
      | nil => rfl
      | cons second tail =>
        have pair : (0xd800 ≤ first ∧ first < 0xdc00 ∧ 0xdc00 ≤ second ∧ second < 0xe000) ↔
            (first / 1024 = 54 ∧ second / 1024 = 55) := by omega
        dsimp only
        by_cases hp : 0xd800 ≤ first ∧ first < 0xdc00 ∧ 0xdc00 ≤ second ∧ second < 0xe000
        · have prefixes := pair.mp hp
          have signed : (first : Int) / 1024 = 54 ∧ (second : Int) / 1024 = 55 := by omega
          rw [ite_eq_left hp, ite_eq_left prefixes, decodeRune_spec, ite_eq_left signed]
          congr 1
          omega
        · rw [ite_eq_right hp, ite_eq_right (mt pair.mpr hp)]

end GoLean.Unicode.Utf16
