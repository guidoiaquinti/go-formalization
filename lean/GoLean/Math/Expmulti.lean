import GoLean.Math.Exp2

namespace GoLean.Math

/-- Go 1.27.0 src/math/exp.go:184–202. The standalone private helper uses
    public Ldexp after the rational reduction. Arm64 contracts the Horner
    stages and the final polynomial subtraction. Amd64 v3/v4 contract
    the Horner stages only. -/
def expmultiWord (arm64 compileFused : Bool) (hi lo : BitVec 64) (k : Int) : Nat :=
  let reduced := exp2Reduced (arm64 || compileFused) arm64 hi.toNat lo.toNat
  ldexp 64 reduced k

set_option exponentiation.threshold 4096 in
set_option maxRecDepth 4096 in
theorem expmultiWord_zero (arm64 fused : Bool) (k : Int) :
    expmultiWord arm64 fused 0 0 k = ldexp 64 0x3ff0000000000000 k := by
  have reduced : exp2Reduced (arm64 || fused) arm64 0 0 = 0x3ff0000000000000 := by
    cases arm64 <;> cases fused <;> decide
  change ldexp 64 (exp2Reduced (arm64 || fused) arm64 0 0) k = _
  rw [reduced]

set_option exponentiation.threshold 4096 in
set_option maxRecDepth 4096 in
theorem expmultiWord_scale_edges (arm64 fused : Bool) :
    expmultiWord arm64 fused 0 0 (-1075) = 0 ∧
      expmultiWord arm64 fused 0 0 (-1074) = 1 ∧
      expmultiWord arm64 fused 0 0 1024 = 0x7ff0000000000000 := by
  rw [expmultiWord_zero, expmultiWord_zero, expmultiWord_zero]
  decide

end GoLean.Math
