import GoLean.Math.BesselRational

namespace GoLean.Math

/-- Go 1.27.0 src/math/j1.go:302–322. Selects the source interval,
    evaluates the rounded ratio, and adds one. None denotes the
    nil-coefficient-table panic below two or on NaN. -/
def poneWord (arm64 fused : Bool) (word : BitVec 64) : Option Nat :=
  let x := word.toNat
  if x < 0x4000000000000000 ∨ 0x7ff0000000000000 < x then none
  else if x = 0x7ff0000000000000 then some 0x3ff0000000000000
  else
    let (p, q) := besselCoefficients x
      ([
        0x40beea7ac32782dd, 0x40ae457da3a532cc, 0x4079c0d4652ea590,
        0x402a7a9d357f7fce, 0x3fbdfffffffffcce, 0x0000000000000000], [
        0x40de1511697a0b2d, 0x40f7d42cb28f17bb, 0x40e20b8697c5bb7f,
        0x40ac85dc964d274f, 0x405c8d458e656cac])
      ([
        0x408085b8bb7e0cb7, 0x40802d16d052d649, 0x405b13b9452602ed,
        0x401b36046e6315e3, 0x3fbdffffe2c10043, 0x3dad0667dae1ca7d], [
        0x40978030036f5e51, 0x40bea4b0b8a5bb15, 0x40b4e9445706b6fb,
        0x408efb361b066701, 0x404da3eaa8af633d])
      ([
        0x4048478f8ea83ee5, 0x4056c3854d2c1837, 0x40418f489da6d129,
        0x400f76bce85ead8a, 0x3fbdfff55b21d17b, 0x3e29fc21a7ad9edd], [
        0x4059f26d7c2eed53, 0x408bd67da32e31e9, 0x40905b7c5037d523,
        0x40750c3307f1a75f, 0x40416549a134069c])
      ([
        0x40144b49a574c1fe, 0x4031b1a8177f8ee2, 0x40287c377f71a964,
        0x4002f2b7f98faec0, 0x3fbdff42be760d83, 0x3e7ce9d4f65544f4], [
        0x4020bab1f44e5192, 0x405d6b7ada1884a9, 0x406d08d8d5a2dbd9,
        0x405f529314f92cd5, 0x40356fbd8ad5ecdc])
    some (addWords 0x3ff0000000000000 (besselRational arm64 fused x p q))

theorem poneWord_panic (arm64 fused : Bool) (word : BitVec 64)
    (outside : word.toNat < 0x4000000000000000 ∨ 0x7ff0000000000000 < word.toNat) :
    poneWord arm64 fused word = none := by
  simp [poneWord, outside]

theorem poneWord_infinity (arm64 fused : Bool) :
    poneWord arm64 fused (BitVec.ofNat 64 0x7ff0000000000000) = some 0x3ff0000000000000 := by
  cases arm64 <;> cases fused <;> decide

end GoLean.Math
