import GoLean.Math.Exp

namespace GoLean.Math

private def sinhRational (fused : Bool) (x : Nat) : Nat :=
  let square := multiplyWords x x
  let numerator := multiplyAddWords fused 0xc03a4e3de8540779 square 0xc0a69c6c36da2dfb
  let numerator := multiplyAddWords fused numerator square 0xc0f5f38b8605d22d
  let numerator := multiplyAddWords fused numerator square 0xc1233fdeba64bb4f
  let denominator := multiplyAddWords fused x x 0xc065b5b9fcd003bb
  let denominator := multiplyAddWords fused denominator square 0x40cdb7963eae91e1
  let denominator := multiplyAddWords fused denominator square 0xc1233fdeba64bb4f
  divideWords (multiplyWords numerator x) denominator

/-- Go 1.27.0 src/math/sinh.go:26–70. Compiler fusion applies in the
    rational branch, and amd64 runtime fusion applies independently in Exp.
    None represents NaN by class. -/
def sinhWord (arm64 compileFused expFused : Bool) (word : BitVec 64) : Option Nat :=
  let x := word.toNat
  let magnitude := x % 2 ^ 63
  let signed := fun result => if 2 ^ 63 ≤ x then negateWord result else result
  if 0x7ff0000000000000 < magnitude then none
  else if magnitude = 0 then some x
  else if 0x3fe0000000000000 < magnitude then
    let exponential := expWord arm64 expFused (BitVec.ofNat 64 magnitude)
    let value :=
      if exponential = 0x7ff0000000000000 then exponential
      else if 0x4035000000000000 < magnitude then multiplyWords exponential 0x3fe0000000000000
      else
        let reciprocal := divideWords 0x3ff0000000000000 exponential
        multiplyWords (subtractWords exponential reciprocal) 0x3fe0000000000000
    some (signed value)
  else some (signed (sinhRational (arm64 || compileFused) magnitude))

theorem sinhWord_zero (arm64 compileFused expFused : Bool) (word : BitVec 64)
    (zero : word.toNat % 2 ^ 63 = 0) :
    sinhWord arm64 compileFused expFused word = some word.toNat := by
  simp [sinhWord, zero]

theorem sinhWord_nan (arm64 compileFused expFused : Bool) (word : BitVec 64)
    (nan : 0x7ff0000000000000 < word.toNat % 2 ^ 63) :
    sinhWord arm64 compileFused expFused word = none := by
  simp [sinhWord, nan]

theorem sinhWord_infinities (arm64 compileFused expFused : Bool) :
    sinhWord arm64 compileFused expFused (BitVec.ofNat 64 0x7ff0000000000000) = some 0x7ff0000000000000 ∧
      sinhWord arm64 compileFused expFused (BitVec.ofNat 64 0xfff0000000000000) = some 0xfff0000000000000 := by
  cases arm64 <;> cases compileFused <;> cases expFused <;> decide

end GoLean.Math
