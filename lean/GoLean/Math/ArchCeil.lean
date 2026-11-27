import GoLean.Math.Trunc

namespace GoLean.Math

/-- Go 1.27.0 src/math/floor_arm64.s:14–19 and src/math/floor_amd64.s:31–56.
    FRINTPD rounds toward positive infinity. The amd64 path preserves large
    words before integer conversion, restores the sign of a zero integer
    part, and increments only a positive nonintegral input. NaNs are observed
    by class. -/
def archCeil (amd64 : Bool) (word : BitVec 64) : Nat :=
  let bits := word.toNat
  let magnitude := bits % 2 ^ 63
  if amd64 && (0x4330000000000000 ≤ magnitude) then bits
  else
    let e := bits / 2 ^ 52 % 2048
    let whole := trunc word
    if !word.msb && whole != bits then
      whole + (if e < 1023 then 1023 * 2 ^ 52 else 2 ^ (1075 - e))
    else whole

theorem archCeil_spec (amd64 : Bool) (word : BitVec 64) :
    archCeil amd64 word =
      let bits := word.toNat
      let e := bits / 2 ^ 52 % 2048
      if e < 1023 then
        let signedZero := bits - bits % 2 ^ 63
        if !word.msb && bits % 2 ^ 63 != 0 then signedZero + 1023 * 2 ^ 52
        else signedZero
      else if e < 1075 then
        let unit := 2 ^ (1075 - e)
        (bits / unit + (if !word.msb && bits % unit != 0 then 1 else 0)) * unit
      else bits := by
  unfold archCeil
  dsimp only
  split
  · rename_i guard
    have large : 0x4330000000000000 ≤ word.toNat % 2 ^ 63 := by
      simpa using (Bool.and_eq_true_iff.mp guard).2
    have exp : word.toNat / 2 ^ 52 % 2048 = word.toNat % 2 ^ 63 / 2 ^ 52 := by omega
    have integral : ¬ word.toNat % 2 ^ 63 / 2 ^ 52 < 1075 := by omega
    have notSmall : ¬ word.toNat % 2 ^ 63 / 2 ^ 52 < 1023 := by omega
    simp [exp, integral, notSmall]
  · rw [trunc_spec]
    dsimp only
    split
    · rename_i small
      have eq : (word.toNat - word.toNat % 2 ^ 63 = word.toNat) ↔
          word.toNat % 2 ^ 63 = 0 := by omega
      simp [bne, eq]
    · rename_i large
      split
      · rename_i fractional
        let unit := 2 ^ (1075 - word.toNat / 2 ^ 52 % 2048)
        have eq : (word.toNat / unit * unit = word.toNat) ↔ word.toNat % unit = 0 := by
          have h := Nat.mod_add_div' word.toNat unit
          omega
        change (if !word.msb && (word.toNat / unit * unit != word.toNat)
          then word.toNat / unit * unit + unit else word.toNat / unit * unit) =
          (word.toNat / unit + (if !word.msb && word.toNat % unit != 0 then 1 else 0)) * unit
        by_cases z : word.toNat % unit = 0
        · simp [z, eq.mpr z]
        · cases word.msb <;> simp [z, mt eq.mp z, Nat.add_mul]

      · rename_i integral
        simp

end GoLean.Math
