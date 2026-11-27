import GoLean.Math.Exp
import GoLean.Math.Log
import GoLean.Math.IsOddInt
import GoLean.Math.Sqrt

namespace GoLean.Math

private def powZero (base : Nat) (exponent : BitVec 64) (negative : Bool) : Nat :=
  let sign := if 2 ^ 63 ≤ base ∧ isOddInt exponent then 2 ^ 63 else 0
  sign + if negative then 0x7ff0000000000000 else 0

private def powSquares : Nat → Nat → Int → Nat → Int → Nat × Int
  | 0, _, _, product, scale => (product, scale)
  | i + 1, fraction, exponent, product, scale =>
    if exponent < -4096 ∨ 4096 < exponent then (product, scale + exponent)
    else
      let product := if (i + 1) % 2 = 1 then multiplyWords product fraction else product
      let scale := if (i + 1) % 2 = 1 then scale + exponent else scale
      let squared := multiplyWords fraction fraction
      let shift := squared < 0x3fe0000000000000
      let fraction := if shift then addWords squared squared else squared
      let exponent := 2 * exponent - if shift then 1 else 0
      powSquares ((i + 1) / 2) fraction exponent product scale
termination_by i => i
decreasing_by omega

private def powFinite (arm64 fused : Bool) (x units : Nat) (negative : Bool) : Option Nat :=
  let scale := 2 ^ 1074
  let whole := units / scale
  let remainder := units % scale
  if remainder ≠ 0 ∧ 2 ^ 63 ≤ x then none
  else if 2 ^ 63 ≤ whole then
    if x = 0xbff0000000000000 then some 0x3ff0000000000000
    else if decide (x % 2 ^ 63 < 0x3ff0000000000000) == !negative then some 0
    else some 0x7ff0000000000000
  else
    let (whole, fraction) :=
      if scale / 2 < remainder then (whole + 1, (remainder : Int) - scale) else (whole, (remainder : Int))
    let fractionPower :=
      if fraction = 0 then 0x3ff0000000000000
      else
        let logarithm := (logWord arm64 (BitVec.ofNat 64 x)).getD 0
        (expWord arm64 fused (BitVec.ofNat 64 (multiplyWords (roundSigned fraction 1) logarithm)))
    let (fraction, exponent) := frexp x
    let (product, exponent) := powSquares whole fraction exponent fractionPower 0
    if negative then some (ldexp 64 (divideWords 0x3ff0000000000000 product) (-exponent))
    else some (ldexp 64 product exponent)

/-- Go 1.27.0 src/math/pow.go:48–158. Preserve ordered exceptions,
    integral exponent parity, the half powers, fractional correction and
    repeated significand squaring with a separate exponent accumulator. -/
def powWords (arm64 fused : Bool) (x y : BitVec 64) : Option Nat :=
  let base := x.toNat
  let exponent := y.toNat
  let xmag := base % 2 ^ 63
  let ymag := exponent % 2 ^ 63
  let negative := decide (2 ^ 63 ≤ exponent)
  if ymag = 0 ∨ base = 0x3ff0000000000000 then some 0x3ff0000000000000
  else if exponent = 0x3ff0000000000000 then
    if 0x7ff0000000000000 < xmag then none else some base
  else if 0x7ff0000000000000 < xmag ∨ 0x7ff0000000000000 < ymag then none
  else if xmag = 0 then some (powZero base y negative)
  else if ymag = 0x7ff0000000000000 then
    if xmag = 0x3ff0000000000000 then some 0x3ff0000000000000
    else if decide (xmag < 0x3ff0000000000000) == !negative then some 0
    else some 0x7ff0000000000000
  else if xmag = 0x7ff0000000000000 then
    some (powZero (base - xmag) y (!negative))
  else if exponent = 0x3fe0000000000000 then sqrtWord x
  else if exponent = 0xbfe0000000000000 then (sqrtWord x).map (divideWords 0x3ff0000000000000)
  else powFinite arm64 fused base (magnitudeUnits ymag) negative

theorem powWords_zeroExponent (arm64 fused : Bool) (x y : BitVec 64)
    (zero : y.toNat % 2 ^ 63 = 0) :
    powWords arm64 fused x y = some 0x3ff0000000000000 := by
  simp [powWords, zero]

theorem powWords_oneBase (arm64 fused : Bool) (x y : BitVec 64)
    (one : x.toNat = 0x3ff0000000000000) :
    powWords arm64 fused x y = some 0x3ff0000000000000 := by
  simp [powWords, one]

theorem powWords_oneExponent (arm64 fused : Bool) (x y : BitVec 64)
    (one : y.toNat = 0x3ff0000000000000)
    (notNaN : x.toNat % 2 ^ 63 ≤ 0x7ff0000000000000) :
    powWords arm64 fused x y = some x.toNat := by
  simp [powWords, one, Nat.not_lt.mpr notNaN, eq_comm]

theorem powWords_signedZeroBase (arm64 fused : Bool) (x y : BitVec 64)
    (base : x.toNat = 0x8000000000000000) (nonzero : y.toNat % 2 ^ 63 ≠ 0)
    (notOne : y.toNat ≠ 0x3ff0000000000000)
    (notNaN : y.toNat % 2 ^ 63 ≤ 0x7ff0000000000000) :
    powWords arm64 fused x y =
      some ((if isOddInt y then 2 ^ 63 else 0) +
        (if 2 ^ 63 ≤ y.toNat then 0x7ff0000000000000 else 0)) := by
  simp [powWords, powZero, base, nonzero, notOne, Nat.not_lt.mpr notNaN]

end GoLean.Math
