import GoLean.Strings.SpecialCase
import GoLean.Strings.ToUpper

namespace GoLean.Strings

/-- Go 1.27.0 src/strings/strings.go:751–755. Apply the special uppercase
    mapping through Map, falling back to the pinned uppercase table on miss. -/
def toUpperSpecial (ranges : Array SpecialCaseRange) (input : List (BitVec 8)) : List Nat :=
  (mapRunes (specialCaseRune .upper ranges upperRune) input).2

theorem toUpperSpecial_fallback (input : List (BitVec 8)) :
    toUpperSpecial #[] input = (mapRunes (fun rune => Int.ofNat (upperRune rune)) input).2 := rfl

theorem toUpperSpecial_identityOverride :
    toUpperSpecial #[⟨105, 105, 0, 0, 0⟩] [105] = [105] := rfl

theorem toUpperSpecial_alternating :
    toUpperSpecial #[⟨65, 68, 0x110000, 0, 0⟩] [65, 66, 67, 68] = [65, 65, 67, 67] := rfl

theorem toUpperSpecial_negativeDrop :
    toUpperSpecial #[⟨65, 65, BitVec.ofInt 32 (-66), 0, 0⟩] [65] = [] := rfl

theorem toUpperSpecial_midpoint :
    toUpperSpecial #[⟨97, 122, BitVec.ofInt 32 (-100), 0, 0⟩,
      ⟨97, 122, 0, 0, 0⟩, ⟨97, 122, 10, 0, 0⟩] [97] = [97] := rfl

end GoLean.Strings
