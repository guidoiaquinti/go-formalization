import GoLean.Strings.IndexRune
import GoLean.Unicode.Utf8.DecodeRuneInString

namespace GoLean.Strings

/-- Model the membership relation of Go's 256-slot ASCII table. -/
def allASCII (bytes : List Nat) : Bool := bytes.all (· < 128)

private def firstAnyRune (chars : List Nat) (offset : Nat) (input : List Nat) : Nat → Option Nat
  | 0 => none
  | fuel + 1 =>
    match input with
    | [] => none
    | lead :: rest =>
      let (rune, width) := GoLean.Unicode.Utf8.decodeNonempty lead rest
      if (indexRune chars rune).isSome then some offset
      else firstAnyRune chars (offset + width) (input.drop width) fuel

/-- Go 1.27.0 strings/strings.go:184–214. An offset is a byte index.
    The ASCII membership table is represented by its membership predicate. -/
def indexAny (input chars : List Nat) : Option Nat :=
  match chars with
  | [] => none
  | [byte] => indexRune input (if byte < 128 then byte else 0xfffd)
  | _ :: _ :: _ =>
    if input.length > 8 && allASCII chars then input.findIdx? (chars.contains ·)
    else firstAnyRune chars 0 input input.length

theorem indexAny_emptyChars (input : List Nat) : indexAny input [] = none := rfl

theorem indexAny_singleHighByte (input : List Nat) (byte : Nat) (high : 128 ≤ byte) :
    indexAny input [byte] = indexRune input 0xfffd := by
  simp [indexAny, Nat.not_lt.mpr high]

theorem indexAny_asciiFirst (lead c d : Nat) (rest chars : List Nat)
    (long : 8 < (lead :: rest).length) (ascii : allASCII (c :: d :: chars) = true)
    (member : (c :: d :: chars).contains lead = true) :
    indexAny (lead :: rest) (c :: d :: chars) = some 0 := by
  simp only [indexAny, long, decide_true, ascii, Bool.true_and, ite_true, List.findIdx?_cons, member]

end GoLean.Strings
