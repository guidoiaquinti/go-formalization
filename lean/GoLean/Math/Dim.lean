import Init.Data.Float.Float

namespace GoLean.Math

/-- Go subtracts before testing for a nonpositive result. Keeping that
    comparison direction lets a NaN result pass through.
    https://github.com/golang/go/blob/go1.27.0/src/math/dim.go#L7-L27 -/
def dim (x y : Float) : Float :=
  let difference := x - y
  if difference ≤ 0 then 0 else difference

/-- Clamp a nonpositive subtraction to the zero literal. Whenever the
    comparison is false, keep every bit of the subtraction result. -/
theorem dim_spec (x y : Float) :
    ((x - y) ≤ 0 → (dim x y).toBits = (0 : Float).toBits) ∧
      (¬ ((x - y) ≤ 0) → (dim x y).toBits = (x - y).toBits) := by
  constructor
  · intro nonpositive
    simp [dim, nonpositive]
  · intro compareFalse
    simp [dim, compareFalse]

end GoLean.Math
