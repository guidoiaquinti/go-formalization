import GoLean.Math.Asin

namespace GoLean.Math

/-- Go 1.27.0 src/math/asin.go:58–67. The public Asin result has already
    been rounded before the subtraction from the rounded π/2 constant.
    None represents NaN by class. -/
def acosWord (arm64 fused : Bool) (word : BitVec 64) : Option Nat :=
  (asinWord arm64 fused word).map (subtractWords 0x3ff921fb54442d18)

theorem acosWord_outside (arm64 fused : Bool) (word : BitVec 64)
    (outside : 0x3ff0000000000000 < word.toNat % 2 ^ 63) :
    acosWord arm64 fused word = none := by
  simp [acosWord, asinWord_outside arm64 fused word outside]

set_option exponentiation.threshold 4096 in
set_option maxRecDepth 4096 in
theorem acosWord_zeros (arm64 fused : Bool) :
    acosWord arm64 fused 0 = some 0x3ff921fb54442d18 ∧
      acosWord arm64 fused (BitVec.ofNat 64 0x8000000000000000) = some 0x3ff921fb54442d18 := by
  cases arm64 <;> cases fused <;> decide

set_option exponentiation.threshold 4096 in
set_option maxRecDepth 4096 in
theorem acosWord_endpoints (arm64 fused : Bool) :
    acosWord arm64 fused (BitVec.ofNat 64 0x3ff0000000000000) = some 0 ∧
      acosWord arm64 fused (BitVec.ofNat 64 0xbff0000000000000) = some 0x400921fb54442d18 := by
  cases arm64 <;> cases fused <;> decide

end GoLean.Math
