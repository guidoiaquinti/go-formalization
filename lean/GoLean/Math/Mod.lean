import Init.Data.BitVec.Lemmas
import Init.Data.Nat.Lemmas

namespace GoLean.Math

private def modShift (r y : Nat) : Nat :=
  let shift := r.log2 - y.log2
  if r < y * 2 ^ shift then shift - 1 else shift

private theorem modStep_le (r y : Nat) (hy : 0 < y) (below : y ≤ r) :
    y * 2 ^ modShift r y ≤ r := by
  unfold modShift
  dsimp only
  split
  · rename_i tooLarge
    have bound := Nat.lt_log2_self (n := y)
    have low := Nat.log2_self_le (n := r) (by omega)
    have shiftPositive : 0 < r.log2 - y.log2 := by
      by_cases hzero : r.log2 - y.log2 = 0
      · simp [hzero] at tooLarge
        omega
      · omega
    have powers : y.log2 + 1 + (r.log2 - y.log2 - 1) = r.log2 := by omega
    have scaled := Nat.mul_lt_mul_of_pos_right bound
      (Nat.two_pow_pos (r.log2 - y.log2 - 1))
    rw [← Nat.pow_add, powers] at scaled
    exact Nat.le_trans (Nat.le_of_lt scaled) low
  · omega

/-- Go 1.27.0 src/math/mod.go:29–52. Use integer coefficients in units of
    the smallest subnormal. The exponent difference selects the first
    shifted divisor, and the significand comparison lowers it when it
    would exceed the current residual. -/
def modReduce (r y : Nat) : Nat :=
  if hy : 0 < y then
    if below : y ≤ r then modReduce (r - y * 2 ^ modShift r y) y
    else r
  else r
termination_by r
decreasing_by
  have positive : 0 < y * 2 ^ modShift r y := Nat.mul_pos hy (Nat.two_pow_pos _)
  omega

theorem modReduce_spec (r y : Nat) (hy : 0 < y) :
    modReduce r y = r % y := by
  rw [modReduce.eq_def, dite_eq_left hy]
  by_cases below : y ≤ r
  · rw [dite_eq_left below, modReduce_spec _ _ hy, Nat.sub_mul_mod (modStep_le r y hy below)]
  · rw [dite_eq_right below, Nat.mod_eq_of_lt (by omega)]
termination_by r
decreasing_by
  have positive : 0 < y * 2 ^ modShift r y := Nat.mul_pos hy (Nat.two_pow_pos _)
  omega

def magnitudeUnits (bits : Nat) : Nat :=
  let exponent := bits / 2 ^ 52 % 2048
  let fraction := bits % 2 ^ 52
  if exponent = 0 then fraction else (2 ^ 52 + fraction) * 2 ^ (exponent - 1)

def encodeUnits (units : Nat) : Nat :=
  if units < 2 ^ 52 then units
  else
    let top := units.log2
    (top - 51) * 2 ^ 52 + (units / 2 ^ (top - 52)) % 2 ^ 52

/-- The pinned wrapper selects the generic path on the supported native
    profiles. None represents NaN. Preserve the original sign of x even
    when its remainder magnitude is zero. -/
def modWords (x y : BitVec 64) : Option Nat :=
  let xm := x.toNat % 2 ^ 63
  let ym := y.toNat % 2 ^ 63
  if ym = 0 ∨ 0x7ff0000000000000 ≤ xm ∨ 0x7ff0000000000000 < ym then none
  else if ym = 0x7ff0000000000000 then some x.toNat
  else some (x.toNat - xm + encodeUnits (modReduce (magnitudeUnits xm) (magnitudeUnits ym)))

private theorem magnitudeUnits_positive (bits : Nat) (pos : 0 < bits)
    (bounded : bits < 2 ^ 63) : 0 < magnitudeUnits bits := by
  unfold magnitudeUnits
  dsimp only
  split
  · omega
  · exact Nat.mul_pos (by omega) (Nat.two_pow_pos _)

theorem modWords_spec (x y : BitVec 64) :
    modWords x y =
      let xm := x.toNat % 2 ^ 63
      let ym := y.toNat % 2 ^ 63
      if ym = 0 ∨ 0x7ff0000000000000 ≤ xm ∨ 0x7ff0000000000000 < ym then none
      else if ym = 0x7ff0000000000000 then some x.toNat
      else some (x.toNat - xm + encodeUnits (magnitudeUnits xm % magnitudeUnits ym)) := by
  unfold modWords
  dsimp only
  split
  · rfl
  · rename_i neither
    split
    · rfl
    · have yPositive : 0 < y.toNat % 2 ^ 63 := by omega
      rw [modReduce_spec _ _ (magnitudeUnits_positive _ yPositive (Nat.mod_lt _ (by decide)))]

end GoLean.Math
