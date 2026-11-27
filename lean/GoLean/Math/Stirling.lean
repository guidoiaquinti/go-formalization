import GoLean.Math.Pow

namespace GoLean.Math

/-- Go 1.27.0 src/math/gamma.go:94–119, for the positive domain used
    by Gamma. Separate factors retain the caller's overflow control. -/
def stirlingWords (arm64 compileFused expFused : Bool) (word : BitVec 64) :
    Option (Nat × Nat) :=
  let x := word.toNat
  if x < 0x4040800000000000 ∨ 0x7ff0000000000000 < x then none
  else if 0x4069000000000000 < x then some (0x7ff0000000000000, 0x3ff0000000000000)
  else do
    let reciprocal := divideWords 0x3ff0000000000000 x
    let p := multiplyAddWords arm64 0x3f49cc72592d7293 reciprocal 0xbf2e166b27e61d7c
    let p := multiplyAddWords arm64 p reciprocal 0xbf65f72607d44fd7
    let p := multiplyAddWords arm64 p reciprocal 0x3f6c71c71b98c5fd
    let p := multiplyAddWords arm64 p reciprocal 0x3fb5555555555986
    let correction := multiplyAddWords (arm64 || compileFused) reciprocal p 0x3ff0000000000000
    let leading := multiplyWords 0x40040d931ff62706 correction
    let exponential := expWord arm64 expFused word
    if 0x4061e083ba3443d4 < x then
      let power := multiplyAddWords arm64 x 0x3fe0000000000000 0xbfd0000000000000
      let second ← powWords arm64 expFused word (BitVec.ofNat 64 power)
      return (second, multiplyWords leading (divideWords second exponential))
    else
      let power := subtractWords x 0x3fe0000000000000
      let first ← powWords arm64 expFused word (BitVec.ofNat 64 power)
      return (divideWords first exponential, leading)

theorem stirlingWords_saturated (arm64 compileFused expFused : Bool) (word : BitVec 64)
    (large : 0x4069000000000000 < word.toNat)
    (valid : word.toNat ≤ 0x7ff0000000000000) :
    stirlingWords arm64 compileFused expFused word =
      some (0x7ff0000000000000, 0x3ff0000000000000) := by
  have above : ¬ word.toNat < 0x4040800000000000 := by omega
  simp [stirlingWords, above, Nat.not_lt.mpr valid, large]

theorem stirlingWords_outside (arm64 compileFused expFused : Bool) (word : BitVec 64)
    (outside : word.toNat < 0x4040800000000000 ∨ 0x7ff0000000000000 < word.toNat) :
    stirlingWords arm64 compileFused expFused word = none := by
  simp [stirlingWords, outside]

end GoLean.Math
