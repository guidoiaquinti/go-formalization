import GoLean.Math.Sin

namespace GoLean.Math

private def sincosPositive (arm64 fused : Bool) (magnitude : Nat) : Nat × Nat :=
  let (octant, reduced) := if magnitude < 0x41c0000000000000 then reduceSmallSin arm64 magnitude
    else reduceLargeSin magnitude
  let cosine := cosPolynomial arm64 fused reduced
  let sine := sinPolynomial arm64 fused reduced
  let reflectedOctant := octant % 4
  let (sine, cosine) := if reflectedOctant = 1 || reflectedOctant = 2 then (cosine, sine)
    else (sine, cosine)
  let reflect := 4 ≤ octant
  (if reflect then negateWord sine else sine,
    if decide reflect != decide (1 < reflectedOctant) then negateWord cosine else cosine)

/-- Go 1.27.0 src/math/sincos.go:16–72. Compute both polynomials from one
    argument reduction, then exchange the outputs and reflect their signs.
    None represents two NaNs, observed by class. -/
def sincosWords (arm64 fused : Bool) (word : BitVec 64) : Option (Nat × Nat) :=
  let raw := word.toNat
  let magnitude := raw % 2 ^ 63
  if magnitude = 0 then some (raw, 0x3ff0000000000000)
  else if 0x7ff0000000000000 ≤ magnitude then none
  else
    let (sine, cosine) := sincosPositive arm64 fused magnitude
    some (if 2 ^ 63 ≤ raw then negateWord sine else sine, cosine)

theorem sincosWords_zero (arm64 fused : Bool) (word : BitVec 64)
    (zero : word.toNat % 2 ^ 63 = 0) :
    sincosWords arm64 fused word = some (word.toNat, 0x3ff0000000000000) := by
  simp [sincosWords, zero]

theorem sincosWords_nonfinite (arm64 fused : Bool) (word : BitVec 64)
    (nonfinite : 0x7ff0000000000000 ≤ word.toNat % 2 ^ 63) :
    sincosWords arm64 fused word = none := by
  simp [sincosWords, nonfinite, show word.toNat % 2 ^ 63 ≠ 0 by omega]

theorem sincosWords_opposite (arm64 fused : Bool) (x y : BitVec 64)
    (positive : x.toNat < 2 ^ 63) (negative : 2 ^ 63 ≤ y.toNat)
    (same : x.toNat % 2 ^ 63 = y.toNat % 2 ^ 63)
    (finite : x.toNat % 2 ^ 63 < 0x7ff0000000000000)
    (nonzero : x.toNat % 2 ^ 63 ≠ 0) :
    sincosWords arm64 fused y =
      (sincosWords arm64 fused x).map (fun (sine, cosine) => (negateWord sine, cosine)) := by
  simp [sincosWords, ← same, nonzero, Nat.not_le.mpr finite, Nat.not_le.mpr positive, negative]

end GoLean.Math
