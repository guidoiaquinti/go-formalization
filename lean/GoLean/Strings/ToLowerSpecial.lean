import GoLean.Strings.SpecialCase
import GoLean.Strings.ToLower

namespace GoLean.Strings

/-- Go 1.27.0 src/strings/strings.go:757–761. Apply the special lowercase
    mapping through Map, falling back to the pinned lowercase table on miss. -/
def toLowerSpecial (ranges : Array SpecialCaseRange) (input : List (BitVec 8)) : List Nat :=
  (mapRunes (specialCaseRune .lower ranges lowerRune) input).2

theorem toLowerSpecial_fallback (input : List (BitVec 8)) :
    toLowerSpecial #[] input = (mapRunes (fun rune => Int.ofNat (lowerRune rune)) input).2 := rfl

theorem toLowerSpecial_identityOverride :
    toLowerSpecial #[⟨73, 73, 0, 0, 0⟩] [73] = [73] := rfl

theorem toLowerSpecial_alternating :
    toLowerSpecial #[⟨65, 68, 0, 0x110000, 0⟩] [65, 66, 67, 68] = [66, 66, 68, 68] := rfl

theorem toLowerSpecial_negativeDrop :
    toLowerSpecial #[⟨65, 65, 0, BitVec.ofInt 32 (-66), 0⟩] [65] = [] := rfl

theorem toLowerSpecial_midpoint :
    toLowerSpecial #[⟨65, 90, 0, BitVec.ofInt 32 (-100), 0⟩,
      ⟨65, 90, 0, 0, 0⟩, ⟨65, 90, 0, 10, 0⟩] [65] = [65] := rfl

end GoLean.Strings
