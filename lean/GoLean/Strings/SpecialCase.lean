import Init.Data.BitVec.Basic

namespace GoLean.Strings

inductive CaseVariant where
  | upper
  | lower
  | title

structure SpecialCaseRange where
  lo : BitVec 32
  hi : BitVec 32
  upperDelta : BitVec 32
  lowerDelta : BitVec 32
  titleDelta : BitVec 32

private def findSpecialRange (ranges : Array SpecialCaseRange) (rune : BitVec 32)
    (lo hi : Nat) : Nat → Option SpecialCaseRange
  | 0 => none
  | fuel + 1 => do
    if lo ≥ hi then none
    else
      let mid := (lo + hi) / 2
      let entry ← ranges[mid]?
      if entry.lo.toInt ≤ rune.toInt && rune.toInt ≤ entry.hi.toInt then some entry
      else if rune.toInt < entry.lo.toInt then findSpecialRange ranges rune lo mid fuel
      else findSpecialRange ranges rune (mid + 1) hi fuel

private def caseDelta (variant : CaseVariant) (entry : SpecialCaseRange) : BitVec 32 :=
  match variant with
  | .upper => entry.upperDelta
  | .lower => entry.lowerDelta
  | .title => entry.titleDelta

/-- Go 1.27.0 src/unicode/letter.go:209–260, 300–325.
    Midpoint search chooses a range before applying its signed rune
    arithmetic. The fallback is used only when the search misses. -/
def specialCaseRune (variant : CaseVariant) (ranges : Array SpecialCaseRange)
    (fallback : Nat → Nat) (rune : Nat) : Int :=
  let raw := BitVec.ofNat 32 rune
  match findSpecialRange ranges raw 0 ranges.size ranges.size with
  | none => Int.ofNat (fallback rune)
  | some entry =>
    let delta := caseDelta variant entry
    if delta.toInt > 0x10ffff then
      let parity : BitVec 32 := match variant with
        | .lower => 1
        | _ => 0
      (entry.lo + (((raw - entry.lo) &&& 0xfffffffe) ||| parity)).toInt
    else (raw + delta).toInt

theorem specialCaseRune_empty (variant : CaseVariant) (fallback : Nat → Nat) (rune : Nat) :
    specialCaseRune variant #[] fallback rune = Int.ofNat (fallback rune) :=
  rfl

end GoLean.Strings
