import GoLean.Strings.Count

namespace GoLean.Unicode.Utf8

/-- The width test matches the UTF-8 acceptance ranges. The arithmetic
    reconstructs the payload bits of a valid sequence, or yields RuneError
    for one malformed byte.
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf8/utf8.go#L219-L266 -/
def decodeNonempty (lead : Nat) (rest : List Nat) : Nat × Nat :=
  let width := Strings.runeWidth lead rest
  let rune := match width with
    | 2 => lead % 32 * 64 + rest[0]?.getD 0 % 64
    | 3 => lead % 16 * 4096 + rest[0]?.getD 0 % 64 * 64 + rest[1]?.getD 0 % 64
    | 4 => lead % 8 * 262144 + rest[0]?.getD 0 % 64 * 4096 +
        rest[1]?.getD 0 % 64 * 64 + rest[2]?.getD 0 % 64
    | _ => if lead < 128 then lead else 0xfffd
  (rune, width)

/-- Go strings are raw byte strings, including incomplete or malformed sequences. -/
def decodeRuneInString : List Nat → Nat × Nat
  | [] => (0xfffd, 0)
  | lead :: rest => decodeNonempty lead rest

theorem decodeRuneInString_empty : decodeRuneInString [] = (0xfffd, 0) := rfl

theorem decodeRuneInString_ascii (lead : Nat) (rest : List Nat) (ascii : lead < 128) :
    decodeRuneInString (lead :: rest) = (lead, 1) := by
  have h2 : ¬ 0xc2 ≤ lead := by omega
  have h3 : ¬ 0xe0 ≤ lead := by omega
  have h4 : ¬ 0xf0 ≤ lead := by omega
  match rest with
  | [] | [_] | [_, _] | _ :: _ :: _ :: _ =>
    simp [decodeRuneInString, decodeNonempty, Strings.runeWidth, h2, h3, h4, ascii]

theorem decodeRuneInString_error (lead : Nat) (rest : List Nat) (nonASCII : 128 ≤ lead)
    (rejected : Strings.runeWidth lead rest = 1) :
    decodeRuneInString (lead :: rest) = (0xfffd, 1) := by
  simp [decodeRuneInString, decodeNonempty, rejected, Nat.not_lt_of_ge nonASCII]

theorem decodeRuneInString_two (lead : Nat) (rest : List Nat)
    (width : Strings.runeWidth lead rest = 2) :
    decodeRuneInString (lead :: rest) =
      (lead % 32 * 64 + rest[0]?.getD 0 % 64, 2) := by
  simp [decodeRuneInString, decodeNonempty, width]

theorem decodeRuneInString_three (lead : Nat) (rest : List Nat)
    (width : Strings.runeWidth lead rest = 3) :
    decodeRuneInString (lead :: rest) =
      (lead % 16 * 4096 + rest[0]?.getD 0 % 64 * 64 + rest[1]?.getD 0 % 64, 3) := by
  simp [decodeRuneInString, decodeNonempty, width]

theorem decodeRuneInString_four (lead : Nat) (rest : List Nat)
    (width : Strings.runeWidth lead rest = 4) :
    decodeRuneInString (lead :: rest) =
      (lead % 8 * 262144 + rest[0]?.getD 0 % 64 * 4096 +
        rest[1]?.getD 0 % 64 * 64 + rest[2]?.getD 0 % 64, 4) := by
  simp [decodeRuneInString, decodeNonempty, width]

end GoLean.Unicode.Utf8
