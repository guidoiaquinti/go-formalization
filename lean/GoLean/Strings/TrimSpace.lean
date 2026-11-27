import GoLean.Strings.TrimFunc

namespace GoLean.Strings

private def trimASCIIWhitespace (byte : BitVec 8) : Bool :=
  let n := byte.toNat
  (9 ≤ n && n ≤ 13) || n == 32

/-- Go 1.27.0 unicode/graphic.go:128–137 and the Unicode 17.0.0
    White_Space table. -/
private def trimUnicodeWhitespace (c : Nat) : Bool :=
  (9 ≤ c && c ≤ 13) || c == 32 || c == 0x85 || c == 0xa0 || c == 0x1680 ||
  (0x2000 ≤ c && c ≤ 0x200a) || c == 0x2028 || c == 0x2029 ||
  c == 0x202f || c == 0x205f || c == 0x3000

private def trimSpaceRight (input : List (BitVec 8)) : List (BitVec 8) :=
  let remaining := input.reverse.dropWhile trimASCIIWhitespace
  match remaining with
  | [] => []
  | last :: _ =>
    if last.toNat < 128 then remaining.reverse
    else ((trimRightFunc remaining.reverse (some trimUnicodeWhitespace)).getD ([], [])).1

/-- Go 1.27.0 src/strings/strings.go:1089–1116. Find the first
    nonspace ASCII byte before scanning the right side. A non-ASCII
    byte switches to the two-sided or right-only Unicode scan depending
    on which side encountered it. -/
def trimSpace (input : List (BitVec 8)) : List (BitVec 8) :=
  let remaining := input.dropWhile trimASCIIWhitespace
  match remaining with
  | [] => []
  | head :: _ =>
    if head.toNat < 128 then trimSpaceRight remaining
    else ((trimFunc remaining (some trimUnicodeWhitespace)).getD ([], [])).1

theorem trimSpace_empty : trimSpace [] = [] := rfl

theorem trimSpace_skipLeft (head : BitVec 8) (tail : List (BitVec 8))
    (space : trimASCIIWhitespace head = true) :
    trimSpace (head :: tail) = trimSpace tail := by
  simp [trimSpace, List.dropWhile, space]

theorem trimSpace_keepASCIIEnds (head last : BitVec 8) (middle : List (BitVec 8))
    (headASCII : head.toNat < 128) (lastASCII : last.toNat < 128)
    (headNonspace : trimASCIIWhitespace head = false)
    (lastNonspace : trimASCIIWhitespace last = false) :
    trimSpace (head :: (middle ++ [last])) = head :: (middle ++ [last]) := by
  simp [trimSpace, trimSpaceRight, List.dropWhile, headASCII, lastASCII, headNonspace, lastNonspace]

end GoLean.Strings
