import GoLean.Math.FmaZero

namespace GoLean.Math

/-- The FMA helper returns the one-bit indicator that an input has a set bit.
    https://github.com/golang/go/blob/go1.27.0/src/math/fma.go#L18-L25 -/
def fmaNonzero (word : BitVec 64) : BitVec 64 :=
  if word = 0 then 0 else 1

/-- Complementing the zero indicator yields the branchless expression
    documented beside the pinned implementation. -/
theorem fmaNonzero_spec (word : BitVec 64) :
    fmaNonzero word = 1 - ((((word >>> (1 : Nat)) ||| (word &&& 1)) - 1) >>> (63 : Nat)) := by
  rw [← fmaZero_spec, fmaNonzero, fmaZero]
  by_cases z : word = 0
  · rw [ite_eq_left z, ite_eq_left z]
    decide
  · rw [ite_eq_right z, ite_eq_right z]
    decide

end GoLean.Math
