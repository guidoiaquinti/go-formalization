import Init.Data.Int.Lemmas
import Init.Omega

namespace GoLean.Math.Big

/-- Canonical signed value observable through Go's normalized Int. -/
def intSet (x : Int) : Int := x

def intAbs (x : Int) : Int := x.natAbs

def intNeg (x : Int) : Int := -x

def intAdd (x y : Int) : Int := x + y

def intSub (x y : Int) : Int := x - y

def intMul (x y : Int) : Int := x * y

def intCmp (x y : Int) : Int := if x < y then -1 else if x > y then 1 else 0

def intCmpAbs (x y : Int) : Int := intCmp x.natAbs y.natAbs

theorem intSet_preserves (x : Int) : intSet x = x := rfl

theorem intAbs_magnitude (x : Int) : intAbs x = x.natAbs := rfl

theorem intNeg_involution (x : Int) : intNeg (intNeg x) = x := by simp [intNeg]

theorem intAdd_cancel (x : Int) : intAdd x (intNeg x) = 0 := by
  unfold intAdd intNeg
  omega

theorem intSub_self (x : Int) : intSub x x = 0 := by simp [intSub]

theorem intMul_zero (x : Int) : intMul x 0 = 0 := by simp [intMul]

theorem intCmp_order (x y : Int) :
    (intCmp x y = -1 ↔ x < y) ∧ (intCmp x y = 1 ↔ x > y) := by
  unfold intCmp
  split
  · omega
  · split <;> omega

theorem intCmpAbs_sign (x y : Int) : intCmpAbs (-x) y = intCmpAbs x y := by
  simp [intCmpAbs]

end GoLean.Math.Big
