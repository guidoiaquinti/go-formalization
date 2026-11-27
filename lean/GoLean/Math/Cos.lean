import GoLean.Math.Sin

namespace GoLean.Math

/-- Go 1.27.0 src/math/sin.go:118–176. The cosine uses the same
    reductions and coefficient chains as sine, with a different
    octant-to-polynomial mapping and reflection rule. -/
def cosWord (arm64 fused : Bool) (word : BitVec 64) : Option Nat :=
  let magnitude := word.toNat % 2 ^ 63
  if 0x7ff0000000000000 ≤ magnitude then none
  else
    let (octant, reduced) := if magnitude < 0x41c0000000000000 then reduceSmallSin arm64 magnitude
      else reduceLargeSin magnitude
    let result := if octant % 4 = 2 then sinPolynomial arm64 fused reduced else cosPolynomial arm64 fused reduced
    some (if decide (4 ≤ octant) != decide (2 ≤ octant % 4) then negateWord result else result)

theorem cosWord_nonfinite (arm64 fused : Bool) (word : BitVec 64)
    (nonfinite : 0x7ff0000000000000 ≤ word.toNat % 2 ^ 63) :
    cosWord arm64 fused word = none := by
  simp [cosWord, nonfinite]

theorem cosWord_even (arm64 fused : Bool) (x y : BitVec 64)
    (same : x.toNat % 2 ^ 63 = y.toNat % 2 ^ 63) :
    cosWord arm64 fused x = cosWord arm64 fused y := by
  simp only [cosWord, same]

set_option exponentiation.threshold 4096 in
set_option maxRecDepth 4096 in
theorem cosWord_zero (arm64 fused : Bool) :
    cosWord arm64 fused (0 : BitVec 64) = some 0x3ff0000000000000 ∧
    cosWord arm64 fused (BitVec.ofNat 64 (2 ^ 63)) = some 0x3ff0000000000000 := by
  cases arm64 <;> cases fused <;> decide

end GoLean.Math
