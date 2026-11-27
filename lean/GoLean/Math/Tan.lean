import GoLean.Math.Sin

namespace GoLean.Math

private def tanPolynomial (arm64 fused : Bool) (z : Nat) : Nat :=
  let zz := multiplyWords z z
  if zz ≤ 0x3d06849b86a12b9b then z
  else
    let p := multiplyAddWords arm64 0xc0c992d8d24f3f38 zz 0x413199eca5fc9ddd
    let p := multiplyAddWords arm64 p zz 0xc1711fead3299176
    let q := multiplyAddWords (arm64 || fused) z z 0x40cab8a5eeb36572
    let q := multiplyAddWords arm64 q zz 0xc13427bc582abc96
    let q := multiplyAddWords arm64 q zz 0x4177d98fc2ead8ef
    let q := multiplyAddWords arm64 q zz 0xc189afe03cbe5a31
    let correction := divideWords (multiplyWords zz p) q
    multiplyAddWords (arm64 || fused) z correction z

private def tanReciprocal (word : Nat) : Nat :=
  if word % 2 ^ 63 = 0 then negateWord word + 0x7ff0000000000000
  else divideWords 0xbff0000000000000 word

/-- Go 1.27.0 src/math/tan.go:83–140. Reuse the split and Payne-Hanek
    reductions, evaluate the rational approximation outside the tiny
    reduced interval, and take a negative reciprocal in odd quadrants. -/
def tanWord (arm64 fused : Bool) (word : BitVec 64) : Option Nat :=
  let raw := word.toNat
  let magnitude := raw % 2 ^ 63
  if magnitude = 0 then some raw
  else if 0x7ff0000000000000 ≤ magnitude then none
  else
    let (octant, reduced) := if magnitude < 0x41c0000000000000 then reduceSmallSin arm64 magnitude
      else reduceLargeSin magnitude
    let withinQuadrant := tanPolynomial arm64 fused reduced
    let result := if octant % 4 = 2 then tanReciprocal withinQuadrant else withinQuadrant
    some (if 2 ^ 63 ≤ raw then negateWord result else result)

theorem tanWord_zero (arm64 fused : Bool) (word : BitVec 64)
    (zero : word.toNat % 2 ^ 63 = 0) :
    tanWord arm64 fused word = some word.toNat := by
  simp [tanWord, zero]

theorem tanWord_nonfinite (arm64 fused : Bool) (word : BitVec 64)
    (nonfinite : 0x7ff0000000000000 ≤ word.toNat % 2 ^ 63) :
    tanWord arm64 fused word = none := by
  simp [tanWord, nonfinite, show word.toNat % 2 ^ 63 ≠ 0 by omega]

theorem tanWord_opposite (arm64 fused : Bool) (x y : BitVec 64)
    (positive : x.toNat < 2 ^ 63) (negative : 2 ^ 63 ≤ y.toNat)
    (same : x.toNat % 2 ^ 63 = y.toNat % 2 ^ 63)
    (finite : x.toNat % 2 ^ 63 < 0x7ff0000000000000)
    (nonzero : x.toNat % 2 ^ 63 ≠ 0) :
    tanWord arm64 fused y = (tanWord arm64 fused x).map negateWord := by
  simp only [tanWord, ← same, nonzero, Nat.not_le.mpr finite, negative,
    Nat.not_le.mpr positive, ite_true, ite_false, Option.map_some]

end GoLean.Math
