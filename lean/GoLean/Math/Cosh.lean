import GoLean.Math.Exp

namespace GoLean.Math

/-- Go 1.27.0 src/math/sinh.go:79–93. The large branch halves the
    already-rounded Exp result. None represents NaN by class. -/
def coshWord (arm64 expFused : Bool) (word : BitVec 64) : Option Nat :=
  let magnitude := word.toNat % 2 ^ 63
  if 0x7ff0000000000000 < magnitude then none
  else
    let exponential := expWord arm64 expFused (BitVec.ofNat 64 magnitude)
    if exponential = 0x7ff0000000000000 then some exponential
    else if 0x4035000000000000 < magnitude then
      some (multiplyWords exponential 0x3fe0000000000000)
    else
      let reciprocal := divideWords 0x3ff0000000000000 exponential
      some (multiplyWords (addWords exponential reciprocal) 0x3fe0000000000000)

theorem coshWord_even (arm64 expFused : Bool) (x y : BitVec 64)
    (magnitude : x.toNat % 2 ^ 63 = y.toNat % 2 ^ 63) :
    coshWord arm64 expFused x = coshWord arm64 expFused y := by
  simp only [coshWord, magnitude]

theorem coshWord_nan (arm64 expFused : Bool) (word : BitVec 64)
    (nan : 0x7ff0000000000000 < word.toNat % 2 ^ 63) :
    coshWord arm64 expFused word = none := by
  simp [coshWord, nan]

set_option exponentiation.threshold 4096 in
set_option maxRecDepth 4096 in
theorem coshWord_specials (arm64 expFused : Bool) :
    coshWord arm64 expFused 0 = some 0x3ff0000000000000 ∧
      coshWord arm64 expFused (BitVec.ofNat 64 0x8000000000000000) = some 0x3ff0000000000000 ∧
      coshWord arm64 expFused (BitVec.ofNat 64 0x7ff0000000000000) = some 0x7ff0000000000000 ∧
      coshWord arm64 expFused (BitVec.ofNat 64 0xfff0000000000000) = some 0x7ff0000000000000 := by
  cases arm64 <;> cases expFused <;> decide

end GoLean.Math
