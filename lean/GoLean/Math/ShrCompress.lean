import GoLean.Math.Shr
import GoLean.Math.FmaNonzero

namespace GoLean.Math

private def lowMask (n : Nat) : BitVec 64 :=
  BitVec.ofNat 64 (2 ^ n - 1)

/-- The generic helper first shifts the two words, then ORs the nonzero
    indicator for discarded bits into the low result. Dedicated source
    branches handle counts of zero, sixty-four, and at least 128.
    https://github.com/golang/go/blob/go1.27.0/src/math/fma.go#L18-L25
    https://github.com/golang/go/blob/go1.27.0/src/math/fma.go#L39-L63 -/
def shrcompress (high low : BitVec 64) (n : Nat) : BitVec 64 × BitVec 64 :=
  if n = 0 then (high, low)
  else if n = 64 then (0, high ||| fmaNonzero low)
  else if n ≥ 128 then (0, fmaNonzero (high ||| low))
  else
    let shifted := shr high low n
    let discarded :=
      if n < 64 then low &&& lowMask n
      else (high &&& lowMask (n - 64)) ||| low
    (shifted.1, shifted.2 ||| fmaNonzero discarded)

private theorem joined (high low : BitVec 64) :
    (high ++ low).toNat = high.toNat * 2 ^ 64 + low.toNat := by
  rw [BitVec.toNat_append, ← Nat.shiftLeft_add_eq_or_of_lt low.isLt, Nat.shiftLeft_eq]

private theorem mask_zero (word : BitVec 64) (h : n ≤ 64) :
    word &&& lowMask n = 0 ↔ word.toNat % 2 ^ n = 0 := by
  have bound : 2 ^ n - 1 < 2 ^ 64 := by
    have := Nat.pow_le_pow_right Nat.zero_lt_two h
    have := Nat.two_pow_pos n
    omega
  rw [BitVec.toNat_eq, BitVec.toNat_and, lowMask, BitVec.toNat_ofNat,
    Nat.mod_eq_of_lt bound, Nat.and_two_pow_sub_one_eq_mod]
  rfl

private theorem short_remainder (high low : BitVec 64) (h : n ≤ 64) :
    (high ++ low).toNat % 2 ^ n = low.toNat % 2 ^ n := by
  rw [joined, Nat.add_comm, Nat.add_mod, Nat.mul_mod,
    Nat.mod_eq_zero_of_dvd (Nat.pow_dvd_pow 2 h), Nat.mul_zero, Nat.zero_mod,
    Nat.add_zero, Nat.mod_mod]

private theorem long_remainder_zero (high low : BitVec 64) (h : 64 ≤ n) :
    (high ++ low).toNat % 2 ^ n = 0 ↔
      high.toNat % 2 ^ (n - 64) = 0 ∧ low = 0 := by
  obtain ⟨k, rfl⟩ := Nat.le.dest h
  rw [show 64 + k - 64 = k by omega, Nat.pow_add]
  let rem := (high ++ low).toNat % (2 ^ 64 * 2 ^ k)
  have lower : rem % 2 ^ 64 = low.toNat := by
    simp only [rem, Nat.mod_mul_right_mod, short_remainder high low (Nat.le_refl _)]
    exact Nat.mod_eq_of_lt low.isLt
  have upper : rem / 2 ^ 64 = high.toNat % 2 ^ k := by
    simp only [rem, Nat.mod_mul_right_div_self, joined]
    rw [Nat.add_comm, Nat.mul_comm high.toNat,
      Nat.add_mul_div_left _ _ (Nat.two_pow_pos _), Nat.div_eq_of_lt low.isLt, Nat.zero_add]
  have combined := Nat.mod_add_div' rem (2 ^ 64)
  rw [lower, upper] at combined
  change rem = 0 ↔ _
  rw [BitVec.toNat_eq]
  change _ ↔ _ ∧ low.toNat = 0
  omega

private theorem append_sticky (high low discarded : BitVec 64) :
    high ++ (low ||| fmaNonzero discarded) =
      (high ++ low) ||| (if discarded = 0 then 0#128 else 1#128) := by
  have extension : (0#64 ++ fmaNonzero discarded) =
      (if discarded = 0 then 0#128 else 1#128) := by
    by_cases empty : discarded = 0
    · simp [fmaNonzero, empty]
    · rw [fmaNonzero, ite_eq_right empty, ite_eq_right empty]
      decide
  rw [← extension, BitVec.or_append, BitVec.or_zero]

/-- The two result words are the shifted combined input, with bit zero
    additionally set whenever the discarded low bits are nonzero. -/
theorem shrcompress_spec (high low : BitVec 64) (n : Nat) :
    (shrcompress high low n).1 ++ (shrcompress high low n).2 =
      ((high ++ low) >>> n) |||
        (if (high ++ low).toNat % 2 ^ n = 0 then 0#128 else 1#128) := by
  unfold shrcompress
  by_cases zero : n = 0
  · subst n
    simp [Nat.mod_one]
  · rw [ite_eq_right zero]
    by_cases middle : n = 64
    · subst n
      have shifted : (high ++ low) >>> 64 = (0#64 ++ high) := by
        rw [← shr_spec]
        simp [shr, BitVec.ushiftRight_eq_zero (Nat.le_refl 64)]
      rw [ite_eq_left rfl, shifted, append_sticky,
        short_remainder high low (Nat.le_refl _), Nat.mod_eq_of_lt low.isLt]
      simp only [show low = 0 ↔ low.toNat = 0 from BitVec.toNat_eq]
      rfl
    · rw [ite_eq_right middle]
      by_cases large : n ≥ 128
      · rw [ite_eq_left large, BitVec.ushiftRight_eq_zero large, BitVec.zero_or]
        have bound : (high ++ low).toNat < 2 ^ n :=
          Nat.lt_of_lt_of_le (high ++ low).isLt (Nat.pow_le_pow_right Nat.zero_lt_two large)
        have empty : (high ||| low) = 0 ↔ (high ++ low).toNat = 0 := by
          change (high ||| low) = 0#64 ↔ _
          rw [BitVec.or_eq_zero_iff, BitVec.toNat_eq, BitVec.toNat_eq, joined]
          change high.toNat = 0 ∧ low.toNat = 0 ↔ _
          omega
        rw [Nat.mod_eq_of_lt bound]
        by_cases allZero : (high ||| low) = 0
        · simp [fmaNonzero, allZero, empty.mp allZero]
        · simp only [fmaNonzero, ite_eq_right allZero, ite_eq_right (mt empty.mpr allZero)]
          decide
      · rw [ite_eq_right large]
        dsimp only
        rw [append_sticky, shr_spec]
        by_cases narrow : n < 64
        · rw [ite_eq_left narrow, short_remainder high low (Nat.le_of_lt narrow)]
          simp only [mask_zero low (Nat.le_of_lt narrow)]
        · rw [ite_eq_right narrow]
          have discardedEmpty : ((high &&& lowMask (n - 64)) ||| low) = 0 ↔
              high.toNat % 2 ^ (n - 64) = 0 ∧ low = 0 := by
            change ((high &&& lowMask (n - 64)) ||| low) = 0#64 ↔ _
            rw [BitVec.or_eq_zero_iff]
            exact and_congr (mask_zero high (n := n - 64) (by omega)) Iff.rfl
          simp only [long_remainder_zero high low (by omega : 64 ≤ n), discardedEmpty]

end GoLean.Math
