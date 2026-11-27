import GoLean.Strings.IndexAny
import GoLean.Unicode.Utf8.DecodeLastRuneInString

namespace GoLean.Strings

private def singleByteRune (byte : Nat) : Nat :=
  if byte < 128 then byte else 0xfffd

private def lastAnyRune (accepts : Nat → Bool) (input : List Nat) : Nat → Option Nat
  | 0 => none
  | fuel + 1 =>
    if input.isEmpty then none
    else
      let (rune, width) := GoLean.Unicode.Utf8.decodeLastRuneInString input
      let offset := input.length - width
      if accepts rune then some offset
      else lastAnyRune accepts (input.take offset) fuel

/-- Go 1.27.0 strings/strings.go:218–260. Return the last matching rune's
    starting byte index. The ASCII table is represented by its membership relation. -/
def lastIndexAny (input chars : List Nat) : Option Nat :=
  let member := fun rune : Nat => (indexRune chars rune).isSome
  if chars.isEmpty then none
  else
    match input with
    | [byte] => if member (singleByteRune byte) then some 0 else none
    | _ =>
      if input.length > 8 && allASCII chars then
        (input.reverse.findIdx? (chars.contains ·)).map (fun i => input.length - 1 - i)
      else
        match chars with
        | [byte] => lastAnyRune (· == singleByteRune byte) input input.length
        | _ => lastAnyRune member input input.length

theorem lastIndexAny_emptyChars (input : List Nat) : lastIndexAny input [] = none := rfl

theorem lastIndexAny_singleHighByte (byte c : Nat) (chars : List Nat) (high : 128 ≤ byte) :
    lastIndexAny [byte] (c :: chars) =
      if (indexRune (c :: chars) 0xfffd).isSome then some 0 else none := by
  simp [lastIndexAny, singleByteRune, Nat.not_lt.mpr high]

theorem lastIndexAny_asciiTail (input : List Nat) (last c : Nat) (chars : List Nat)
    (long : 8 < (input ++ [last]).length) (ascii : allASCII (c :: chars) = true)
    (member : (c :: chars).contains last = true) :
    lastIndexAny (input ++ [last]) (c :: chars) = some input.length := by
  cases input with
  | nil => simp at long
  | cons head rest =>
    change lastIndexAny (head :: (rest ++ [last])) (c :: chars) = _
    have size : 8 < (head :: (rest ++ [last])).length := long
    have tail : rest ++ [last] ≠ [] := by simp
    cases h : rest ++ [last] with
    | nil => contradiction
    | cons inner suffix =>
      rw [h] at size
      change lastIndexAny (head :: inner :: suffix) (c :: chars) = _
      simp only [lastIndexAny, List.isEmpty_cons, Bool.false_eq_true, ite_false,
        size, decide_true, ascii, Bool.true_and, ite_true]
      rw [← h]
      change (List.findIdx? _ ((head :: rest) ++ [last]).reverse).map _ = _
      simp only [List.reverse_append, List.reverse_singleton, List.singleton_append,
        List.findIdx?_cons, member, ite_true, Option.map_some, List.length_cons, List.length_append,
        List.length_nil, Nat.zero_add, Nat.add_sub_cancel, Nat.sub_zero]

end GoLean.Strings
