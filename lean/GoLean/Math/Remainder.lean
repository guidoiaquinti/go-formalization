import GoLean.Math.Mod

namespace GoLean.Math

private def centerRemainder (r y : Int) (tiny : Bool) : Int :=
  if (if tiny then 2 * r > y else r > y / 2) then
    let once := r - y
    if (if tiny then 2 * once ≥ y else once ≥ y / 2) then once - y else once
  else r

/-- Go 1.27.0 src/math/remainder.go:45–88. Use exact counts of the
    smallest subnormal to represent the native reduction and comparisons.
    Reduce modulo twice the divisor only if doubling fits. The tiny path
    doubles the residual instead of halving the divisor. -/
def remainderWords (x y : BitVec 64) : Option Nat :=
  let xm := x.toNat % 2 ^ 63
  let ym := y.toNat % 2 ^ 63
  if ym = 0 ∨ 0x7ff0000000000000 ≤ xm ∨ 0x7ff0000000000000 < ym then none
  else if ym = 0x7ff0000000000000 then some x.toNat
  else
    let sign := x.toNat - xm
    if xm = ym then some sign
    else
      let dividend := magnitudeUnits xm
      let divisor := magnitudeUnits ym
      let reduced := if ym ≤ 0x7fdfffffffffffff then modReduce dividend (2 * divisor) else dividend
      let result := centerRemainder reduced divisor (ym < 0x0020000000000000)
      let resultSign := if result < 0 then 2 ^ 63 - sign else sign
      some (resultSign + encodeUnits result.natAbs)

private def nearestResidue (x y : Nat) : Int :=
  let r := x % y
  if y < 2 * r ∨ (2 * r = y ∧ x / y % 2 = 1) then (r : Int) - y
  else r

private theorem centerRemainder_spec (x y : Nat) (tiny : Bool)
    (pos : 0 < y) (halfExact : tiny = true ∨ y % 2 = 0) :
    centerRemainder (x % (2 * y) : Nat) y tiny = nearestResidue x y := by
  have period : y ∣ 2 * y := ⟨2, Nat.mul_comm ..⟩
  have low := Nat.mod_mod_of_dvd x period
  have quotient := Nat.mod_mul_left_div_self x y 2
  have together := Nat.mod_add_div (x % (2 * y)) y
  rw [low, quotient] at together
  have residueBound := Nat.mod_lt x pos
  have parityBound : x / y % 2 < 2 := Nat.mod_lt _ (by decide)
  unfold nearestResidue
  have bound := Nat.mod_lt x (show 0 < 2 * y by omega)
  generalize x % (2 * y) = big at *
  generalize x % y = small at *
  generalize x / y % 2 = parity at *
  have decomposition : big = small + y * parity := together.symm
  clear together low quotient
  have evenOrOdd : parity = 0 ∨ parity = 1 := by omega
  rcases evenOrOdd with eq | eq
  all_goals subst parity
  all_goals simp only [Nat.mul_zero, Nat.mul_one, Nat.add_zero] at decomposition
  all_goals cases tiny
  all_goals simp at halfExact
  all_goals simp [centerRemainder]
  all_goals split <;> (repeat' split) <;> omega

private def signedEncode (sign : Nat) (result : Int) : Nat :=
  (if result < 0 then 2 ^ 63 - sign else sign) + encodeUnits result.natAbs

theorem remainderWords_finite (x y : BitVec 64)
    (valid : ¬ (y.toNat % 2 ^ 63 = 0 ∨ 0x7ff0000000000000 ≤ x.toNat % 2 ^ 63 ∨
      0x7ff0000000000000 < y.toNat % 2 ^ 63))
    (finite : y.toNat % 2 ^ 63 ≠ 0x7ff0000000000000)
    (unequal : x.toNat % 2 ^ 63 ≠ y.toNat % 2 ^ 63)
    (dyadicPositive : 0 < magnitudeUnits (y.toNat % 2 ^ 63))
    (halfExact : y.toNat % 2 ^ 63 < 0x0020000000000000 ∨ magnitudeUnits (y.toNat % 2 ^ 63) % 2 = 0)
    (safePeriod : 0x7fdfffffffffffff < y.toNat % 2 ^ 63 →
      magnitudeUnits (x.toNat % 2 ^ 63) < 2 * magnitudeUnits (y.toNat % 2 ^ 63)) :
    remainderWords x y =
      let xm := x.toNat % 2 ^ 63
      let ym := y.toNat % 2 ^ 63
      some (signedEncode (x.toNat - xm) (nearestResidue (magnitudeUnits xm) (magnitudeUnits ym))) := by
  unfold remainderWords
  dsimp only
  rw [ite_eq_right valid, ite_eq_right finite, ite_eq_right unequal]
  have reduction : (if y.toNat % 2 ^ 63 ≤ 0x7fdfffffffffffff then
        modReduce (magnitudeUnits (x.toNat % 2 ^ 63)) (2 * magnitudeUnits (y.toNat % 2 ^ 63))
        else magnitudeUnits (x.toNat % 2 ^ 63)) =
      magnitudeUnits (x.toNat % 2 ^ 63) % (2 * magnitudeUnits (y.toNat % 2 ^ 63)) := by
    split
    · exact modReduce_spec _ _ (by omega)
    · exact (Nat.mod_eq_of_lt (safePeriod (by omega))).symm
  rw [reduction, centerRemainder_spec _ _ _ dyadicPositive (by simpa using halfExact)]
  rfl

end GoLean.Math
