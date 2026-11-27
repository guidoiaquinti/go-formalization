import GoLean.Strings.Count

namespace GoLean.Strings

private def replaceMatches (old replacement : List Nat) (input : List Nat) : Nat → List Nat
  | 0 => input
  | n + 1 =>
    match countFirstMatch old input with
    | none => input
    | some offset =>
      input.take offset ++ replacement ++
        replaceMatches old replacement (input.drop (offset + old.length)) n

private def insertAtRuneEnds (replacement : List Nat) (input : List Nat) : Nat → List Nat
  | 0 => input
  | n + 1 =>
    match input with
    | [] => []
    | lead :: rest =>
      let width := runeWidth lead rest
      input.take width ++ replacement ++ insertAtRuneEnds replacement (input.drop width) n

/-- Go 1.27.0 src/strings/strings.go:1136–1171. Count and clamp before
    replacing. A zero-width pattern inserts at the initial offset and then
    at decoded rune boundaries, including one boundary per malformed byte.
    The builder and its no-allocation exits are represented by their bytes. -/
def replace (input old replacement : List Nat) (limit : Int) : List Nat :=
  if old = replacement || limit == 0 then input
  else
    let matching := count input old
    if matching = 0 then input
    else
      let n := if limit < 0 then matching else min matching limit.toNat
      if old.isEmpty then replacement ++ insertAtRuneEnds replacement input (n - 1)
      else replaceMatches old replacement input n

theorem replace_zero (input old replacement : List Nat) :
    replace input old replacement 0 = input := by
  simp [replace]

theorem replace_same (input old : List Nat) (limit : Int) :
    replace input old old limit = input := by
  simp [replace]

theorem replace_missing (input old replacement : List Nat) (limit : Int)
    (absent : count input old = 0) : replace input old replacement limit = input := by
  simp [replace, absent]

theorem replace_emptyOnce (input replacement : List Nat) :
    replace input [] replacement 1 = replacement ++ input := by
  cases replacement with
  | nil => simp [replace]
  | cons head tail => simp [replace, count, Nat.min_eq_right, insertAtRuneEnds]

theorem replace_emptyMalformed :
    replace [0xe2, 0x82, 0xef, 0xbf, 0xbd] [] [0xff] (-1) =
      [0xff, 0xe2, 0xff, 0x82, 0xff, 0xef, 0xbf, 0xbd, 0xff] := by
  have first : runeWidth 226 [130, 239, 191, 189] = 1 := rfl
  have second : runeWidth 130 [239, 191, 189] = 1 := rfl
  have third : runeWidth 239 [191, 189] = 3 := rfl
  simp [replace, count, countRunes, insertAtRuneEnds, first, second, third]

end GoLean.Strings
