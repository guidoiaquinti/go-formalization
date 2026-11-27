import GoLean.Math.Cbrt
import GoLean.Math.Ldexp

namespace GoLean.Math

def negateWord (word : Nat) : Nat := word ^^^ 2 ^ 63

def subtractWords (x y : Nat) : Nat :=
  addWords x (negateWord y)

def exp2Reduced (fused fuseSubtract : Bool) (hi lo : Nat) : Nat :=
  let r := subtractWords hi lo
  let t := multiplyWords r r
  let p4 := multiplyAddWords fused t 0x3e66376972bea4d0 0xbebbbd41c5d26bf1
  let p3 := multiplyAddWords fused t p4 0x3f11566aaf25de2c
  let p2 := multiplyAddWords fused t p3 0xbf66c16c16bebd93
  let p1 := multiplyAddWords fused t p2 0x3fc5555555555555
  let c := multiplyAddWords fuseSubtract t (negateWord p1) r
  let fraction := divideWords (multiplyWords r c) (subtractWords 0x4000000000000000 c)
  subtractWords 0x3ff0000000000000 (subtractWords (subtractWords lo fraction) hi)

private def exp2Exponent (raw : Nat) : Int :=
  let half := if raw < 2 ^ 63 then 0x3fe0000000000000 else 0xbfe0000000000000
  let integer := signedUnits (addWords raw half)
  integer.tdiv (2 ^ 1074)

def scaleExp2 (arm64 : Bool) (word : Nat) (k : Int) : Nat :=
  if arm64 then
    let exponent := (word / 2 ^ 52 : Int) + k
    let fraction := word % 2 ^ 52
    if 1 ≤ exponent then exponent.toNat * 2 ^ 52 + fraction
    else multiplyWords ((exponent + 52).toNat * 2 ^ 52 + fraction) 0x3cb0000000000000
  else ldexp 64 word k

/-- Go 1.27.0 src/math/exp.go:137–202 and src/math/exp_arm64.s:112–185.
    Models the reduction, split constants, polynomial, rational correction,
    and power-of-two scale. The arm64 assembly uses a fused polynomial and
    its own 52-bit subnormal staging. Generic amd64 v3/v4 contracts the
    polynomial additions and rounds the final multiply and subtract separately.
    The arm64 generic profile contracts both operations and uses generic scaling. -/
def exp2Word (arm64 fused fuseSubtract : Bool) (word : BitVec 64) : Nat :=
  let raw := word.toNat
  let magnitude := raw % 2 ^ 63
  if 0x7ff0000000000000 < magnitude then raw
  else if raw < 2 ^ 63 ∧ 0x408fffffffffffff < magnitude then 0x7ff0000000000000
  else if 2 ^ 63 ≤ raw ∧ 0x4090c80000000000 < magnitude then 0
  else
    let k := exp2Exponent raw
    let t := subtractWords raw (roundSigned (k * 2 ^ 1074) 1)
    let hi := multiplyWords t 0x3fe62e42fee00000
    let lo := multiplyWords (negateWord t) 0x3dea39ef35793c76
    scaleExp2 arm64 (exp2Reduced (arm64 || fused) (arm64 || fuseSubtract) hi lo) k

theorem exp2Word_nan (arm64 fused fuseSubtract : Bool) (word : BitVec 64)
    (nan : 0x7ff0000000000000 < word.toNat % 2 ^ 63) :
    exp2Word arm64 fused fuseSubtract word = word.toNat := by
  simp [exp2Word, nan]

theorem exp2Word_cutoff (arm64 fused fuseSubtract : Bool) (word : BitVec 64)
    (number : word.toNat % 2 ^ 63 ≤ 0x7ff0000000000000) :
    (2 ^ 63 ≤ word.toNat → 0x4090c80000000000 < word.toNat % 2 ^ 63 →
      exp2Word arm64 fused fuseSubtract word = 0) ∧
    (word.toNat < 2 ^ 63 → 0x408fffffffffffff < word.toNat % 2 ^ 63 →
      exp2Word arm64 fused fuseSubtract word = 0x7ff0000000000000) := by
  constructor
  · intro sign bound
    simp [exp2Word, Nat.not_lt.mpr number, Nat.not_lt.mpr sign, sign, bound]
  · intro sign bound
    simp [exp2Word, Nat.not_lt.mpr number, sign, bound]

set_option exponentiation.threshold 4096 in
set_option maxRecDepth 4096 in
theorem exp2Word_zeros (arm64 fused fuseSubtract : Bool) :
    exp2Word arm64 fused fuseSubtract (0 : BitVec 64) = 0x3ff0000000000000 ∧
      exp2Word arm64 fused fuseSubtract (BitVec.ofNat 64 (2 ^ 63)) = 0x3ff0000000000000 := by
  cases arm64 <;> cases fused <;> cases fuseSubtract <;> decide

end GoLean.Math
