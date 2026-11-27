import GoLean.Strings.IsSeparator

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:790–811 has the same signed-rune
    classification as the strings helper. -/
def isSeparator (r : Int) : Bool :=
  GoLean.Strings.isSeparator r

theorem isSeparator_ascii (r : Int) (ascii : r ≤ 127) :
    isSeparator r = true ↔ ¬ ((48 ≤ r ∧ r ≤ 57) ∨ (97 ≤ r ∧ r ≤ 122) ∨
      (65 ≤ r ∧ r ≤ 90) ∨ r = 95) :=
  GoLean.Strings.isSeparator_ascii r ascii

theorem isSeparator_negative (r : Int) (negative : r < 0) : isSeparator r = true :=
  GoLean.Strings.isSeparator_negative r negative

theorem isSeparator_aboveSpace (r : Int) (above : 0x3000 < r) : isSeparator r = false :=
  GoLean.Strings.isSeparator_aboveSpace r above

end GoLean.Bytes
