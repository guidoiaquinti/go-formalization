import GoLean.Math.Log1p

namespace GoLean.Math

/-- Go 1.27.0 src/math/atanh.go:48–84, for the generic public dispatch.
    The two arithmetic profiles select the corresponding Log1p kernel.
    None represents NaN by class. -/
def atanhWord (arm64 fused : Bool) (word : BitVec 64) : Option Nat :=
  let x := word.toNat
  let magnitude := x % 2 ^ 63
  let signed := fun result => if 2 ^ 63 ≤ x then negateWord result else result
  if 0x3ff0000000000000 < magnitude then none
  else if magnitude = 0x3ff0000000000000 then some (signed 0x7ff0000000000000)
  else if magnitude < 0x3e30000000000000 then some x
  else
    let doubled := addWords magnitude magnitude
    let denominator := subtractWords 0x3ff0000000000000 magnitude
    let argument :=
      if magnitude < 0x3fe0000000000000 then
        addWords doubled (divideWords (multiplyWords doubled magnitude) denominator)
      else divideWords doubled denominator
    (log1pWord arm64 fused (BitVec.ofNat 64 argument)).map
      (fun result => signed (multiplyWords 0x3fe0000000000000 result))

theorem atanhWord_tiny (arm64 fused : Bool) (word : BitVec 64)
    (tiny : word.toNat % 2 ^ 63 < 0x3e30000000000000) :
    atanhWord arm64 fused word = some word.toNat := by
  unfold atanhWord
  dsimp only
  rw [ite_eq_right (by omega), ite_eq_right (by omega), ite_eq_left tiny]

theorem atanhWord_outside (arm64 fused : Bool) (word : BitVec 64)
    (outside : 0x3ff0000000000000 < word.toNat % 2 ^ 63) :
    atanhWord arm64 fused word = none := by
  simp [atanhWord, outside]

theorem atanhWord_endpoints (arm64 fused : Bool) :
    atanhWord arm64 fused (BitVec.ofNat 64 0x3ff0000000000000) = some 0x7ff0000000000000 ∧
      atanhWord arm64 fused (BitVec.ofNat 64 0xbff0000000000000) = some 0xfff0000000000000 := by
  cases arm64 <;> cases fused <;> decide

end GoLean.Math
