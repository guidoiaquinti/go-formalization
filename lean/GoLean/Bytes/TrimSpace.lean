import GoLean.Bytes.TrimFunc

namespace GoLean.Bytes

private def trimASCIIWhitespace (byte : BitVec 8) : Bool :=
  let n := byte.toNat
  (9 ≤ n && n ≤ 13) || n == 32

/-- Go 1.27.0 unicode/graphic.go:128–137 and Unicode 17.0.0 White_Space. -/
private def trimUnicodeWhitespace (c : Nat) : Bool :=
  (9 ≤ c && c ≤ 13) || c == 32 || c == 0x85 || c == 0xa0 || c == 0x1680 ||
  (0x2000 ≤ c && c ≤ 0x200a) || c == 0x2028 || c == 0x2029 ||
  c == 0x202f || c == 0x205f || c == 0x3000

private def trimSpaceFallback (input : List (BitVec 8)) (spare : Nat) :
    Option (List (BitVec 8) × Nat) :=
  (trimFunc input spare (some trimUnicodeWhitespace)).bind (·.1)

/-- Go 1.27.0 src/bytes/bytes.go:1111–1140. Either non-ASCII edge
    delegates to TrimFunc. The right scan shortens length alone, so the
    discarded trailing ASCII whitespace still contributes to capacity. -/
def trimSpace (input : List (BitVec 8)) (spare : Nat) : Option (List (BitVec 8) × Nat) :=
  let remaining := input.dropWhile trimASCIIWhitespace
  match remaining with
  | [] => none
  | head :: _ =>
    if head.toNat < 128 then
      let reversed := (remaining.reverse.dropWhile trimASCIIWhitespace)
      match reversed with
      | [] => none
      | last :: _ =>
        if last.toNat < 128 then some (reversed.reverse, remaining.length + spare)
        else trimSpaceFallback reversed.reverse (remaining.length - reversed.length + spare)
    else trimSpaceFallback remaining spare

theorem trimSpace_empty (spare : Nat) : trimSpace [] spare = none := rfl

theorem trimSpace_skipLeft (head : BitVec 8) (tail : List (BitVec 8)) (spare : Nat)
    (space : trimASCIIWhitespace head = true) :
    trimSpace (head :: tail) spare = trimSpace tail spare := by
  simp [trimSpace, List.dropWhile, space]

theorem trimSpace_keepASCIIEnds (head last : BitVec 8) (middle : List (BitVec 8))
    (spare : Nat) (headASCII : head.toNat < 128) (lastASCII : last.toNat < 128)
    (headNonspace : trimASCIIWhitespace head = false)
    (lastNonspace : trimASCIIWhitespace last = false) :
    trimSpace (head :: (middle ++ [last])) spare =
      some (head :: (middle ++ [last]), (head :: (middle ++ [last])).length + spare) := by
  simp [trimSpace, List.dropWhile, headASCII, lastASCII, headNonspace, lastNonspace]

end GoLean.Bytes
