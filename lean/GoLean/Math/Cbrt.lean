import GoLean.Math.Hypot

namespace GoLean.Math

private def scale : Nat := 2 ^ 1074

def signedUnits (word : Nat) : Int :=
  let units : Int := magnitudeUnits (word % 2 ^ 63)
  if 2 ^ 63 ≤ word then -units else units

def roundSigned (numerator : Int) (denominator : Nat) : Nat :=
  (if numerator < 0 then 2 ^ 63 else 0) + roundUnits numerator.natAbs denominator

def addWords (x y : Nat) : Nat :=
  roundSigned (signedUnits x + signedUnits y) 1

def multiplyWords (x y : Nat) : Nat :=
  roundSigned (signedUnits x * signedUnits y) scale

def divideWords (x y : Nat) : Nat :=
  let divisor := signedUnits y
  let numerator := signedUnits x * scale
  roundSigned (if divisor < 0 then -numerator else numerator) divisor.natAbs

def multiplyAddWords (fused : Bool) (x y z : Nat) : Nat :=
  if fused then roundSigned (signedUnits x * signedUnits y + signedUnits z * scale) scale
  else addWords (multiplyWords x y) z

private def cbrtMagnitude (fused : Bool) (x : Nat) : Nat :=
  let rough :=
    if x < 0x0010000000000000 then multiplyWords 0x4350000000000000 x / 3 + (696219795 <<< 32)
    else x / 3 + (715094163 <<< 32)
  let ratio := divideWords (multiplyWords rough rough) x
  let sum := multiplyAddWords fused ratio rough 0x3fe15f15f15f15f1
  let correction := addWords (addWords sum 0x3ff6a0ea0ea0ea0f) (divideWords 0xbfe691de2532c834 sum)
  let improved := multiplyWords rough (addWords 0x3fd6db6db6db6db7 (divideWords 0x3ff9b6db6db6db6e correction))
  let chopped := (improved &&& (0xffffffffc <<< 28)) + (1 <<< 30)
  let squared := multiplyWords chopped chopped
  let estimate := divideWords x squared
  let twice := addWords chopped chopped
  let delta := addWords estimate (chopped + 2 ^ 63)
  let step := divideWords delta (addWords twice estimate)
  multiplyAddWords fused chopped step chopped

/-- Go 1.27.0 src/math/cbrt.go:26–82. Binary64 stages follow the word
    estimate, rational correction, low-bit chopping and Newton update. The
    fused profile applies to the compiler-contracted multiply additions. -/
def cbrtWords (fused : Bool) (word : BitVec 64) : Nat :=
  let x := word.toNat % 2 ^ 63
  if x = 0 ∨ 0x7ff0000000000000 ≤ x then word.toNat
  else word.toNat - x + cbrtMagnitude fused x

theorem cbrtWords_special (fused : Bool) (word : BitVec 64)
    (h : word.toNat % 2 ^ 63 = 0 ∨ 0x7ff0000000000000 ≤ word.toNat % 2 ^ 63) :
    cbrtWords fused word = word.toNat := by
  simp [cbrtWords, h]

theorem cbrtWords_opposite (fused : Bool) (x y : BitVec 64)
    (positive : x.toNat < 2 ^ 63) (negative : 2 ^ 63 ≤ y.toNat)
    (magnitude : x.toNat % 2 ^ 63 = y.toNat % 2 ^ 63) :
    cbrtWords fused y = cbrtWords fused x + 2 ^ 63 := by
  have offset : y.toNat = x.toNat + 2 ^ 63 := by
    have bound := y.isLt
    omega
  unfold cbrtWords
  simp only [← magnitude]
  split <;> omega

end GoLean.Math
