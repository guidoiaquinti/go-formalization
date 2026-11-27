import GoLean.Strings.SpecialCase
import GoLean.Strings.ToTitle

namespace GoLean.Strings

/-- Go 1.27.0 src/strings/strings.go:763–767. Apply the special title-case
    mapping through Map, falling back to the pinned title-case table on miss. -/
def toTitleSpecial (ranges : Array SpecialCaseRange) (input : List (BitVec 8)) : List Nat :=
  (mapRunes (specialCaseRune .title ranges titleRune) input).2

theorem toTitleSpecial_fallback (input : List (BitVec 8)) :
    toTitleSpecial #[] input = toTitle input := rfl

theorem toTitleSpecial_identityOverride :
    toTitleSpecial #[⟨105, 105, 0, 0, 0⟩] [105] = [105] := rfl

theorem toTitleSpecial_alternating :
    toTitleSpecial #[⟨65, 68, 0, 0, 0x110000⟩] [65, 66, 67, 68] = [65, 65, 67, 67] := rfl

theorem toTitleSpecial_negativeDrop :
    toTitleSpecial #[⟨65, 65, 0, 0, BitVec.ofInt 32 (-66)⟩] [65] = [] := rfl

theorem toTitleSpecial_midpoint :
    toTitleSpecial #[⟨97, 122, 0, 0, BitVec.ofInt 32 (-100)⟩,
      ⟨97, 122, 0, 0, 0⟩, ⟨97, 122, 0, 0, 10⟩] [97] = [97] := rfl

theorem toTitleSpecial_titleDelta :
    toTitleSpecial #[⟨65, 65, 1, 2, 3⟩] [65] = [68] := rfl

end GoLean.Strings
