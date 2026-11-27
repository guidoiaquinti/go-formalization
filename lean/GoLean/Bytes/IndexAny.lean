import GoLean.Bytes.ContainsRunePrivate
import GoLean.Bytes.IndexRune
import GoLean.Bytes.MakeASCIISet
import GoLean.Bytes.ShouldUseASCIISet
import GoLean.Unicode.Utf8.DecodeRune

namespace GoLean.Bytes

private def indexAnyRunes (chars : List Nat) (offset : Nat) (input : List Nat) :
    Nat → Option Nat
  | 0 => none
  | fuel + 1 =>
    match input with
    | [] => none
    | lead :: _ =>
      let (rune, width) := GoLean.Unicode.Utf8.decodeRune input
      let member := if lead < 128 then chars.contains lead
        else containsRunePrivate chars rune
      if member then some offset
      else indexAnyRunes chars (offset + width) (input.drop width) fuel

/-- Go 1.27.0 src/bytes/bytes.go:213–285. Preserve singleton priority,
    the ASCII-table gate, and original byte offsets. Unicode membership
    abstracts the source's equivalent encoded search shortcuts. -/
def indexAny (input chars : List (BitVec 8)) : Option Nat :=
  if chars.isEmpty then none
  else
    match input with
    | [byte] =>
      let member := if byte.toNat < 128 then chars.contains byte
        else containsRunePrivate (chars.map BitVec.toNat) 0xfffd
      if member then some 0 else none
    | _ =>
      match chars with
      | [byte] =>
        indexRune (input.map BitVec.toNat) (if byte.toNat < 128 then byte.toNat else 0xfffd)
      | _ =>
        let (table, ascii) := makeASCIISet chars
        if shouldUseASCIISet input.length && ascii then input.findIdx? table
        else indexAnyRunes (chars.map BitVec.toNat) 0 (input.map BitVec.toNat) input.length

theorem indexAny_emptyChars (input : List (BitVec 8)) : indexAny input [] = none := rfl

theorem indexAny_emptyInput (chars : List (BitVec 8)) : indexAny [] chars = none := by
  cases chars with
  | nil => rfl
  | cons c rest =>
    cases rest with
    | nil =>
      by_cases ascii : c.toNat < 128
      · simp [indexAny, ascii, indexRune, GoLean.Strings.indexRune, GoLean.Strings.firstRuneError,
          GoLean.Strings.encodeRune, GoLean.Strings.index, GoLean.Strings.countFirstMatch]
      · simp [indexAny, ascii, indexRune, GoLean.Strings.indexRune, GoLean.Strings.firstRuneError]
    | cons d tail =>
      simp [indexAny, shouldUseASCIISet, GoLean.Strings.shouldUseASCIISet, indexAnyRunes]

theorem indexAny_singleHigh (byte c : BitVec 8) (chars : List (BitVec 8))
    (high : 128 ≤ byte.toNat) :
    indexAny [byte] (c :: chars) =
      if containsRunePrivate ((c :: chars).map BitVec.toNat) 0xfffd then some 0 else none := by
  simp [indexAny, Nat.not_lt.mpr high]

end GoLean.Bytes
