import GoLean.Math.J1
import GoLean.Math.Log

namespace GoLean.Math

private def y1Large (arm64 fused : Bool) (word : BitVec 64) : Option Nat := do
  let (sine, cosine) ← besselOnePhase arm64 fused word
  let root ← sqrtWord word
  let phase ←
    if 0x4800000000000000 < word.toNat then some sine
    else (do
      let u ← poneWord arm64 fused word
      let v ← qoneWord arm64 fused word
      some (if arm64 then multiplyAddWords true u sine (multiplyWords v cosine)
        else multiplyAddWords fused v cosine (multiplyWords u sine)))
  some (divideWords (multiplyWords 0x3fe20dd750429b6d phase) root)

private def y1Small (arm64 compileFused : Bool) (word : BitVec 64) : Option Nat := do
  let x := word.toNat
  let factor := 0x3fe45f306dc9c883
  if x ≤ 0x3c90000000000000 then some (divideWords (negateWord factor) x)
  else
    let logarithm ← logWord arm64 word
    let j ← j1Word arm64 compileFused word
    let fused := arm64 || compileFused
    let square := multiplyWords x x
    let fma := multiplyAddWords fused square
    let numerator := fma
      (fma (fma (fma 0xbe78ac00569105b8 0x3ef8ab038fa6b88e) 0xbf5f55e54844f50f)
        0x3fa9d3c776292cd1) 0xbfc91866143cbc8a
    let denominator := fma
      (fma (fma (fma (fma 0x3db25039daca772a 0x3e3abf1d5ba69a86) 0x3eb6c05a894e8ca6)
        0x3f2a8c896c257764) 0x3f94650d3f4da9f0) 0x3ff0000000000000
    let ratio := divideWords numerator denominator
    let reciprocal := divideWords 0x3ff0000000000000 x
    let logCorrection := multiplyAddWords arm64 j logarithm (negateWord reciprocal)
    some (if arm64 then multiplyAddWords true x ratio (multiplyWords factor logCorrection)
      else multiplyAddWords compileFused factor logCorrection (multiplyWords x ratio))

/-- Go 1.27.0 src/math/j1.go:156–223. None denotes the output NaN class. -/
def y1Word (arm64 fused : Bool) (word : BitVec 64) : Option Nat :=
  let x := word.toNat
  let magnitude := x % 2 ^ 63
  if magnitude = 0 then some 0xfff0000000000000
  else if 2 ^ 63 ≤ x ∨ 0x7ff0000000000000 < magnitude then none
  else if x = 0x7ff0000000000000 then some 0
  else if x < 0x4000000000000000 then y1Small arm64 fused word
  else y1Large arm64 fused word

theorem y1Word_zero (arm64 fused : Bool) (word : BitVec 64)
    (zero : word.toNat % 2 ^ 63 = 0) :
    y1Word arm64 fused word = some 0xfff0000000000000 := by
  simp [y1Word, zero]

theorem y1Word_negative (arm64 fused : Bool) (word : BitVec 64)
    (negative : 2 ^ 63 ≤ word.toNat) (nonzero : word.toNat % 2 ^ 63 ≠ 0) :
    y1Word arm64 fused word = none := by
  simp [y1Word, nonzero, negative]

theorem y1Word_inf (arm64 fused : Bool) :
    y1Word arm64 fused (BitVec.ofNat 64 0x7ff0000000000000) = some 0 := rfl

end GoLean.Math
