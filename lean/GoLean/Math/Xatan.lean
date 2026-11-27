import GoLean.Math.Atan

namespace GoLean.Math

/-- Go 1.27.0 src/math/atan.go:55–72. Model the rational core on the signed
    interval |x| ≤ 0.66, including negative arguments supplied by satan.
    Outside this interval the raw word evaluator does not model IEEE overflow. -/
def xatanWord (fused : Bool) (word : BitVec 64) : Nat :=
  atanSmallWord fused word.toNat

set_option exponentiation.threshold 4096 in
set_option maxRecDepth 4096 in
theorem xatanWord_zeros (fused : Bool) :
    xatanWord fused 0 = 0 ∧ xatanWord fused (BitVec.ofNat 64 0x8000000000000000) = 0 := by
  cases fused <;> decide

set_option exponentiation.threshold 4096 in
set_option maxRecDepth 4096 in
theorem xatanWord_minSubnormal (fused : Bool) :
    xatanWord fused 1 = 1 ∧
      xatanWord fused (BitVec.ofNat 64 0x8000000000000001) = 0x8000000000000001 := by
  cases fused <;> decide

end GoLean.Math
