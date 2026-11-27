import GoLean.Math.BesselRational

namespace GoLean.Math

/-- Go 1.27.0 src/math/j1.go:405–424. Selects the source interval,
    evaluates the rounded ratio, adds 3/8, and divides by x. None denotes the
    nil-coefficient-table panic below two or on NaN. -/
def qoneWord (arm64 fused : Bool) (word : BitVec 64) : Option Nat :=
  let x := word.toNat
  if x < 0x4000000000000000 ∨ 0x7ff0000000000000 < x then none
  else if x = 0x7ff0000000000000 then some 0
  else
    let (p, q) := besselCoefficients x
      ([
        0xc0e7a6d065d09c6a, 0xc0c724e740f87415, 0xc087bcd053e4b576,
        0xc0304591a26779f7, 0xbfba3ffffffffdf3, 0x0000000000000000], [
        0xc111f9690ea5aa18, 0x412457d27719ad5c, 0x4125f65372869c19,
        0x4100579ab0b75e98, 0x40be9162d0d88419, 0x40642ca6de5bcde5])
      ([
        0xc0a468e388fda79d, 0xc09574c66931734f, 0xc066f56d6ca7b9b0,
        0xc0201ce6ca03ad4b, 0xbfba3fffcb597fef, 0xbdb6fa431aa1a098], [
        0xc0b26f2efcffa004, 0x40db4b04cf7c364b, 0x40e8576daabad197,
        0x40d10f1f0d64ce29, 0x409f1f31e77bf839, 0x405451b2ff5a11b2])
      ([
        0xc06b66b95f5c1bf6, 0xc06c87d34718d55f, 0xc04cec71c25d16da,
        0xc01270c23302d9ff, 0xbfba3feb51aeed54, 0xbe35cfa9d38fc84f], [
        0xc060e670290a311f, 0x409dbc7a0dd4df4b, 0x40b5abbaa61d54a6,
        0x40aa684e448e7c9a, 0x40850eebc031ee3e, 0x4047d523ccd367e4])
      ([
        0xc0355f3639cf6e52, 0xc04529a3de104aaa, 0xc033a9e2c168907f,
        0xc006048469bb4eda, 0xbfba3e8e9148b010, 0xbe87f12644c626d2], [
        0xc013d686e71be86b, 0x40637e5e3c3ed8d4, 0x40871b2548d4c029,
        0x4087ac05ce49a0f7, 0x406f9f68db821cba, 0x403d888a78ae64ff])
    let ratio := besselRational arm64 fused x p q
    some (divideWords (addWords 0x3fd8000000000000 ratio) x)

theorem qoneWord_panic (arm64 fused : Bool) (word : BitVec 64)
    (outside : word.toNat < 0x4000000000000000 ∨ 0x7ff0000000000000 < word.toNat) :
    qoneWord arm64 fused word = none := by
  simp [qoneWord, outside]

theorem qoneWord_infinity (arm64 fused : Bool) :
    qoneWord arm64 fused (BitVec.ofNat 64 0x7ff0000000000000) = some 0 := by
  cases arm64 <;> cases fused <;> decide

end GoLean.Math
