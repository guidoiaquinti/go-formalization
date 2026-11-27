import GoLean.Math.Erfinv

namespace GoLean.Math

/-- Go 1.27.0 src/math/erfinv.go:127–129. The subtraction is rounded
    to binary64 before dispatch to Erfinv. Nonfinite input always
    gives NaN by class. -/
def erfcinvWord (arm64 fused : Bool) (word : BitVec 64) : Option Nat :=
  if 0x7ff0000000000000 ≤ word.toNat % 2 ^ 63 then none
  else erfinvWord arm64 fused (BitVec.ofNat 64 (subtractWords 0x3ff0000000000000 word.toNat))

theorem erfcinvWord_nonfinite (arm64 fused : Bool) (word : BitVec 64)
    (nonfinite : 0x7ff0000000000000 ≤ word.toNat % 2 ^ 63) :
    erfcinvWord arm64 fused word = none := by
  simp [erfcinvWord, nonfinite]

set_option exponentiation.threshold 4096 in
set_option maxRecDepth 4096 in
theorem erfcinvWord_endpoints (arm64 fused : Bool) :
    erfcinvWord arm64 fused 0 = some 0x7ff0000000000000 ∧
      erfcinvWord arm64 fused (BitVec.ofNat 64 0x8000000000000000) = some 0x7ff0000000000000 ∧
      erfcinvWord arm64 fused (BitVec.ofNat 64 0x3ff0000000000000) = some 0 ∧
      erfcinvWord arm64 fused (BitVec.ofNat 64 0x4000000000000000) = some 0xfff0000000000000 := by
  cases arm64 <;> cases fused <;> decide

set_option exponentiation.threshold 4096 in
set_option maxRecDepth 4096 in
theorem erfcinvWord_tinyNegative (arm64 fused : Bool) :
    erfcinvWord arm64 fused (BitVec.ofNat 64 0xbca0000000000000) = some 0x7ff0000000000000 ∧
      erfcinvWord arm64 fused (BitVec.ofNat 64 0xbca0000000000001) = none ∧
      erfcinvWord arm64 fused (BitVec.ofNat 64 0x8000000000000001) = some 0x7ff0000000000000 := by
  cases arm64 <;> cases fused <;> decide

end GoLean.Math
