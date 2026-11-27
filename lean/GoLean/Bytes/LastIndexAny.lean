import GoLean.Bytes.ContainsRunePrivate
import GoLean.Bytes.MakeASCIISet
import GoLean.Bytes.ShouldUseASCIISet
import GoLean.Unicode.Utf8.DecodeLastRune

namespace GoLean.Bytes

private def lastAnyRune (accepts : Nat → Bool) (input : List Nat) : Nat → Option Nat
  | 0 => none
  | fuel + 1 =>
    if input.isEmpty then none
    else
      let (rune, width) := GoLean.Unicode.Utf8.decodeLastRune input
      let offset := input.length - width
      if accepts rune then some offset
      else lastAnyRune accepts (input.take offset) fuel

/-- Go 1.27.0 src/bytes/bytes.go:291–369. Try the ASCII table before
    singleton branches, then decode backward using original byte widths.
    The general predicate abstracts equivalent encoded membership probes. -/
def lastIndexAny (input chars : List (BitVec 8)) : Option Nat :=
  if chars.isEmpty then none
  else
    let (table, ascii) := makeASCIISet chars
    if shouldUseASCIISet input.length && ascii then
      (input.reverse.findIdx? table).map (fun index => input.length - 1 - index)
    else
      let members := chars.map BitVec.toNat
      match input with
      | [byte] =>
        let accepts := if byte.toNat < 128 then chars.contains byte
          else containsRunePrivate members 0xfffd
        if accepts then some 0 else none
      | _ =>
        let accepts := match chars with
          | [byte] => fun rune => rune == (if byte.toNat < 128 then byte.toNat else 0xfffd)
          | _ => fun rune => containsRunePrivate members rune
        lastAnyRune accepts (input.map BitVec.toNat) input.length

theorem lastIndexAny_emptyChars (input : List (BitVec 8)) : lastIndexAny input [] = none := rfl

theorem lastIndexAny_emptyInput (chars : List (BitVec 8)) : lastIndexAny [] chars = none := by
  simp [lastIndexAny, shouldUseASCIISet, GoLean.Strings.shouldUseASCIISet, lastAnyRune]

theorem lastIndexAny_singleHigh (byte c : BitVec 8) (chars : List (BitVec 8))
    (high : 128 ≤ byte.toNat) :
    lastIndexAny [byte] (c :: chars) =
      if containsRunePrivate ((c :: chars).map BitVec.toNat) 0xfffd then some 0 else none := by
  simp [lastIndexAny, shouldUseASCIISet, GoLean.Strings.shouldUseASCIISet, Nat.not_lt.mpr high]

end GoLean.Bytes
