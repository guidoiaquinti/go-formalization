import GoLean.Strings.Count

namespace GoLean.Unicode.Utf8

/-- The string range advances by one decoded width at each iteration.
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf8/utf8.go#L446-L452 -/
def runeCountInString (input : List Nat) : Nat := Strings.countRunes input

/-- Go counts leading ASCII bytes directly, then delegates the remaining bytes.
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf8/utf8.go#L432-L444 -/
def runeCount : List Nat → Nat
  | [] => 0
  | b :: rest => if b < 128 then 1 + runeCount rest else runeCountInString (b :: rest)

private theorem asciiWidth (b : Nat) (rest : List Nat) (h : b < 128) :
    Strings.runeWidth b rest = 1 := by
  have h2 : ¬ 0xc2 ≤ b := by omega
  have h3 : ¬ 0xe0 ≤ b := by omega
  have h4 : ¬ 0xf0 ≤ b := by omega
  match rest with
  | [] | [_] | [_, _] | _ :: _ :: _ :: _ => simp [Strings.runeWidth, h2, h3, h4]

theorem runeCountInString_step (lead : Nat) (rest : List Nat) :
    runeCountInString (lead :: rest) =
        1 + runeCountInString (rest.drop (Strings.runeWidth lead rest - 1)) := by
  unfold runeCountInString
  rw [Strings.countRunes.eq_def]

theorem runeCountInString_bound (input : List Nat) :
    runeCountInString input ≤ input.length := by
  cases input with
  | nil => simp [runeCountInString, Strings.countRunes]
  | cons b rest =>
    rw [runeCountInString_step]
    have ih := runeCountInString_bound (rest.drop (Strings.runeWidth b rest - 1))
    have len : (rest.drop (Strings.runeWidth b rest - 1)).length ≤ rest.length := by simp
    simp only [List.length_cons]
    omega
termination_by input.length
decreasing_by simp [List.length_drop]; omega

theorem runeCount_spec (input : List Nat) : runeCount input = runeCountInString input := by
  induction input with
  | nil => simp [runeCount, runeCountInString, Strings.countRunes]
  | cons b rest ih =>
    by_cases h : b < 128
    · rw [runeCount, ite_eq_left h, ih, runeCountInString_step, asciiWidth b rest h]
      rfl
    · rw [runeCount, ite_eq_right h]

end GoLean.Unicode.Utf8
