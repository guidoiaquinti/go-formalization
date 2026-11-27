import GoLean.Math.J0
import GoLean.Math.Log

namespace GoLean.Math

private def y0Large (arm64 fused : Bool) (word : BitVec 64) : Option Nat := do
  let (difference, sum) ← besselZeroPhase arm64 fused word
  let root ← sqrtWord word
  let phase ←
    if 0x4800000000000000 < word.toNat then some difference
    else (do
      let u ← pzeroWord arm64 fused word
      let v ← qzeroWord arm64 fused word
      some (if arm64 then multiplyAddWords true v sum (multiplyWords u difference)
        else multiplyAddWords fused u difference (multiplyWords v sum)))
  some (divideWords (multiplyWords 0x3fe20dd750429b6d phase) root)

private def y0Small (arm64 compileFused : Bool) (word : BitVec 64) : Option Nat := do
  let logarithm ← logWord arm64 word
  let fused := arm64 || compileFused
  let factor := 0x3fe45f306dc9c883
  let constant := 0xbfb2e4d699cbd01f
  if word.toNat ≤ 0x3e40000000000000 then
    some (multiplyAddWords fused factor logarithm constant)
  else
    let j ← j0Word arm64 compileFused word
    let square := multiplyWords word.toNat word.toNat
    let fma := multiplyAddWords fused square
    let numerator := fma
      (fma (fma (fma (fma (fma 0xbdc5e43d693fb3c8 0x3e5500573b4eabd4)
        0xbecffea773d25cad) 0x3f36c54d20b29b6b) 0xbf8c4ce8b16cfa97)
        0x3fc69d019de9e3fc) constant
    let denominator := fma
      (fma (fma (fma 0x3dfe50183bd6d9ef 0x3e91642d7ff202fd)
        0x3f13ecbbf578c6c1) 0x3f8a127091c9c71a) 0x3ff0000000000000
    some (multiplyAddWords arm64 (multiplyWords factor j) logarithm (divideWords numerator denominator))

/-- Go 1.27.0 src/math/j0.go:156–228. None denotes the output NaN class. -/
def y0Word (arm64 fused : Bool) (word : BitVec 64) : Option Nat :=
  let x := word.toNat
  let magnitude := x % 2 ^ 63
  if magnitude = 0 then some 0xfff0000000000000
  else if 2 ^ 63 ≤ x ∨ 0x7ff0000000000000 < magnitude then none
  else if x = 0x7ff0000000000000 then some 0
  else if x < 0x4000000000000000 then y0Small arm64 fused word
  else y0Large arm64 fused word

theorem y0Word_zero (arm64 fused : Bool) (word : BitVec 64)
    (zero : word.toNat % 2 ^ 63 = 0) :
    y0Word arm64 fused word = some 0xfff0000000000000 := by
  simp [y0Word, zero]

theorem y0Word_negative (arm64 fused : Bool) (word : BitVec 64)
    (negative : 2 ^ 63 ≤ word.toNat) (nonzero : word.toNat % 2 ^ 63 ≠ 0) :
    y0Word arm64 fused word = none := by
  simp [y0Word, nonzero, negative]

theorem y0Word_inf (arm64 fused : Bool) :
    y0Word arm64 fused (BitVec.ofNat 64 0x7ff0000000000000) = some 0 := rfl

end GoLean.Math
