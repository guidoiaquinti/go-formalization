import GoLean.Bytes.SimpleCase
import GoLean.Strings.SpecialCase

namespace GoLean.Bytes

open GoLean.Strings (CaseVariant SpecialCaseRange specialCaseRune upperRune lowerRune titleRune)

private def mapSpecial (variant : CaseVariant) (fallback : Nat → Nat)
    (ranges : Array SpecialCaseRange) (input : List (BitVec 8)) : CaseBytes :=
  let mapped := ((mapRunes (specialCaseRune variant ranges fallback) input).2).getD []
  ⟨mapped, if mapped.length ≤ input.length then some input.length else none⟩

/-- Go 1.27.0 src/bytes/bytes.go:743–745. Use SpecialCase.ToUpper with
    the byte Map, retaining signed rune results through its drop check. -/
def toUpperSpecialBytes (ranges : Array SpecialCaseRange) (input : List (BitVec 8)) : CaseBytes :=
  mapSpecial .upper upperRune ranges input

/-- Go 1.27.0 src/bytes/bytes.go:749–751. Use the lowercase delta and fallback. -/
def toLowerSpecialBytes (ranges : Array SpecialCaseRange) (input : List (BitVec 8)) : CaseBytes :=
  mapSpecial .lower lowerRune ranges input

/-- Go 1.27.0 src/bytes/bytes.go:755–757. Use the titlecase delta and fallback. -/
def toTitleSpecialBytes (ranges : Array SpecialCaseRange) (input : List (BitVec 8)) : CaseBytes :=
  mapSpecial .title titleRune ranges input

theorem toUpperSpecialBytes_fallback (input : List (BitVec 8)) :
    toUpperSpecialBytes #[] input = mapCaseBytes upperRune input := rfl

theorem toLowerSpecialBytes_fallback (input : List (BitVec 8)) :
    toLowerSpecialBytes #[] input = mapCaseBytes lowerRune input := rfl

theorem toTitleSpecialBytes_fallback (input : List (BitVec 8)) :
    toTitleSpecialBytes #[] input = toTitleBytes input := rfl

theorem toUpperSpecialBytes_drop :
    toUpperSpecialBytes #[⟨97, 97, BitVec.ofInt 32 (-98), 0, 0⟩] [97] = ⟨[], some 1⟩ := rfl

theorem toLowerSpecialBytes_identity :
    toLowerSpecialBytes #[⟨65, 65, 0, 0, 0⟩] [65] = ⟨[65], some 1⟩ := rfl

theorem toTitleSpecialBytes_delta :
    toTitleSpecialBytes #[⟨65, 65, 1, 2, 3⟩] [65] = ⟨[68], some 1⟩ := rfl

end GoLean.Bytes
