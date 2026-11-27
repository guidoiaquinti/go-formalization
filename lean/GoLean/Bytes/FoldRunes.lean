import GoLean.Strings.EqualFold
import GoLean.Unicode.Utf8.DecodeRune
import GoLean.Unicode.Utf8.RuneCount

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:1213–1285. The ASCII scan and ordered
    SimpleFold orbit match the string implementation. The slice decoder
    yields the same rune and width as the string decoder. -/
def equalFoldBytes (first second : List Nat) : Bool :=
  GoLean.Strings.equalFold first second

theorem equalFoldBytes_ascii (a b : Nat) (as bs : List Nat) (aASCII : a < 128) (bASCII : b < 128) :
    equalFoldBytes (a :: as) (b :: bs) = true ↔
      (a = b ∨ (65 ≤ a ∧ a ≤ 90 ∧ b = a + 32) ∨ (65 ≤ b ∧ b ≤ 90 ∧ a = b + 32)) ∧
        equalFoldBytes as bs = true :=
  GoLean.Strings.equalFold_asciiHead a b as bs aASCII bASCII

theorem equalFoldBytes_replacement : equalFoldBytes [0xff] [0xef, 0xbf, 0xbd] = true :=
  GoLean.Strings.equalFold_invalidRuneError

theorem equalFoldBytes_byteCount : equalFoldBytes [0xff, 0x80] [0xef, 0xbf, 0xbd] = false :=
  GoLean.Strings.equalFold_invalidByteCount

structure RunesResult where
  values : List Nat
  capacity : Nat

private def decodeRunes (input : List Nat) : Nat → List Nat
  | 0 => []
  | remaining + 1 =>
    match input with
    | [] => []
    | _ :: _ =>
      let (rune, width) := GoLean.Unicode.Utf8.decodeRune input
      rune :: decodeRunes (input.drop width) remaining

/-- Go 1.27.0 src/bytes/bytes.go:1144–1154. Count first for the exact
    allocation, then decode and store each rune from the unread suffix.
    The resulting slice is always nonnil, including for empty input. -/
def runesBytes (input : List Nat) : RunesResult :=
  let capacity := GoLean.Unicode.Utf8.runeCount input
  ⟨decodeRunes input input.length, capacity⟩

theorem runesBytes_empty : runesBytes [] = ⟨[], 0⟩ := rfl

theorem runesBytes_asciiHead (lead : Nat) (rest : List Nat) (ascii : lead < 128) :
    (runesBytes (lead :: rest)).values = lead :: (runesBytes rest).values := by
  simp [runesBytes, decodeRunes, GoLean.Unicode.Utf8.decodeRune, ascii]

theorem runesBytes_malformed_and_literal :
    (runesBytes [0xe2, 0x82, 0xef, 0xbf, 0xbd]).values = [0xfffd, 0xfffd, 0xfffd] := rfl

end GoLean.Bytes
