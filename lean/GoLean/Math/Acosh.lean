import GoLean.Math.Log
import GoLean.Math.Log1p

namespace GoLean.Math

/-- Go 1.27.0 src/math/acosh.go:43–65, for the generic public dispatch.
    Arm64 contracts both radicands. Amd64 v3/v4 contracts the near-one
    radicand while retaining separate rounding in square-minus-one.
    None represents NaN by class. -/
def acoshWord (arm64 fused : Bool) (word : BitVec 64) : Option Nat :=
  let x := word.toNat
  let one := 0x3ff0000000000000
  if x < one ∨ 0x7ff0000000000000 < x then none
  else if x = one then some 0
  else if x = 0x7ff0000000000000 then some x
  else if 0x41b0000000000000 ≤ x then
    (logWord arm64 word).map (fun result => addWords result 0x3fe62e42fefa39ef)
  else if 0x4000000000000000 < x then do
    let radicand := multiplyAddWords arm64 x x 0xbff0000000000000
    let root ← sqrtWord (BitVec.ofNat 64 radicand)
    let reciprocal := divideWords one (addWords x root)
    logWord arm64 (BitVec.ofNat 64 (subtractWords (addWords x x) reciprocal))
  else do
    let t := subtractWords x one
    let radicand := multiplyAddWords (arm64 || fused) t t (addWords t t)
    let root ← sqrtWord (BitVec.ofNat 64 radicand)
    log1pWord arm64 fused (BitVec.ofNat 64 (addWords t root))

theorem acoshWord_outside (arm64 fused : Bool) (word : BitVec 64)
    (outside : word.toNat < 0x3ff0000000000000 ∨ 0x7ff0000000000000 < word.toNat) :
    acoshWord arm64 fused word = none := by
  simp [acoshWord, outside]

theorem acoshWord_edges (arm64 fused : Bool) :
    acoshWord arm64 fused (BitVec.ofNat 64 0x3ff0000000000000) = some 0 ∧
      acoshWord arm64 fused (BitVec.ofNat 64 0x7ff0000000000000) = some 0x7ff0000000000000 := by
  cases arm64 <;> cases fused <;> decide

end GoLean.Math
