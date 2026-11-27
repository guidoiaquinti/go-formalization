import GoLean.Math.Atan

namespace GoLean.Math

/-- Go 1.27.0 src/math/atan2.go:37–76. Preserve special-case order and
    quadrant reconstruction from the rounded quotient. The latter can
    select positive π for a lower-left input whose quotient rounds to zero. -/
def atan2Word (fused : Bool) (y x : BitVec 64) : Option Nat :=
  let ymag := y.toNat % 2 ^ 63
  let xmag := x.toNat % 2 ^ 63
  let negativeX := 2 ^ 63 ≤ x.toNat
  let signed := fun word => if 2 ^ 63 ≤ y.toNat then negateWord word else word
  let inf := 0x7ff0000000000000
  let pi := 0x400921fb54442d18
  let halfPi := 0x3ff921fb54442d18
  if inf < ymag ∨ inf < xmag then none
  else if ymag = 0 then some (signed (if negativeX then pi else 0))
  else if xmag = 0 then some (signed halfPi)
  else if xmag = inf then
    if negativeX then
      some (signed (if ymag = inf then 0x4002d97c7f3321d2 else pi))
    else
      some (signed (if ymag = inf then 0x3fe921fb54442d18 else 0))
  else if ymag = inf then some (signed halfPi)
  else do
    let q ← atanWord fused (BitVec.ofNat 64 (divideWords y.toNat x.toNat))
    if negativeX then
      return if q = 0 ∨ 2 ^ 63 ≤ q then addWords q pi else subtractWords q pi
    else
      return q

theorem atan2Word_nan (fused : Bool) (y x : BitVec 64)
    (nan : 0x7ff0000000000000 < y.toNat % 2 ^ 63 ∨
      0x7ff0000000000000 < x.toNat % 2 ^ 63) :
    atan2Word fused y x = none := by
  simp [atan2Word, nan]

theorem atan2Word_zeroY (fused : Bool) (y x : BitVec 64)
    (zero : y.toNat % 2 ^ 63 = 0)
    (notNaN : x.toNat % 2 ^ 63 ≤ 0x7ff0000000000000) :
    atan2Word fused y x =
      some (let angle := if 2 ^ 63 ≤ x.toNat then 0x400921fb54442d18 else 0
            if 2 ^ 63 ≤ y.toNat then negateWord angle else angle) := by
  simp [atan2Word, zero, Nat.not_lt.mpr notNaN]

set_option exponentiation.threshold 4096 in
set_option maxRecDepth 4096 in
theorem atan2Word_lowerLeftZero (fused : Bool) (y x : BitVec 64)
    (yNonzero : y.toNat % 2 ^ 63 ≠ 0) (xNonzero : x.toNat % 2 ^ 63 ≠ 0)
    (yFinite : y.toNat % 2 ^ 63 < 0x7ff0000000000000)
    (xFinite : x.toNat % 2 ^ 63 < 0x7ff0000000000000)
    (xNegative : 2 ^ 63 ≤ x.toNat) (underflow : divideWords y.toNat x.toNat = 0) :
    atan2Word fused y x = some 0x400921fb54442d18 := by
  have sum : addWords 0 0x400921fb54442d18 = 0x400921fb54442d18 := by decide
  simp [atan2Word, yNonzero, xNonzero, Nat.not_lt.mpr (Nat.le_of_lt yFinite),
    Nat.not_lt.mpr (Nat.le_of_lt xFinite), Nat.ne_of_lt yFinite, Nat.ne_of_lt xFinite,
    xNegative, underflow, atanWord, sum]

end GoLean.Math
