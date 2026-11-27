import GoLean.Math.Atan

namespace GoLean.Math

/-- Compute the rounded radicand and use the same square root for both
    arctangent reductions. The pinned arm64 compiler contracts 1-x*x. -/
def asinPositiveWord (arm64 fused : Bool) (magnitude : Nat) : Option Nat := do
  let radicand := multiplyAddWords arm64 (negateWord magnitude) magnitude 0x3ff0000000000000
  let root ← sqrtWord (BitVec.ofNat 64 radicand)
  if 0x3fe6666666666666 < magnitude then
    let ratio := divideWords root magnitude
    return subtractWords 0x3ff921fb54442d18 (atanPositiveWord fused ratio)
  else
    return atanPositiveWord fused (divideWords magnitude root)

/-- Go 1.27.0 src/math/asin.go:27–51. Preserve the signed zero, domain
    rejection, rounded square root and both arctangent reductions.
    None represents NaN by class. -/
def asinWord (arm64 fused : Bool) (word : BitVec 64) : Option Nat :=
  let raw := word.toNat
  let magnitude := raw % 2 ^ 63
  if magnitude = 0 then some raw
  else if 0x3ff0000000000000 < magnitude then none
  else (asinPositiveWord arm64 fused magnitude).map fun result =>
    if 2 ^ 63 ≤ raw then negateWord result else result

theorem asinWord_zero (arm64 fused : Bool) (word : BitVec 64)
    (zero : word.toNat % 2 ^ 63 = 0) : asinWord arm64 fused word = some word.toNat := by
  simp [asinWord, zero]

theorem asinWord_outside (arm64 fused : Bool) (word : BitVec 64)
    (outside : 0x3ff0000000000000 < word.toNat % 2 ^ 63) :
    asinWord arm64 fused word = none := by
  simp [asinWord, outside, show word.toNat % 2 ^ 63 ≠ 0 by omega]

theorem asinWord_opposite (arm64 fused : Bool) (x y : BitVec 64)
    (positive : x.toNat < 2 ^ 63) (negative : 2 ^ 63 ≤ y.toNat)
    (same : x.toNat % 2 ^ 63 = y.toNat % 2 ^ 63)
    (inRange : x.toNat % 2 ^ 63 ≤ 0x3ff0000000000000)
    (nonzero : x.toNat % 2 ^ 63 ≠ 0) :
    asinWord arm64 fused y = (asinWord arm64 fused x).map negateWord := by
  cases result : asinPositiveWord arm64 fused (x.toNat % 2 ^ 63) <;>
    simp [asinWord, ← same, nonzero, Nat.not_lt.mpr inRange, Nat.not_le.mpr positive, negative, result]

end GoLean.Math
