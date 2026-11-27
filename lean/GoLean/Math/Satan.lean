import GoLean.Math.Atan

namespace GoLean.Math

/-- Go 1.27.0 src/math/atan.go:76–88. The reduction receives nonnegative
    magnitudes, including positive infinity. Both zero encodings are
    observed as positive zero in direct calls to the internal helper. -/
def satanWord (fused : Bool) (word : BitVec 64) : Nat :=
  atanPositiveWord fused (word.toNat % 2 ^ 63)

set_option exponentiation.threshold 4096 in
set_option maxRecDepth 4096 in
theorem satanWord_zeros (fused : Bool) :
    satanWord fused 0 = 0 ∧ satanWord fused (BitVec.ofNat 64 0x8000000000000000) = 0 := by
  cases fused <;> decide

set_option exponentiation.threshold 4096 in
set_option maxRecDepth 4096 in
theorem satanWord_infinity (fused : Bool) :
    satanWord fused (BitVec.ofNat 64 0x7ff0000000000000) = 0x3ff921fb54442d18 := by
  cases fused <;> decide

end GoLean.Math
