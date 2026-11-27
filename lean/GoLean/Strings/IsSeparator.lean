namespace GoLean.Strings

private def asciiWordRune (r : Int) : Bool :=
  (48 ≤ r && r ≤ 57) || (97 ≤ r && r ≤ 122) || (65 ≤ r && r ≤ 90) || r == 95

/-- The part of Unicode 17.0.0 White_Space above ASCII. -/
private def nonASCIISpace (r : Int) : Bool :=
  r == 0x85 || r == 0xa0 || r == 0x1680 ||
  (0x2000 ≤ r && r ≤ 0x200a) || r == 0x2028 || r == 0x2029 ||
  r == 0x202f || r == 0x205f || r == 0x3000

/-- Go 1.27.0 src/strings/strings.go:821–842. Above ASCII the letter and
    digit tests cannot reject any White_Space member, so classification
    reduces to that table. Signed values below zero take the ASCII branch. -/
def isSeparator (r : Int) : Bool :=
  if r ≤ 127 then !asciiWordRune r else nonASCIISpace r

theorem isSeparator_ascii (r : Int) (ascii : r ≤ 127) :
    isSeparator r = true ↔ ¬ ((48 ≤ r ∧ r ≤ 57) ∨ (97 ≤ r ∧ r ≤ 122) ∨
      (65 ≤ r ∧ r ≤ 90) ∨ r = 95) := by
  simp [isSeparator, ascii, asciiWordRune]
  omega

theorem isSeparator_negative (r : Int) (negative : r < 0) : isSeparator r = true := by
  rw [isSeparator_ascii r (by omega)]
  omega

theorem isSeparator_aboveSpace (r : Int) (above : 0x3000 < r) :
    isSeparator r = false := by
  have nonASCII : ¬ r ≤ 127 := by omega
  simp [isSeparator, nonASCII, nonASCIISpace]
  omega

end GoLean.Strings
