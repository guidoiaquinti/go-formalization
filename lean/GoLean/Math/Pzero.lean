import GoLean.Math.BesselRational

namespace GoLean.Math

/-- Go 1.27.0 src/math/j0.go:307–327. Selects the source interval,
    evaluates the rounded ratio, and adds one. None denotes the
    nil-coefficient-table panic below two or on NaN. -/
def pzeroWord (arm64 fused : Bool) (word : BitVec 64) : Option Nat :=
  let x := word.toNat
  if x < 0x4000000000000000 ∨ 0x7ff0000000000000 < x then none
  else if x = 0x7ff0000000000000 then some 0x3ff0000000000000
  else
    let (p, q) := besselCoefficients x
      ([
        0xc0b4850b36cc643d, 0xc0a36a6ecd4dcafc, 0xc07011027b19e863,
        0xc02029d0b44fa779, 0xbfb1fffffffffd32, 0x0000000000000000], [
        0x40e741774f2c49dc, 0x40fc810f8f9fa9bd, 0x40e3d2bb6eb6b05f,
        0x40adf37d50596938, 0x405d223307a96751])
      ([
        0xc075a6ef28a38bd7, 0xc074b3b36742cc63, 0xc050eb2f5a7d1783,
        0xc010a370f90c6bbf, 0xbfb1ffffe69afbc6, 0xbda918b147e495cc], [
        0x40a2cc1dc70be864, 0x40c2ccb8fa76fa38, 0x40b75af88fbe1d60,
        0x40906d025c7e2864, 0x404e60810c98c5de])
      ([
        0xc03f72aca892d80f, 0xc04d0a22420a1a45, 0xc035f74a4cb94e14,
        0xc00345b2aea48074, 0xbfb1fff6f7c0e24b, 0xbe25e1036fe1aa86], [
        0x4065b296fc379081, 0x40919ffcb8c39b7e, 0x4092a66e6d1061d6,
        0x40769839464a7c0e, 0x4041ed9284077dd3])
      ([
        0xc009de81af8fe70f, 0xc02662e6c5246303, 0xc01e8af3edafa7f3,
        0xbff736398a24a843, 0xbfb1ff62495e1e42, 0xbe77d316e927026d], [
        0x402d50b344391809, 0x40633c033ab6faff, 0x4070e78642ea079b,
        0x4061069e0ee8878f, 0x40363865908b5959])
    some (addWords 0x3ff0000000000000 (besselRational arm64 fused x p q))

theorem pzeroWord_panic (arm64 fused : Bool) (word : BitVec 64)
    (outside : word.toNat < 0x4000000000000000 ∨ 0x7ff0000000000000 < word.toNat) :
    pzeroWord arm64 fused word = none := by
  simp [pzeroWord, outside]

theorem pzeroWord_infinity (arm64 fused : Bool) :
    pzeroWord arm64 fused (BitVec.ofNat 64 0x7ff0000000000000) = some 0x3ff0000000000000 := by
  cases arm64 <;> cases fused <;> decide

end GoLean.Math
