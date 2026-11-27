import GoLean.Math.Log
import GoLean.Math.Log1p

namespace GoLean.Math

/-- Go 1.27.0 src/math/asinh.go:40–77, for the generic public dispatch.
    Arm64 and amd64 v3/v4 fuse the two square-plus-one radicands.
    None represents NaN by class. -/
def asinhWord (arm64 fused : Bool) (word : BitVec 64) : Option Nat :=
  let x := word.toNat
  let magnitude := x % 2 ^ 63
  let signed := fun result => if 2 ^ 63 ≤ x then negateWord result else result
  if 0x7ff0000000000000 < magnitude then none
  else if magnitude = 0x7ff0000000000000 then some x
  else if magnitude < 0x3e30000000000000 then some x
  else if 0x41b0000000000000 < magnitude then
    (logWord arm64 (BitVec.ofNat 64 magnitude)).map
      (fun result => signed (addWords result 0x3fe62e42fefa39ef))
  else do
    let one := 0x3ff0000000000000
    let radicand := multiplyAddWords (arm64 || fused) magnitude magnitude one
    let root ← sqrtWord (BitVec.ofNat 64 radicand)
    let value ←
      if 0x4000000000000000 < magnitude then
        let reciprocal := divideWords one (addWords root magnitude)
        logWord arm64 (BitVec.ofNat 64 (addWords (addWords magnitude magnitude) reciprocal))
      else
        let square := multiplyWords magnitude magnitude
        let correction := divideWords square (addWords one root)
        log1pWord arm64 fused (BitVec.ofNat 64 (addWords magnitude correction))
    return signed value

theorem asinhWord_tiny (arm64 fused : Bool) (word : BitVec 64)
    (tiny : word.toNat % 2 ^ 63 < 0x3e30000000000000) :
    asinhWord arm64 fused word = some word.toNat := by
  unfold asinhWord
  dsimp only
  rw [ite_eq_right (by omega), ite_eq_right (by omega), ite_eq_left tiny]

theorem asinhWord_nan (arm64 fused : Bool) (word : BitVec 64)
    (nan : 0x7ff0000000000000 < word.toNat % 2 ^ 63) :
    asinhWord arm64 fused word = none := by
  simp [asinhWord, nan]

theorem asinhWord_infinity (arm64 fused : Bool) (word : BitVec 64)
    (infinity : word.toNat % 2 ^ 63 = 0x7ff0000000000000) :
    asinhWord arm64 fused word = some word.toNat := by
  simp [asinhWord, infinity]

end GoLean.Math
