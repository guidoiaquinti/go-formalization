import GoLean.Math.Exp2

namespace GoLean.Math

/-- The pinned small-argument rational kernel. Its denominator begins
    with x*x + Q0, which the fused profiles evaluate with one rounding. -/
def atanSmallWord (fused : Bool) (x : Nat) : Nat :=
  let square := multiplyWords x x
  let p := multiplyAddWords fused 0xbfec007fa1f72594 square 0xc03028545b6b807a
  let p := multiplyAddWords fused p square 0xc052c08c36880273
  let p := multiplyAddWords fused p square 0xc05eb8bf2d05ba25
  let p := multiplyAddWords fused p square 0xc0503669fd28ec8e
  let q := multiplyAddWords fused x x 0x4038dbc45b14603c
  let q := multiplyAddWords fused q square 0x4064a0dd43b8fa25
  let q := multiplyAddWords fused q square 0x407b0e18d2e2be3b
  let q := multiplyAddWords fused q square 0x407e563f13b049ea
  let q := multiplyAddWords fused q square 0x4068519efbbd62ec
  let correction := divideWords (multiplyWords square p) q
  multiplyAddWords fused x correction x

/-- Reduce a nonnegative magnitude through the pinned thresholds and
    add the split π/4 or π/2 constant with separate roundings. -/
def atanPositiveWord (fused : Bool) (magnitude : Nat) : Nat :=
  if magnitude ≤ 0x3fe51eb851eb851f then atanSmallWord fused magnitude
  else if 0x4003504f333f9de6 < magnitude then
    let reciprocal := if magnitude = 0x7ff0000000000000 then 0
      else divideWords 0x3ff0000000000000 magnitude
    addWords (subtractWords 0x3ff921fb54442d18 (atanSmallWord fused reciprocal))
      0x3c91a62633145c07
  else
    let reduced := divideWords (subtractWords magnitude 0x3ff0000000000000)
      (addWords magnitude 0x3ff0000000000000)
    addWords (addWords 0x3fe921fb54442d18 (atanSmallWord fused reduced)) 0x3c81a62633145c07

/-- Go 1.27.0 src/math/atan.go:53–112. Preserve the three reduction
    branches and reconstruction rounding. None represents NaN by class. -/
def atanWord (fused : Bool) (word : BitVec 64) : Option Nat :=
  let raw := word.toNat
  let magnitude := raw % 2 ^ 63
  if magnitude = 0 then some raw
  else if 0x7ff0000000000000 < magnitude then none
  else
    let result := atanPositiveWord fused magnitude
    some (if 2 ^ 63 ≤ raw then negateWord result else result)

theorem atanWord_zero (fused : Bool) (word : BitVec 64)
    (zero : word.toNat % 2 ^ 63 = 0) : atanWord fused word = some word.toNat := by
  simp [atanWord, zero]

theorem atanWord_nan (fused : Bool) (word : BitVec 64)
    (nan : 0x7ff0000000000000 < word.toNat % 2 ^ 63) : atanWord fused word = none := by
  simp [atanWord, nan, show word.toNat % 2 ^ 63 ≠ 0 by omega]

theorem atanWord_opposite (fused : Bool) (x y : BitVec 64)
    (positive : x.toNat < 2 ^ 63) (negative : 2 ^ 63 ≤ y.toNat)
    (same : x.toNat % 2 ^ 63 = y.toNat % 2 ^ 63)
    (notNaN : x.toNat % 2 ^ 63 ≤ 0x7ff0000000000000)
    (nonzero : x.toNat % 2 ^ 63 ≠ 0) :
    atanWord fused y = (atanWord fused x).map negateWord := by
  simp [atanWord, ← same, nonzero, Nat.not_lt.mpr notNaN, Nat.not_le.mpr positive, negative]

end GoLean.Math
