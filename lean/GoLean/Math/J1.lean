import GoLean.Math.Sincos
import GoLean.Math.Cos
import GoLean.Math.Pone
import GoLean.Math.Qone
import GoLean.Math.Sqrt

namespace GoLean.Math

/-- Order-one Bessel phase with the source's double-angle correction. -/
def besselOnePhase (arm64 fused : Bool) (word : BitVec 64) : Option (Nat × Nat) := do
  let (sine, cosine) ← sincosWords arm64 fused word
  let difference := subtractWords sine cosine
  let negativeSum := subtractWords (negateWord sine) cosine
  let x := word.toNat
  if x < 0x7fdfffffffffffff then
    let doubleCosine ← cosWord arm64 fused (BitVec.ofNat 64 (addWords x x))
    let product := multiplyWords sine cosine
    if 0 < product && product < 2 ^ 63 then
      some (negativeSum, divideWords doubleCosine negativeSum)
    else some (divideWords doubleCosine difference, difference)
  else some (negativeSum, difference)

private def j1Large (arm64 fused : Bool) (word : BitVec 64) : Option Nat := do
  let (sine, cosine) ← besselOnePhase arm64 fused word
  let root ← sqrtWord word
  let phase ←
    if 0x4800000000000000 < word.toNat then some cosine
    else (do
      let u ← poneWord arm64 fused word
      let v ← qoneWord arm64 fused word
      some (multiplyAddWords arm64 (negateWord v) sine (multiplyWords u cosine)))
  some (divideWords (multiplyWords 0x3fe20dd750429b6d phase) root)

private def j1Small (fused : Bool) (x : Nat) : Nat :=
  let square := multiplyWords x x
  let fma := multiplyAddWords fused square
  let numerator := fma
    (fma (fma 0x3e6aaafa46ca0bd9 0xbef0c5c6ba169668) 0x3f570d9f98472c61)
    0xbfb0000000000000
  let denominator := fma
    (fma (fma (fma (fma 0x3dab2acfcfb97ed8 0x3e35ac88c97dff2c) 0x3eb3bff8333f8498)
      0x3f285f56b9cdf664) 0x3f939d0b12637e53) 0x3ff0000000000000
  let ratio := divideWords (multiplyWords (multiplyWords square numerator) x) denominator
  multiplyAddWords fused 0x3fe0000000000000 x ratio

/-- Go 1.27.0 src/math/j1.go:75–146. None denotes the output NaN class.
    The tiny branch returns from the magnitude before restoring the sign. -/
def j1Word (arm64 fused : Bool) (word : BitVec 64) : Option Nat :=
  let x := word.toNat % 2 ^ 63
  if 0x7ff0000000000000 < x then none
  else if x = 0x7ff0000000000000 ∨ x = 0 then some 0
  else if x < 0x3e40000000000000 then some (multiplyWords 0x3fe0000000000000 x)
  else
    let result := if x < 0x4000000000000000 then some (j1Small (arm64 || fused) x)
      else j1Large arm64 fused (BitVec.ofNat 64 x)
    result.map fun y => if 2 ^ 63 ≤ word.toNat then negateWord y else y

theorem j1Word_specialZero (arm64 fused : Bool) (word : BitVec 64)
    (special : word.toNat % 2 ^ 63 = 0x7ff0000000000000 ∨ word.toNat % 2 ^ 63 = 0) :
    j1Word arm64 fused word = some 0 := by
  simp [j1Word, special, show ¬0x7ff0000000000000 < word.toNat % 2 ^ 63 by omega]

theorem j1Word_tinyMagnitude (arm64 fused : Bool) (word : BitVec 64)
    (nonzero : word.toNat % 2 ^ 63 ≠ 0)
    (tiny : word.toNat % 2 ^ 63 < 0x3e40000000000000) :
    j1Word arm64 fused word = some (multiplyWords 0x3fe0000000000000 (word.toNat % 2 ^ 63)) := by
  simp [j1Word, tiny, nonzero,
    show ¬0x7ff0000000000000 < word.toNat % 2 ^ 63 by omega,
    show word.toNat % 2 ^ 63 ≠ 0x7ff0000000000000 by omega]

end GoLean.Math
