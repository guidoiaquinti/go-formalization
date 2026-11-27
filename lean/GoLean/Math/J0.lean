import GoLean.Math.Sincos
import GoLean.Math.Cos
import GoLean.Math.Pzero
import GoLean.Math.Qzero
import GoLean.Math.Sqrt

namespace GoLean.Math

/-- The shared order-zero Bessel phase. Correct a cancellation-prone sum or
    difference with the double angle identity before applying the amplitudes. -/
def besselZeroPhase (arm64 fused : Bool) (word : BitVec 64) : Option (Nat × Nat) := do
  let (sine, cosine) ← sincosWords arm64 fused word
  let difference := addWords sine (negateWord cosine)
  let sum := addWords sine cosine
  let x := word.toNat
  if x < 0x7fdfffffffffffff then
    let doubleCosine ← cosWord arm64 fused (BitVec.ofNat 64 (addWords x x))
    let product := multiplyWords sine cosine
    let opposite := 2 ^ 63 < product
    let correction := negateWord doubleCosine
    if opposite then some (difference, divideWords correction difference)
    else some (divideWords correction sum, sum)
  else some (difference, sum)

private def j0Large (arm64 fused : Bool) (word : BitVec 64) : Option Nat := do
  let (difference, sum) ← besselZeroPhase arm64 fused word
  let root ← sqrtWord word
  let phase ←
    if 0x4800000000000000 < word.toNat then some sum
    else (do
      let u ← pzeroWord arm64 fused word
      let v ← qzeroWord arm64 fused word
      some (multiplyAddWords arm64 (negateWord v) difference (multiplyWords u sum)))
  some (divideWords (multiplyWords 0x3fe20dd750429b6d phase) root)

private def j0Small (arm64 compileFused : Bool) (x : Nat) : Nat :=
  let one := 0x3ff0000000000000
  if x < 0x3e40000000000000 then one
  else if x < 0x3f20000000000000 then
    multiplyAddWords arm64 (multiplyWords 0x3fd0000000000000 x) (negateWord x) one
  else
    let fused := arm64 || compileFused
    let square := multiplyWords x x
    let fma := multiplyAddWords fused square
    let numerator := multiplyWords square (fma
      (fma (fma 0xbe33d5e773d63fce 0x3ebeb1d10c503919) 0xbf28e6a5b61ac6e9)
      0x3f8ffffffffffffd)
    let denominator := fma
      (fma (fma (fma 0x3e1408bcf4745d8f 0x3ea13b54ce84d5a9) 0x3f1ea6d2dd57dbf4)
        0x3f8ffce882c8c2a4) one
    let ratio := divideWords numerator denominator
    if x < one then fma (addWords 0xbfd0000000000000 ratio) one
    else
      let half := multiplyWords 0x3fe0000000000000 x
      let left := addWords one half
      let right := addWords one (negateWord half)
      if arm64 then fma ratio (multiplyWords left right)
      else multiplyAddWords compileFused left right (multiplyWords square ratio)

/-- Go 1.27.0 src/math/j0.go:77–147. None denotes the output NaN class. -/
def j0Word (arm64 fused : Bool) (word : BitVec 64) : Option Nat :=
  let x := word.toNat % 2 ^ 63
  if 0x7ff0000000000000 < x then none
  else if x = 0x7ff0000000000000 then some 0
  else if x < 0x4000000000000000 then some (j0Small arm64 fused x)
  else j0Large arm64 fused (BitVec.ofNat 64 x)

theorem j0Word_even (arm64 fused : Bool) (a b : BitVec 64)
    (same : a.toNat % 2 ^ 63 = b.toNat % 2 ^ 63) :
    j0Word arm64 fused a = j0Word arm64 fused b := by
  simp [j0Word, same]

theorem j0Word_nan (arm64 fused : Bool) (word : BitVec 64)
    (nan : 0x7ff0000000000000 < word.toNat % 2 ^ 63) :
    j0Word arm64 fused word = none := by
  simp [j0Word, nan]

end GoLean.Math
